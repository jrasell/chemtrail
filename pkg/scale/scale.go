package scale

import (
	"net/http"

	"github.com/jrasell/chemtrail/pkg/client"
	serverCfg "github.com/jrasell/chemtrail/pkg/config/server"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/scale/provider"
	aws_asg "github.com/jrasell/chemtrail/pkg/scale/provider/aws-asg"
	noop "github.com/jrasell/chemtrail/pkg/scale/provider/no-op"
	"github.com/jrasell/chemtrail/pkg/scale/resource"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Scale is the interface used to perform scaling checks and trigger activities.
type Scale interface {

	// OKToScale performs a number of safety checks to ensure the scaling request does not break
	// any policy parameters and that the request has a chance to run successfully. The int
	// returned indicates the appropriate HTTP response code, the error will contain any relevant
	// messages which describe the check that failed. If no error is returned, it can be assumed
	// that the request is OK to continue with.
	OKToScale(req *state.ScalingRequest) (int, error)

	// InvokeScaling triggers a scaling activity, all events from this point will be written to the
	// state store. The function is designed to be called asynchronously, therefore there is no
	// return.
	InvokeScaling(req *state.ScalingRequest)
}

type BackendConfig struct {
	Provider      *serverCfg.ProviderConfig
	Logger        zerolog.Logger
	Nomad         *client.Nomad
	NodeResources resource.Handler
	ScaleState    state.ScaleBackend
	PolicyState   state.PolicyBackend
}

type Backend struct {
	logger          zerolog.Logger
	nomad           *client.Nomad
	scaleState      state.ScaleBackend
	policyState     state.PolicyBackend
	resourceHandler resource.Handler

	// clientProvider stores the client provider interface so we can interact and scale these
	// backends. Currently this is only populated during the instantiation of the new backend.
	clientProvider map[state.ClientProvider]provider.ClientProvider

	// eventChan is used to listen and write scaling activity updates to the backend state store.
	eventChan chan *state.EventMessage
}

func NewScaleBackend(cfg *BackendConfig) Scale {
	b := Backend{
		logger:          cfg.Logger,
		nomad:           cfg.Nomad,
		scaleState:      cfg.ScaleState,
		policyState:     cfg.PolicyState,
		resourceHandler: cfg.NodeResources,
		clientProvider:  make(map[state.ClientProvider]provider.ClientProvider),
		eventChan:       make(chan *state.EventMessage, 10),
	}

	// If the AWS provider is enabled, configure the scaling backend.
	if cfg.Provider.AWSASG {
		b.clientProvider[state.AWSAutoScaling] = aws_asg.NewAWSASGProvider(b.logger, b.eventChan)
	}
	b.logger.Debug().Msg("successfully setup AWS AutoScaling provider")

	if cfg.Provider.NoOp {
		b.clientProvider[state.NoOpClientProvider] = noop.NewNoOpProvider(b.logger, b.eventChan)
		b.logger.Debug().Msg("successfully setup notify log provider")
	}

	// Start the event handler.
	go b.eventUpdateHandler()

	return &b
}

// OKToScale satisfies the OKToScale function on the Scale interface.
func (b *Backend) OKToScale(req *state.ScalingRequest) (int, error) {
	// Create a temporary logger so that every log line includes the targeted class.
	logger := helper.LoggerWithNodeClassContext(b.logger, req.Policy.Class)

	logger.Info().
		Object("request", req).
		Msg("performing scaling precondition checks")

	// Perform an initial check to make sure the policy is enabled.
	if !req.Policy.Enabled {
		logger.Warn().Err(errScalingPolicyDisabled).Msg(scalingPreconditionCheckFailedMsg)
		return http.StatusUnprocessableEntity, errScalingPolicyDisabled
	}

	// Check that the node provider is configured for use.
	_, ok := b.clientProvider[req.Policy.Provider]
	if !ok {
		logger.Warn().Err(errScalingProviderNotFound).Msg(scalingPreconditionCheckFailedMsg)
		return http.StatusUnprocessableEntity, errScalingProviderNotFound
	}

	// Check there are actually nodes within the class which has received the request to scale.
	n := b.resourceHandler.GetNodesOfClass(req.Policy.Class)
	if n == nil || len(n) < 1 {
		logger.Warn().Err(errNoNodesFoundInClass).Msg(scalingPreconditionCheckFailedMsg)
		return http.StatusUnprocessableEntity, errNoNodesFoundInClass
	}

	// Check the new count does not break any thresholds.
	code, err := b.checkNewCount(req.Policy, req.Direction)
	if err != nil {
		logger.Warn().Err(err).Msg(scalingPreconditionCheckFailedMsg)
	}
	return code, err
}

// InvokeScaling satisfies the InvokeScaling function on the Scale interface.
func (b *Backend) InvokeScaling(req *state.ScalingRequest) {
	// Create a temporary logger so that every log line includes the targeted class.
	logger := helper.LoggerWithNodeClassContext(b.logger, req.Policy.Class)

	logger.Info().
		Object("request", req).
		Msg("performing scaling activity")

	// Write the initial event and state entry. If this fails, we do not continue.
	if err := b.scaleState.WriteRequest(req); err != nil {
		logger.Error().Err(err).Msg("failed to write initial state entry")
		return
	}
	err := b.invokeScaling(req)

	// Log the outcome of the scaling activity.
	if err != nil {
		b.logger.Error().Err(err).Object("request", req).Msg("scaling activity ended in failure")
	} else {
		b.logger.Info().Object("request", req).Msg("scaling activity ended successfully")
	}

	// Send the final activity update detailing the end state.
	b.eventChan <- &state.EventMessage{
		ID:        req.ID,
		Timestamp: helper.GenerateEventTimestamp(),
		Source:    eventSourceChemtrail,
		Error:     err,
	}
}

func (b *Backend) invokeScaling(req *state.ScalingRequest) error {
	switch req.Direction {
	case state.ScaleDirectionOut:
		return b.clientProvider[req.Policy.Provider].ScaleOut(req)

	case state.ScaleDirectionIn:
		// If we are scaling in, we need to discover the node we will target.
		node := b.resourceHandler.GetLeastAllocatedNodeInClass(req.Policy.Class)
		if node == nil {
			return errors.New("failed to discover least allocated node in class")
		}
		req.TargetNodeID = node.ID

		// If we are using the NoOp provider, we should not remove the node from the cluster.
		if req.Policy.Provider != state.NoOpClientProvider {
			if err := b.removeNodeFromCluster(req.TargetNodeID, req.ID); err != nil {
				return err
			}
		}

		target, err := b.identifyProviderTarget(req)
		if err != nil {
			return err
		}
		return b.clientProvider[req.Policy.Provider].ScaleIn(req, target)

	default:
		return errors.Errorf("unsupported scaling direction for invoke: %s", req.Direction.String())
	}
}

func (b *Backend) identifyProviderTarget(req *state.ScalingRequest) (string, error) {
	switch req.Policy.Provider {
	case state.AWSAutoScaling:
		return b.nodeIDToAWSInstanceID(req.TargetNodeID, req.ID)
	case state.NoOpClientProvider:
		return req.TargetNodeID, nil
	default:
		return "", errors.Errorf("unsupported provider: %s", req.Policy.Provider.String())
	}
}
