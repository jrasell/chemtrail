package memory

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/state"
)

type ScaleBackend struct {
	events map[uuid.UUID]*state.ScalingActivity
	l      sync.RWMutex
}

func NewScaleStateBackend() state.ScaleBackend {
	return &ScaleBackend{
		events: make(map[uuid.UUID]*state.ScalingActivity),
	}
}

// GetScalingActivities satisfies the GetScalingActivities function on the state.ScaleBackend
// interface.
func (s *ScaleBackend) GetScalingActivities() (map[uuid.UUID]*state.ScalingActivity, error) {
	s.l.RLock()
	events := s.events
	s.l.RUnlock()
	return events, nil
}

// GetScalingActivity satisfies the GetScalingActivity function on the state.ScaleBackend
// interface.
func (s *ScaleBackend) GetScalingActivity(id uuid.UUID) (*state.ScalingActivity, error) {
	s.l.RLock()
	event := s.events[id]
	s.l.RUnlock()
	return event, nil
}

// RunStateGarbageCollection satisfies the RunStateGarbageCollection function on the
// state.ScaleBackend interface.
func (s *ScaleBackend) RunStateGarbageCollection() {
	s.l.Lock()

	threshold := helper.GenerateEventTimestamp() - state.GarbageCollectionThreshold

	for id, activity := range s.events {
		if activity.LastUpdate < threshold {
			delete(s.events, id)
		}
	}
	s.l.Unlock()
}

// WriteRequestEvent satisfies the WriteRequestEvent function on the state.ScaleBackend interface.
func (s *ScaleBackend) WriteRequestEvent(message *state.ScalingUpdate) error {
	s.l.RLock()
	event := s.events[message.ID]
	s.l.RUnlock()

	detail := state.Event{
		Timestamp: message.Detail.Timestamp,
		Message:   message.Detail.Message,
		Source:    message.Detail.Source,
	}

	event.Events = append(event.Events, detail)
	event.LastUpdate = message.Detail.Timestamp

	if message.Status != state.ScaleStatusInProgress {
		event.Status = message.Status
	}

	s.l.Lock()
	s.events[message.ID] = event
	s.l.Unlock()

	return nil
}

// WriteRequest satisfies the WriteRequest function on the state.ScaleBackend interface.
func (s *ScaleBackend) WriteRequest(req *state.ScalingRequest) error {
	ts := helper.GenerateEventTimestamp()

	entry := state.ScalingActivity{
		Events: []state.Event{{
			Timestamp: ts,
			Message:   state.ScaleStartMessage,
			Source:    state.ScaleChemtrailSource,
		}},
		Direction:   req.Direction,
		LastUpdate:  ts,
		Status:      state.ScaleStatusStarted,
		Provider:    req.Policy.Provider,
		ProviderCfg: req.Policy.ProviderConfig,
	}

	s.l.Lock()
	s.events[req.ID] = &entry
	s.l.Unlock()

	return nil
}
