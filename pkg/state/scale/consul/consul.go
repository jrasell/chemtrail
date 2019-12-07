package consul

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// baseEventKVPath is Consul path suffix added to the CLI param which identifies where scaling
// state is stored.
const baseEventKVPath = "state/events/"

// ScaleBackend is the Consul implementation of the state.ScaleBackend interface.
type ScaleBackend struct {
	eventPath   string
	gcThreshold int64
	kv          *api.KV
	logger      zerolog.Logger
}

// NewPolicyBackend returns the Consul implementation of the state.ScaleBackend interface.
func NewScaleBackend(log zerolog.Logger, path string, client *api.Client) state.ScaleBackend {
	return &ScaleBackend{
		eventPath:   path + baseEventKVPath,
		gcThreshold: state.GarbageCollectionThreshold,
		kv:          client.KV(),
		logger:      log,
	}
}

// GetScalingActivities satisfies the GetScalingActivities function on the state.ScaleBackend
// interface.
func (s *ScaleBackend) GetScalingActivities() (map[uuid.UUID]*state.ScalingActivity, error) {
	kv, _, err := s.kv.List(s.eventPath, nil)
	if err != nil {
		return nil, err
	}

	out := make(map[uuid.UUID]*state.ScalingActivity)

	if kv == nil {
		return out, nil
	}

	for i := range kv {
		activity := &state.ScalingActivity{}

		if err := json.Unmarshal(kv[i].Value, activity); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
		}

		keySplit := strings.Split(kv[i].Key, "/")

		id, err := uuid.FromString(keySplit[len(keySplit)-1])
		if err != nil {
			return nil, errors.Wrap(err, "failed to get UUID from string")
		}

		out[id] = activity
	}

	return out, nil
}

// GetScalingActivity satisfies the GetScalingActivity function on the state.ScaleBackend
// interface.
func (s *ScaleBackend) GetScalingActivity(id uuid.UUID) (*state.ScalingActivity, error) {
	kv, _, err := s.kv.Get(s.eventPath+id.String(), nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	out := &state.ScalingActivity{}

	if err := json.Unmarshal(kv.Value, out); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
	}

	return out, nil
}

// RunStateGarbageCollection satisfies the RunStateGarbageCollection function on the
// state.ScaleBackend interface.
func (s *ScaleBackend) RunStateGarbageCollection() {

	kv, _, err := s.kv.List(s.eventPath, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("GC failed to list events in Consul backend")
		return
	}

	if kv == nil {
		return
	}

	gc := helper.GenerateEventTimestamp() - s.gcThreshold

	for i := range kv {

		ss := &state.ScalingActivity{}

		if err := json.Unmarshal(kv[i].Value, ss); err != nil {
			s.logger.Error().Err(err).Msg("GC failed to unmarshal event for inspection")
			continue
		}

		switch ss.Status {
		case state.ScaleStatusCompleted, state.ScaleStatusFailed:
			if ss.LastUpdate < gc {
				// Unlike the in-memory, we currently delete keys which have passed the expiration
				// threshold. Delete vs. re-create has not been benchmarked, but my initial opinion is
				// that delete will be more efficient and is at least easier for the MVP.
				if _, err := s.kv.Delete(kv[i].Key, nil); err != nil {
					s.logger.Error().
						Str("key", kv[i].Key).
						Err(err).
						Msg("GC failed to delete stale event in Consul backend")
				}
			}
		default:
			continue
		}
	}
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

	marshal, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	pair := &api.KVPair{
		Key:   s.eventPath + req.ID.String(),
		Value: marshal,
	}

	_, err = s.kv.Put(pair, nil)
	return err
}

// WriteRequestEvent satisfies the WriteRequestEvent function on the state.ScaleBackend interface.
func (s *ScaleBackend) WriteRequestEvent(message *state.ScalingUpdate) error {
	kv, _, err := s.kv.Get(s.eventPath+message.ID.String(), nil)
	if err != nil {
		return err
	}

	// Adding an event to an activity requires the initial state be written to Consul. In the
	// situation where no KV is found, this is an error and should be reported as such.
	if kv == nil {
		return errors.New("scaling activity not found in Consul backend")
	}

	event := &state.ScalingActivity{}

	if err := json.Unmarshal(kv.Value, event); err != nil {
		return errors.Wrap(err, "failed to unmarshal Consul KV value")
	}

	// Build the additional activity information and add this to the list of events whilst updating
	// the timestamp and status.
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

	marshal, err := json.Marshal(event)
	if err != nil {
		return err
	}

	pair := &api.KVPair{
		Key:   s.eventPath + message.ID.String(),
		Value: marshal,
	}

	_, err = s.kv.Put(pair, nil)
	return err
}
