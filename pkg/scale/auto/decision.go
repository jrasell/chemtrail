package auto

import (
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/rs/zerolog"
)

type decision struct {
	direction state.ScaleDirection
}

func (s *Scale) performPolicyChecks(log zerolog.Logger, pol *state.ClientScalingPolicy) (*decision, error) {

	// Gather the allocated resource stats for the node class.
	allocStats, err := s.resourceHandler.GetClassResourceAllocation(pol.Class)
	if err != nil {
		return nil, err
	}

	// Create a decision mapping. This allows us to track the decisions made by the various checks
	// and store information as desired to explain what is happening. This is currently a boolean
	// until the desired data is better understood.
	classDecision := make(map[state.ScaleDirection]bool)

	for name, check := range pol.Checks {

		if !check.Enabled {
			log.Debug().
				Str("check-name", name).
				Msg("scaling policy check administratively disabled")
			continue
		}

		var actual float64

		switch check.Resource {
		case state.ScaleResourceCPU:
			actual = allocStats.CPU
		case state.ScaleResourceMemory:
			actual = allocStats.Memory
		}
		log.Debug().
			Str("check-name", name).
			Str("check-resource", check.Resource.String()).
			Float64("check-resource-threshold", check.ComparisonPercentage).
			Float64("check-resource-actual", actual).
			Str("check-resource-comparison", check.ComparisonOperator.String()).
			Msg("performing scaling policy check analysis")

		// Perform the actual check comparison and blindly add a class decision entry. This will
		// change once we track decision metadata.
		checkDecision := s.performPolicyCheck(check, actual, check.ComparisonPercentage)

		if checkDecision.direction != state.ScaleDirectionNone {
			classDecision[checkDecision.direction] = true
		}
	}
	return s.buildSingleDecision(classDecision), nil
}

func (s *Scale) buildSingleDecision(decisions map[state.ScaleDirection]bool) *decision {
	if len(decisions) == 0 {
		return nil
	}

	// ScaleDirectionOut should always trump in.
	if decisions[state.ScaleDirectionOut] && decisions[state.ScaleDirectionIn] {
		return &decision{direction: state.ScaleDirectionOut}
	}

	if decisions[state.ScaleDirectionOut] {
		return &decision{direction: state.ScaleDirectionOut}
	}
	if decisions[state.ScaleDirectionIn] {
		return &decision{direction: state.ScaleDirectionIn}
	}
	return nil
}

func (s *Scale) performPolicyCheck(check *state.PolicyCheck, actual, threshold float64) *decision {
	var dec *decision

	switch check.ComparisonOperator {
	case state.ComparisonLessThan:
		dec = performLessThanCheck(actual, threshold, check.Action)
	case state.ComparisonGreaterThan:
		dec = performGreaterThanCheck(actual, threshold, check.Action)
	}
	return dec
}

func performGreaterThanCheck(actual, threshold float64, action state.ComparisonAction) *decision {
	if actual > threshold {
		switch action {
		case state.ActionScaleIn:
			return &decision{direction: state.ScaleDirectionIn}
		case state.ActionScaleOut:
			return &decision{direction: state.ScaleDirectionOut}
		default:
		}
	}
	return &decision{direction: state.ScaleDirectionNone}
}

func performLessThanCheck(actual, threshold float64, action state.ComparisonAction) *decision {
	if actual < threshold {
		switch action {
		case state.ActionScaleIn:
			return &decision{direction: state.ScaleDirectionIn}
		case state.ActionScaleOut:
			return &decision{direction: state.ScaleDirectionOut}
		default:
		}
	}
	return &decision{direction: state.ScaleDirectionNone}
}
