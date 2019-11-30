package resource

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
)

func Test_updateHandler_handleAllocMessageTerminal(t *testing.T) {
	testCases := []struct {
		inputAlloc             *api.Allocation
		inputClass             string
		inputHandler           *updateHandler
		expectedNodePoolResult map[string]*classInfo
		name                   string
	}{
		{
			inputAlloc: &api.Allocation{
				ID:           "test-alloc",
				NodeID:       "test-node",
				ClientStatus: "dead",
				Resources:    &api.Resources{CPU: intToPointer(500), MemoryMB: intToPointer(256)},
			},
			inputClass: "test-class",
			inputHandler: &updateHandler{
				nodePool: map[string]*classInfo{
					"test-class": {
						class: "test-class",
						nodes: map[string]*nodeInfo{
							"test-node": {
								ID:    "test-node",
								class: "test-class",
								resourceStats: &resourceStats{
									allocatableResources: &resources{cpu: 5182, memory: 985},
									allocatedResources:   &resources{cpu: 500, memory: 256},
								},
							},
						},
						resourceStats: &resourceStats{
							allocatableResources: &resources{cpu: 5182, memory: 985},
							allocatedResources:   &resources{cpu: 500, memory: 256},
						},
						allocations: map[string]string{"test-alloc": "running"},
					},
				},
			},
			expectedNodePoolResult: map[string]*classInfo{
				"test-class": {
					class: "test-class",
					nodes: map[string]*nodeInfo{
						"test-node": {
							ID:    "test-node",
							class: "test-class",
							resourceStats: &resourceStats{
								allocatableResources: &resources{cpu: 5182, memory: 985},
								allocatedResources:   &resources{cpu: 0, memory: 0},
							},
						},
					},
					resourceStats: &resourceStats{
						allocatableResources: &resources{cpu: 5182, memory: 985},
						allocatedResources:   &resources{cpu: 0, memory: 0},
					},
					allocations: map[string]string{},
				},
			},
			name: "remove only alloc found within node class",
		},

		{
			inputAlloc: &api.Allocation{
				ID:           "removal-alloc",
				NodeID:       "test-node",
				ClientStatus: "dead",
				Resources:    &api.Resources{CPU: intToPointer(500), MemoryMB: intToPointer(256)},
			},
			inputClass: "test-class",
			inputHandler: &updateHandler{
				nodePool: map[string]*classInfo{
					"test-class": {
						class: "test-class",
						nodes: map[string]*nodeInfo{
							"test-node": {
								ID:    "test-node",
								class: "test-class",
								resourceStats: &resourceStats{
									allocatableResources: &resources{cpu: 5182, memory: 985},
									allocatedResources:   &resources{cpu: 1000, memory: 512},
								},
							},
						},
						resourceStats: &resourceStats{
							allocatableResources: &resources{cpu: 5182, memory: 985},
							allocatedResources:   &resources{cpu: 1000, memory: 512},
						},
						allocations: map[string]string{"existing-alloc": "running", "removal-alloc": "running"},
					},
				},
			},
			expectedNodePoolResult: map[string]*classInfo{
				"test-class": {
					class: "test-class",
					nodes: map[string]*nodeInfo{
						"test-node": {
							ID:    "test-node",
							class: "test-class",
							resourceStats: &resourceStats{
								allocatableResources: &resources{cpu: 5182, memory: 985},
								allocatedResources:   &resources{cpu: 500, memory: 256},
							},
						},
					},
					resourceStats: &resourceStats{
						allocatableResources: &resources{cpu: 5182, memory: 985},
						allocatedResources:   &resources{cpu: 500, memory: 256},
					},
					allocations: map[string]string{"existing-alloc": "running"},
				},
			},
			name: "remove alloc found within node class with other allocs",
		},
	}

	for _, tc := range testCases {
		tc.inputHandler.handleAllocMessageTerminal(tc.inputClass, tc.inputAlloc)
		assert.Equal(t, tc.expectedNodePoolResult, tc.inputHandler.nodePool, tc.name)
	}
}

