package auto

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jrasell/chemtrail/pkg/client"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/scale"
	"github.com/jrasell/chemtrail/pkg/scale/resource"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog"
)

type Scale struct {
	threads  int
	interval int
	logger   zerolog.Logger
	nomad    *client.Nomad

	policyBackend   state.PolicyBackend
	scaler          scale.Scale
	resourceHandler resource.Handler
	pool            *ants.PoolWithFunc
	inProgress      bool

	// isRunning is used to track whether the autoscaler loop is being run. This helps determine
	// whether stop should be called.
	isRunning bool

	// doneChan is used to stop the autoscaling execution.
	doneChan chan struct{}
}

func NewAutoScaler(cfg *Config) (*Scale, error) {
	s := Scale{
		threads:         cfg.Threads,
		interval:        cfg.Interval,
		logger:          cfg.Logger,
		nomad:           cfg.Nomad,
		policyBackend:   cfg.Policy,
		scaler:          cfg.Scale,
		resourceHandler: cfg.Resource,
		doneChan:        make(chan struct{}),
	}

	pool, err := s.createWorkerPool()
	if err != nil {
		return nil, err
	}
	s.pool = pool

	return &s, nil
}

func (s *Scale) Run() {
	s.logger.Info().Msg("starting Chemtrail internal auto-scaling engine")

	// Track that the autoscaler is actively running.
	s.isRunning = true

	t := time.NewTicker(time.Second * time.Duration(s.interval))
	defer t.Stop()

	for {
		select {
		case <-t.C:
			s.logger.Info().Msg("triggering new autoscaler evaluation run")

			// Check whether a previous scaling loop is in progress, and if it is we should skip
			// this round. This avoids putting more pressure on a system which may be under load
			// causing slow API responses.
			if s.inProgress {
				s.logger.Info().Msg("scaling run in progress, skipping new assessment")
				break
			}
			s.setScalingInProgressTrue()

			policies, err := s.policyBackend.GetPolicies()
			if err != nil {
				s.logger.Error().Err(err).Msg("autoscaler unable to get scaling policies")
				s.setScalingInProgressFalse()
				break
			}

			if len(policies) == 0 {
				s.logger.Debug().Msg("no scaling policies found in storage backend for autoscaler to iterate")
				s.setScalingInProgressFalse()
				break
			}

			for _, policy := range policies {

				// Check whether the policy is enabled.
				if !policy.Enabled {
					continue
				}

				if err := s.pool.Invoke(policy); err != nil {
					s.logger.Error().Err(err).Msg("failed to invoke autoscaling worker thread")
				}
			}
			s.setScalingInProgressFalse()

		case <-s.doneChan:
			s.isRunning = false
			return
		}
	}
}

// Stop is used to gracefully stop the autoscaling workers.
func (s *Scale) Stop() {

	// Inform sub-process to exit.
	close(s.doneChan)

	for {
		if !s.isRunning && !s.inProgress {
			s.pool.Release()
			s.logger.Info().Msg("successfully drained autoscaler worker pool")
			return
		}
		s.logger.Debug().Msg("autoscaler still has in-flight workers, will continue to check")
		time.Sleep(1 * time.Second)
	}
}

func (s *Scale) setScalingInProgressTrue()  { s.inProgress = true }
func (s *Scale) setScalingInProgressFalse() { s.inProgress = false }

// createWorkerPool is responsible for building the ants goroutine worker pool with the number of
// threads controlled by the operator configured value.
func (s *Scale) createWorkerPool() (*ants.PoolWithFunc, error) {
	return ants.NewPoolWithFunc(s.threads, s.workerPoolFunc(), ants.WithExpiryDuration(60*time.Second))
}

func (s *Scale) workerPoolFunc() func(payload interface{}) {
	return func(payload interface{}) {

		// If this thread starts after the autoscaler has been asked to shutdown, exit. Otherwise
		// perform the work.
		select {
		case <-s.doneChan:
			s.logger.Debug().Msg("exiting autoscaling thread as a result of shutdown request")
			return
		default:
		}

		req, ok := payload.(*state.ClientScalingPolicy)
		if !ok {
			s.logger.Error().Msg("autoscaler worker pool received unexpected payload type")
			return
		}

		// Create a temporary logger so that every log line includes the targeted class.
		logger := helper.LoggerWithNodeClassContext(s.logger, req.Class)

		scalingDecision, err := s.performPolicyChecks(logger, req)
		if err != nil {
			logger.Error().Err(err).Msg("unable to perform node class scaling decision")
			return
		}
		if scalingDecision == nil {
			logger.Debug().Msg("no scaling action required")
			return
		}

		// Generate a UUID which is used as the scaling identifier. If we can't generate this then
		// we exit.
		id, err := uuid.NewV4()
		if err != nil {
			logger.Error().Err(err).Msg("failed to generate scaling ID")
			return
		}

		scalingReq := state.ScalingRequest{ID: id, Direction: scalingDecision.direction, Policy: req}
		if _, err := s.scaler.OKToScale(&scalingReq); err != nil {
			logger.Info().Str("reason", err.Error()).Msg("autoscaling activity not allowed to continue")
			return
		}
		s.scaler.InvokeScaling(&scalingReq)
	}
}

// IsRunning is used to determine if the autoscaler loop is running.
func (s *Scale) IsRunning() bool { return s.isRunning }
