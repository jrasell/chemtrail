package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScaleDirection_String(t *testing.T) {
	testCases := []struct {
		inputScaleDirection ScaleDirection
		expectedOutput      string
		name                string
	}{
		{
			inputScaleDirection: ScaleDirectionIn,
			expectedOutput:      "in",
		},
		{
			inputScaleDirection: ScaleDirectionOut,
			expectedOutput:      "out",
		},
		{
			inputScaleDirection: ScaleDirectionNone,
			expectedOutput:      "none",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputScaleDirection.String()
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func TestScaleStatus_String(t *testing.T) {
	testCases := []struct {
		inputScaleStatus ScaleStatus
		expectedOutput   string
		name             string
	}{
		{
			inputScaleStatus: ScaleStatusStarted,
			expectedOutput:   "started",
		},
		{
			inputScaleStatus: ScaleStatusInProgress,
			expectedOutput:   "in-progress",
		},
		{
			inputScaleStatus: ScaleStatusCompleted,
			expectedOutput:   "completed",
		},
		{
			inputScaleStatus: ScaleStatusFailed,
			expectedOutput:   "failed",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputScaleStatus.String()
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}