func Test_updateHandler_handleAllocMessageRunning(t *testing.T) {
	testCases := []struct {
		inputAlloc             *api.Allocation
		inputClass             string
		inputHandler           *updateHandler
		expectedNodePoolResult map[string]*classInfo
		name                   string
	}{
		{
			inputAlloc: &api.Allocation{
				ID:           "test-alloc",
				NodeID:       "test-node",
				ClientStatus: "running",
				Resources:    &api.Resources{CPU: intToPointer(500), MemoryMB: intToPointer(256)},
			},
			inputClass: "test-class",
			inputHandler: &updateHandler{
				nodePool: map[string]*classInfo{
					"test-class": {
						class: "test-class",
						nodes: map[string]*nodeInfo{
							"test-node": {
								ID:    "test-node",
								class: "test-class",
								resourceStats: &resourceStats{
									allocatableResources: &resources{cpu: 5182, memory: 985},
									allocatedResources:   &resources{cpu: 0, memory: 0},
								},
							},
						},
						resourceStats: &resourceStats{
							allocatableResources: &resources{cpu: 5182, memory: 985},
							allocatedResources:   &resources{cpu: 0, memory: 0},
						},
						allocations: make(map[string]string),
					},
				},
			},
			expectedNodePoolResult: map[string]*classInfo{
				"test-class": {
					class: "test-class",
					nodes: map[string]*nodeInfo{
						"test-node": {
							ID:    "test-node",
							class: "test-class",
							resourceStats: &resourceStats{
								allocatableResources: &resources{cpu: 5182, memory: 985},
								allocatedResources:   &resources{cpu: 500, memory: 256},
							},
						},
					},
					resourceStats: &resourceStats{
						allocatableResources: &resources{cpu: 5182, memory: 985},
						allocatedResources:   &resources{cpu: 500, memory: 256},
					},
					allocations: map[string]string{"test-alloc": "running"},
				},
			},
			name: "first alloc to be found of node class",
		},

		{
			inputAlloc: &api.Allocation{
				ID:           "new-alloc",
				NodeID:       "test-node",
				ClientStatus: "running",
				Resources:    &api.Resources{CPU: intToPointer(500), MemoryMB: intToPointer(256)},
			},
			inputClass: "test-class",
			inputHandler: &updateHandler{
				nodePool: map[string]*classInfo{
					"test-class": {
						class: "test-class",
						nodes: map[string]*nodeInfo{
							"test-node": {
								ID:    "test-node",
								class: "test-class",
								resourceStats: &resourceStats{
									allocatableResources: &resources{cpu: 5182, memory: 985},
									allocatedResources:   &resources{cpu: 500, memory: 256},
								},
							},
						},
						resourceStats: &resourceStats{
							allocatableResources: &resources{cpu: 5182, memory: 985},
							allocatedResources:   &resources{cpu: 500, memory: 256},
						},
						allocations: map[string]string{"existing-alloc": "running"},
					},
				},
			},
			expectedNodePoolResult: map[string]*classInfo{
				"test-class": {
					class: "test-class",
					nodes: map[string]*nodeInfo{
						"test-node": {
							ID:    "test-node",
							class: "test-class",
							resourceStats: &resourceStats{
								allocatableResources: &resources{cpu: 5182, memory: 985},
								allocatedResources:   &resources{cpu: 1000, memory: 512},
							},
						},
					},
					resourceStats: &resourceStats{
						allocatableResources: &resources{cpu: 5182, memory: 985},
						allocatedResources:   &resources{cpu: 1000, memory: 512},
					},
					allocations: map[string]string{"existing-alloc": "running", "new-alloc": "running"},
				},
			},
			name: "new alloc on node with existing alloc",
		},
	}

	for _, tc := range testCases {
		tc.inputHandler.handleAllocMessageRunning(tc.inputClass, tc.inputAlloc)
		assert.Equal(t, tc.expectedNodePoolResult, tc.inputHandler.nodePool, tc.name)
	}
}

func intToPointer(i int) *int { return &i }
