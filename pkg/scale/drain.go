package scale

import (
	"context"
	"strings"
	"time"

	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/state"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/nomad/api"
)

const (
	drainDeadlineMinutes = 5
)

func (b *Backend) removeNodeFromCluster(nodeID string, scaleID uuid.UUID) error {
	b.logger.Info().
		Str("node-id", nodeID).
		Msg("removing node from Nomad cluster")

	drainSpec := api.DrainSpec{Deadline: drainDeadlineMinutes * time.Minute}

	resp, err := b.nomad.Client.Nodes().UpdateDrain(nodeID, &drainSpec, false, nil)
	if err != nil {
		return err
	}
	b.monitorNodeDrain(nodeID, scaleID, resp.LastIndex)

	return nil
}

func (b *Backend) monitorNodeDrain(nodeID string, scaleID uuid.UUID, index uint64) {
	for msg := range b.nomad.Client.Nodes().MonitorDrain(context.Background(), nodeID, index, false) {
		b.eventChan <- &state.EventMessage{
			ID:        scaleID,
			Timestamp: helper.GenerateEventTimestamp(),
			Source:    eventSourceNomad,
			Message:   strings.ToLower(msg.Message),
		}
	}
}
