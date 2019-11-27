package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handler_calculateAllocatedPercentageStats(t *testing.T) {
	handler := &handler{nodeManager: &updateHandler{}}

	testCases := []struct {
		inputStats     *resourceStats
		expectedOutput *AllocatedStats
		name           string
	}{
		{
			inputStats: &resourceStats{
				allocatableResources: &resources{cpu: 1000, memory: 1000},
				allocatedResources:   &resources{cpu: 0, memory: 0},
			},
			expectedOutput: &AllocatedStats{CPU: 0, Memory: 0},
			name:           "zero allocated memory and cpu",
		},
		{
			inputStats: &resourceStats{
				allocatableResources: &resources{cpu: 1000, memory: 1000},
				allocatedResources:   &resources{cpu: 1000, memory: 1000},
			},
			expectedOutput: &AllocatedStats{CPU: 100, Memory: 100},
			name:           "100% allocated memory and cpu",
		},
		{
			inputStats: &resourceStats{
				allocatableResources: &resources{cpu: 100, memory: 100},
				allocatedResources:   &resources{cpu: 10, memory: 10},
			},
			expectedOutput: &AllocatedStats{CPU: 10, Memory: 10},
			name:           "10% allocated memory and cpu",
		},
	}

	for _, tc := range testCases {
		actualOutput := handler.calculateAllocatedPercentageStats(tc.inputStats)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}
