package scale

import (
	"net/http"

	"github.com/jrasell/chemtrail/pkg/state"
)

func (b *Backend) checkNewCount(policy *state.ClientScalingPolicy, dir state.ScaleDirection) (int, error) {
	nodes := b.resourceHandler.GetNodesOfClass(policy.Class)

	switch dir {
	case state.ScaleDirectionIn:
		if (len(nodes) - policy.ScaleInCount) < policy.MinCount {
			return http.StatusPreconditionFailed, errScalingInCountCheckFailed
		}
	case state.ScaleDirectionOut:
		if (len(nodes) + policy.ScaleOutCount) > policy.MaxCount {
			return http.StatusPreconditionFailed, errScalingOutCountCheckFailed
		}
	}

	return http.StatusOK, nil
}
