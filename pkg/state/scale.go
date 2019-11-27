package state

import (
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
)

// ScaleBackend is the interface which storage providers must implement in order to store scaling
// state.
type ScaleBackend interface {

	// GetScalingActivities returns all the currently stored scaling activates from the storage
	// backend. The map is keyed by the scaling ID.
	GetScalingActivities() (map[uuid.UUID]*ScalingActivity, error)

	// GetScalingActivity attempts to retrieve a single scaling activity from the store based on
	// the scaling ID.
	GetScalingActivity(id uuid.UUID) (*ScalingActivity, error)

	// RunStateGarbageCollection triggers a cleanup of the scaling state, removing all entries
	// which are older than the configured threshold.
	RunStateGarbageCollection()

	// WriteRequest is used to write an initial scaling request to the store. This should be called
	// immediately after an activity is requested.
	WriteRequest(req *ScalingRequest) error

	// WriteRequestEvent updates the stored request details with the passed event. This should not
	// overwrite the stored data, but append to the list of stored events.
	WriteRequestEvent(message *ScalingUpdate) error
}

const (
	// GarbageCollectionThreshold is a time in nanoseconds. This is used to determine whether or
	// not scaling events should be garbage collected. This is 24hrs.
	GarbageCollectionThreshold = 172800000000000
)

// ScalingActivity represents an individual scaling operation and is designed to take the overall
// configuration, as well as events during the operation.
type ScalingActivity struct {

	// Events is a list of events which occurred during the scaling operation and provides insight
	// into the operations and actions conducted.
	Events []Event

	// ScaleDirection is the direction of scaling.
	Direction ScaleDirection

	// LastUpdate is the UnixNano timestamp of the last update to the scaling operation and if the
	// status is terminal, indicates the time of the last action of the scaling event.
	LastUpdate int64

	// ScaleStatus indicates the current status of the scaling operation.
	Status ScaleStatus

	// Provider indicated the backend node provider used for this scaling operation.
	Provider ClientProvider

	// ProviderCfg is a key value map containing any runtime specific information for the node
	// provider. This should not include credentials, but instead items such as ASG name,
	// instanceIDs or IP addresses.
	ProviderCfg map[string]string
}

type ScalingUpdate struct {
	ID     uuid.UUID
	Status ScaleStatus
	Detail Event
}

func (su ScalingUpdate) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", su.ID.String()).
		Str("message", su.Detail.Message).
		Int64("timestamp", su.Detail.Timestamp).
		Str("status", su.Status.String()).
		Str("source", su.Detail.Source)
}

type Event struct {
	Timestamp int64
	Message   string
	Source    string
}

type ScalingRequest struct {
	ID           uuid.UUID
	Direction    ScaleDirection
	TargetNodeID string
	Policy       *ClientScalingPolicy
}

func (sr ScalingRequest) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", sr.ID.String()).Str("direction", sr.Direction.String())
}

// ScaleDirection describes the direction which a scaling activity should take.
type ScaleDirection string

// String is a helper method to return the string of the ScaleDirection.
func (sd ScaleDirection) String() string { return string(sd) }

const (
	// ScaleDirectionIn indicates that an activity should be undertaken to decrement the current
	// client count.
	ScaleDirectionIn ScaleDirection = "in"

	// ScaleDirectionOut indicates that an activity should be undertaken to increment the current client count.
	ScaleDirectionOut ScaleDirection = "out"

	// ScaleDirectionNone is an indication that no scaling action is required.
	ScaleDirectionNone ScaleDirection = "none"
)

// ScaleStatus describes the state of a scaling activity as well as the state an activity was in
// when an event was recorded.
type ScaleStatus string

// String is a helper method to return the string of the ScaleStatus.
func (s ScaleStatus) String() string { return string(s) }

const (
	// ScaleStatusStarted represents the initial event of a scaling operation. The source of an
	// event using this status should always be Chemtrail.
	ScaleStatusStarted ScaleStatus = "started"

	// ScaleStatusInProgress is a general purpose status which all downstream activity providers
	// can use to record scaling operation events.
	ScaleStatusInProgress ScaleStatus = "in-progress"

	// ScaleStatusCompleted is a terminally successful scaling status. The source of an event using
	// this status should always be Chemtrail.
	ScaleStatusCompleted ScaleStatus = "completed"

	// ScaleStatusFailed is a terminally unsuccessful scaling status. The source of an event using
	// this status should always be Chemtrail.
	ScaleStatusFailed ScaleStatus = "failed"
)

type EventMessage struct {
	ID        uuid.UUID
	Timestamp int64
	Source    string
	Message   string
	Error     error
}

func (em EventMessage) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", em.ID.String()).
		Int64("timestamp", em.Timestamp).
		Str("source", em.Source)

	if em.Error != nil {
		e.Err(em.Error)
	}
	if em.Message != "" {
		e.Str("message", em.Message)
	}
}
