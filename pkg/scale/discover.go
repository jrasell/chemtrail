package scale

import (
	"github.com/gofrs/uuid"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/pkg/errors"
)

const (
	awsInstanceIDAttr = "unique.platform.aws.instance-id"
)

const (
	eventMsgFailedNodeInfo          = "failed to call Nomad node info API"
	eventMsgAWSInstanceAttrNotFound = "aws instance-id not found within attributes"
)

func (b *Backend) nodeIDToAWSInstanceID(nodeID string, id uuid.UUID) (string, error) {
	node, _, err := b.nomad.Client.Nodes().Info(nodeID, nil)
	if err != nil {
		b.eventChan <- &state.EventMessage{
			ID:        id,
			Timestamp: helper.GenerateEventTimestamp(),
			Source:    eventSourceNomad,
			Message:   eventMsgFailedNodeInfo,
		}
		return "", err
	}

	val, ok := node.Attributes[awsInstanceIDAttr]
	if !ok {
		b.eventChan <- &state.EventMessage{
			ID:        id,
			Timestamp: helper.GenerateEventTimestamp(),
			Source:    eventSourceNomad,
			Message:   eventMsgAWSInstanceAttrNotFound,
		}
		return "", errors.New(eventMsgAWSInstanceAttrNotFound)
	}
	return val, nil
}
