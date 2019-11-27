package auto

import (
	"testing"

	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/stretchr/testify/assert"
)

func Test_performGreaterThanCheck(t *testing.T) {
	testCases := []struct {
		inputActual    float64
		inputThreshold float64
		inputName      string
		inputAction    state.ComparisonAction
		expectedOutput *decision
		name           string
	}{
		{
			inputActual:    90,
			inputThreshold: 91,
			inputName:      "test_cpu_check",
			inputAction:    state.ActionScaleOut,
			expectedOutput: &decision{direction: state.ScaleDirectionNone},
			name:           "expected scale none result",
		},
		{
			inputActual:    90.001,
			inputThreshold: 90,
			inputName:      "test_check",
			inputAction:    state.ActionScaleOut,
			expectedOutput: &decision{direction: state.ScaleDirectionOut},
			name:           "expected scale out result with metric decision",
		},
		{
			inputActual:    111001.1,
			inputThreshold: 111001.01,
			inputName:      "test_check",
			inputAction:    state.ActionScaleIn,
			expectedOutput: &decision{direction: state.ScaleDirectionIn},
			name:           "expected scale in result with metric decision",
		},
	}

	for _, tc := range testCases {
		actualOutput := performGreaterThanCheck(tc.inputActual, tc.inputThreshold, tc.inputAction)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func Test_performLessThanCheck(t *testing.T) {
	testCases := []struct {
		inputActual    float64
		inputThreshold float64
		inputName      string
		inputAction    state.ComparisonAction
		expectedOutput *decision
		name           string
	}{
		{
			inputActual:    91,
			inputThreshold: 90,
			inputName:      "test_cpu_check",
			inputAction:    state.ActionScaleOut,
			expectedOutput: &decision{direction: state.ScaleDirectionNone},
			name:           "expected scale none result",
		},
		{
			inputActual:    90,
			inputThreshold: 90.001,
			inputName:      "test_check",
			inputAction:    state.ActionScaleOut,
			expectedOutput: &decision{direction: state.ScaleDirectionOut},
			name:           "expected scale out result with metric decision",
		},
		{
			inputActual:    111001.01,
			inputThreshold: 111001.1,
			inputName:      "test_check",
			inputAction:    state.ActionScaleIn,
			expectedOutput: &decision{direction: state.ScaleDirectionIn},
			name:           "expected scale in result with metric decision",
		},
	}

	for _, tc := range testCases {
		actualOutput := performLessThanCheck(tc.inputActual, tc.inputThreshold, tc.inputAction)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}
