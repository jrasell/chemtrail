package awsasg

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/state"
)

// event in a type of AWS AutoScaling interaction, which helps dictate to logs and events recorded
// in the Chemtrail server.
type event string

const (
	eventTypeDesc      event = "describe"
	eventTypeDetach    event = "detach"
	eventTypeUpdate    event = "update"
	eventTypeTerminate event = "terminate"
)

// handleEvent is used to managed AWS AutoScaling provider events in a generic manner.
func (a *ClientProvider) handleEvent(e event, err error, resource *string, id uuid.UUID) {

	// Build the base event message with params which are common.
	msg := state.EventMessage{ID: id, Timestamp: helper.GenerateEventTimestamp(), Source: a.Name()}

	switch err {
	case nil:
		a.handleEventSuccess(e, &msg, resource)
	default:
		a.handleEventError(e, &msg, err, resource)
	}
}

func (a *ClientProvider) handleEventError(e event, msg *state.EventMessage, err error, resource *string) {
	var msgString string

	switch e {
	case eventTypeDesc:
		msgString = "failed to describe AWS AutoScaling group"
	case eventTypeDetach:
		msgString = fmt.Sprintf("failed to detach instance %s from AWS AutoScaling group", *resource)
	case eventTypeUpdate:
		msgString = "failed to update count of AWS AutoScaling group"
	case eventTypeTerminate:
		msgString = fmt.Sprintf("failed to terminate AWS EC2 instance %s", *resource)
	default:
	}

	// Log the message to include the provide error message.
	a.log.Error().Err(err).Msg(msgString)

	// Update and send the event message to store in the backend.
	msg.Message = msgString
	a.eventChan <- msg
}

func (a *ClientProvider) handleEventSuccess(e event, msg *state.EventMessage, resource *string) {
	var msgString string

	switch e {
	case eventTypeDesc:
		msgString = "successfully described AWS AutoScaling group"
	case eventTypeDetach:
		msgString = fmt.Sprintf("successfully detched instance %s from AWS AutoScaling group", *resource)
	case eventTypeUpdate:
		msgString = "successfully updated count of AWS AutoScaling group"
	case eventTypeTerminate:
		msgString = fmt.Sprintf("successfully terminated AWS EC2 instance %s", *resource)
	default:
	}

	// Log the message to info including the call that was made successfully.
	a.log.Info().Msg(msgString)

	// Update and send the event message to store in the backend.
	msg.Message = msgString
	a.eventChan <- msg
}
