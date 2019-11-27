package scale

import (
	"github.com/jrasell/chemtrail/pkg/state"
)

const (
	// eventSourceNomad is the source to use when events are as a result of interactions with a
	// Nomad cluster.
	eventSourceNomad = "nomad"

	// eventSourceChemtrail is a special event source which originates from Chemtrail itself. This
	// is typically used to denote the start of end of a scaling activity.
	eventSourceChemtrail = "chemtrail"

	// eventMessageSuccess is the message used when a scaling activity has finished and has been
	// successful.
	eventMessageSuccess = "scaling activity has successfully completed"

	// eventMessageFailure is the message used when a scaling activity has reached a terminal
	// failure and can no longer continue.
	eventMessageFailure = "scaling activity has reached terminal failure"
)

// eventUpdateHandler is responsible for listing to the event channel and writing updates to the
// storage backend. Having central control over this is useful as events can come from a number of
// sources.
func (b *Backend) eventUpdateHandler() {
	for {
		select {
		case update := <-b.eventChan:
			var err error

			switch update.Source {
			case eventSourceChemtrail:
				err = b.handleChemtrailUpdate(update)
			default:
				err = b.handleBackendUpdate(update)
			}

			// In the event of an error log the error and the event. Writes to the state backend
			// will not stop a scaling activity so at least operators can check in the logs for
			// lost messages.
			if err != nil {
				b.logger.Error().
					Object("event", update).
					Err(err).Msg("failed to add scaling activity update")
				return
			}
			b.logger.Debug().Object("event", update).Msg("successfully stored scaling activity update")
		}
	}
}

// handleChemtrailUpdate is responsible for handling the Chemtrail update, which denotes the final
// update of a scaling activity and thus the end state. Depending whether an error is present in
// the message determines the final status and message. Generic messages are used here as it is
// expected that any source specific errors or messages are written to state before we trigger our
// end state.
func (b *Backend) handleChemtrailUpdate(msg *state.EventMessage) error {
	stateUpdate := state.ScalingUpdate{
		ID: msg.ID,
		Detail: state.Event{
			Timestamp: msg.Timestamp,
			Source:    eventSourceChemtrail,
		},
	}

	switch msg.Error {
	case nil:
		stateUpdate.Status = state.ScaleStatusCompleted
		stateUpdate.Detail.Message = eventMessageSuccess
	default:
		stateUpdate.Status = state.ScaleStatusFailed
		stateUpdate.Detail.Message = eventMessageFailure
	}

	return b.scaleState.WriteRequestEvent(&stateUpdate)
}

func (b *Backend) handleBackendUpdate(msg *state.EventMessage) error {
	stateUpdate := state.ScalingUpdate{
		ID:     msg.ID,
		Status: state.ScaleStatusInProgress,
		Detail: state.Event{
			Timestamp: msg.Timestamp,
			Message:   msg.Message,
			Source:    msg.Source,
		},
	}
	return b.scaleState.WriteRequestEvent(&stateUpdate)
}
