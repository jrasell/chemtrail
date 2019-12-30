package noop

import (
	"github.com/gofrs/uuid"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/scale/provider"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/rs/zerolog"
)

// ClientProvider implements the provider.ClientProvider interface.
type ClientProvider struct {
	log       zerolog.Logger
	eventChan chan *state.EventMessage
}

// NewNoOpProvider creates a new log notification client scaling provider.
func NewNoOpProvider(log zerolog.Logger, eventChan chan *state.EventMessage) provider.ClientProvider {
	return &ClientProvider{
		log:       log.With().Str("provider", state.NoOpClientProvider.String()).Logger(),
		eventChan: eventChan,
	}
}

// Name satisfies the provider.ClientProvider Name interface function.
func (a *ClientProvider) Name() string { return state.NoOpClientProvider.String() }

// ScaleIn satisfies the provider.ClientProvider ScaleIn interface function.
func (a *ClientProvider) ScaleIn(req *state.ScalingRequest, _ string) error {
	return a.notifyWrapper(req)
}

// ScaleOut satisfies the provider.ClientProvider ScaleOut interface function.
func (a *ClientProvider) ScaleOut(req *state.ScalingRequest) error {
	return a.notifyWrapper(req)
}

func (a *ClientProvider) notifyWrapper(req *state.ScalingRequest) error {
	a.sendScaleLogMessage(req)
	a.sendEvent(req.ID)
	return nil
}

// sendScaleLogMessage is responsible for sending the notification message containing the scaling
// request information.
func (a *ClientProvider) sendScaleLogMessage(req *state.ScalingRequest) {
	a.log.Info().
		Str("id", req.ID.String()).
		Str("direction", req.Direction.String()).
		Str("target-node", req.TargetNodeID).
		Object("policy", req.Policy).
		Msg("no-op log notification of scaling activity")
}

// sendEvent sends a scaling event to be stored tracking the successful log notification. This
// allows Chemtrail to run essentially in noop mode, while still going through the entire cycle of
// a scaling action.
func (a *ClientProvider) sendEvent(id uuid.UUID) {
	a.eventChan <- &state.EventMessage{
		ID:        id,
		Timestamp: helper.GenerateEventTimestamp(),
		Source:    a.Name(),
		Message:   "successfully triggered log notification",
	}
}
