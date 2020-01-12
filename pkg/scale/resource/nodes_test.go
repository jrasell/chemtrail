package resource

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
)

func Test_updateHandler_handleNodeAvailableMessage(t *testing.T) {
	testCases := []struct {
		inputNode               *api.Node
		inputHandler            *updateHandler
		expectedNodePoolResult  map[string]*classInfo
		expectedNodeClassResult map[string]string
		name                    string
	}{
		{
			inputNode: &api.Node{
				ID:        "test-node",
				NodeClass: "test-class",
				NodeResources: &api.NodeResources{
					Cpu:    api.NodeCpuResources{CpuShares: 5182},
					Memory: api.NodeMemoryResources{MemoryMB: 985},
				},
				ReservedResources: &api.NodeReservedResources{
					Cpu:    api.NodeReservedCpuResources{CpuShares: 0},
					Memory: api.NodeReservedMemoryResources{MemoryMB: 0},
				},
			},
			inputHandler: &updateHandler{
				nodePool:  make(map[string]*classInfo),
				nodeClass: make(map[string]string),
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
					allocations: make(map[string]string),
				},
			},
			expectedNodeClassResult: map[string]string{"test-node": "test-class"},
			name:                    "first node availability message",
		},

		{
			inputNode: &api.Node{
				ID:        "new-node",
				NodeClass: "test-class",
				NodeResources: &api.NodeResources{
					Cpu:    api.NodeCpuResources{CpuShares: 5182},
					Memory: api.NodeMemoryResources{MemoryMB: 985},
				},
				ReservedResources: &api.NodeReservedResources{
					Cpu:    api.NodeReservedCpuResources{CpuShares: 0},
					Memory: api.NodeReservedMemoryResources{MemoryMB: 0},
				},
			},
			inputHandler: &updateHandler{
				nodePool: map[string]*classInfo{
					"test-class": {
						class: "test-class",
						nodes: map[string]*nodeInfo{
							"existing-node": {
								ID:    "existing-node",
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
				nodeClass: map[string]string{"existing-node": "test-class"},
			},
			expectedNodePoolResult: map[string]*classInfo{
				"test-class": {
					class: "test-class",
					nodes: map[string]*nodeInfo{
						"existing-node": {
							ID:    "existing-node",
							class: "test-class",
							resourceStats: &resourceStats{
								allocatableResources: &resources{cpu: 5182, memory: 985},
								allocatedResources:   &resources{cpu: 0, memory: 0},
							},
						},
						"new-node": {
							ID:    "new-node",
							class: "test-class",
							resourceStats: &resourceStats{
								allocatableResources: &resources{cpu: 5182, memory: 985},
								allocatedResources:   &resources{cpu: 0, memory: 0},
							},
						},
					},
					resourceStats: &resourceStats{
						allocatableResources: &resources{cpu: 10364, memory: 1970},
						allocatedResources:   &resources{cpu: 0, memory: 0},
					},
					allocations: make(map[string]string),
				},
			},
			expectedNodeClassResult: map[string]string{"new-node": "test-class", "existing-node": "test-class"},
			name:                    "second node of class to be discovered",
		},

		{
			inputNode: &api.Node{
				ID:        "new-node",
				NodeClass: "new-class",
				NodeResources: &api.NodeResources{
					Cpu:    api.NodeCpuResources{CpuShares: 5182},
					Memory: api.NodeMemoryResources{MemoryMB: 985},
				},
				ReservedResources: &api.NodeReservedResources{
					Cpu:    api.NodeReservedCpuResources{CpuShares: 0},
					Memory: api.NodeReservedMemoryResources{MemoryMB: 0},
				},
			},
			inputHandler: &updateHandler{
				nodePool: map[string]*classInfo{
					"existing-class": {
						class: "existing-class",
						nodes: map[string]*nodeInfo{
							"existing-class-node": {
								ID:    "existing-class-node",
								class: "existing-class",
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
				nodeClass: map[string]string{"existing-class-node": "existing-class"},
			},
			expectedNodePoolResult: map[string]*classInfo{
				"existing-class": {
					class: "existing-class",
					nodes: map[string]*nodeInfo{
						"existing-class-node": {
							ID:    "existing-class-node",
							class: "existing-class",
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
				"new-class": {
					class: "new-class",
					nodes: map[string]*nodeInfo{
						"new-node": {
							ID:    "new-node",
							class: "new-class",
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
			expectedNodeClassResult: map[string]string{"existing-class-node": "existing-class", "new-node": "new-class"},
			name:                    "new node class discovered, state contains existing class",
		},
	}

	for _, tc := range testCases {
		tc.inputHandler.handleNodeAvailableMessage(tc.inputNode)
		assert.Equal(t, tc.expectedNodePoolResult, tc.inputHandler.nodePool, tc.name)
		assert.Equal(t, tc.expectedNodeClassResult, tc.inputHandler.nodeClass, tc.name)
	}
}

func Test_updateHandler_handleNodeUnavailableMessage(t *testing.T) {
	testCases := []struct {
		inputNode               *api.Node
		inputHandler            *updateHandler
		expectedNodePoolResult  map[string]*classInfo
		expectedNodeClassResult map[string]string
		name                    string
	}{
		{
			inputNode: &api.Node{NodeClass: "test-class"},
			inputHandler: &updateHandler{
				nodePool:  make(map[string]*classInfo),
				nodeClass: make(map[string]string),
			},
			expectedNodePoolResult:  map[string]*classInfo{},
			expectedNodeClassResult: map[string]string{},
			name:                    "handle unavailable node as first processed",
		},
		{
			inputNode: &api.Node{ID: "fake-id", NodeClass: "test-class"},
			inputHandler: &updateHandler{
				nodePool: map[string]*classInfo{
					"test-class": {
						class: "test-class",
						nodes: map[string]*nodeInfo{
							"fake-id": {
								ID:    "fake-id",
								class: "test-class",
								resourceStats: &resourceStats{
									allocatableResources: &resources{cpu: 1000, memory: 1000},
									allocatedResources:   &resources{cpu: 0, memory: 0},
								},
							},
						},
						resourceStats: &resourceStats{
							allocatableResources: &resources{cpu: 1000, memory: 1000},
							allocatedResources:   &resources{cpu: 0, memory: 0},
						},
						allocations: make(map[string]string),
					},
				},
				nodeClass: map[string]string{"fake-id": "test-class"},
			},
			expectedNodePoolResult: map[string]*classInfo{
				"test-class": {
					class: "test-class",
					nodes: map[string]*nodeInfo{},
					resourceStats: &resourceStats{
						allocatableResources: &resources{cpu: 0, memory: 0},
						allocatedResources:   &resources{cpu: 0, memory: 0},
					},
					allocations: make(map[string]string),
				},
			},
			expectedNodeClassResult: map[string]string{},
			name:                    "handle only node being tracked of class",
		},
		{
			inputNode: &api.Node{ID: "node-for-removal", NodeClass: "test-class"},
			inputHandler: &updateHandler{
				nodePool: map[string]*classInfo{
					"test-class": {
						class: "test-class",
						nodes: map[string]*nodeInfo{
							"node-for-removal": {
								ID:    "node-for-removal",
								class: "test-class",
								resourceStats: &resourceStats{
									allocatableResources: &resources{cpu: 1000, memory: 1000},
									allocatedResources:   &resources{cpu: 0, memory: 0},
								},
							},
							"node-to-keep-trucking": {
								ID:    "node-to-keep-trucking",
								class: "test-class",
								resourceStats: &resourceStats{
									allocatableResources: &resources{cpu: 1000, memory: 1000},
									allocatedResources:   &resources{cpu: 0, memory: 0},
								},
							},
						},
						resourceStats: &resourceStats{
							allocatableResources: &resources{cpu: 2000, memory: 2000},
							allocatedResources:   &resources{cpu: 0, memory: 0},
						},
						allocations: make(map[string]string),
					},
				},
				nodeClass: map[string]string{"node-for-removal": "test-class", "node-to-keep-trucking": "test-class"},
			},
			expectedNodePoolResult: map[string]*classInfo{
				"test-class": {
					class: "test-class",
					nodes: map[string]*nodeInfo{
						"node-to-keep-trucking": {
							ID:    "node-to-keep-trucking",
							class: "test-class",
							resourceStats: &resourceStats{
								allocatableResources: &resources{cpu: 1000, memory: 1000},
								allocatedResources:   &resources{cpu: 0, memory: 0},
							},
						},
					},
					resourceStats: &resourceStats{
						allocatableResources: &resources{cpu: 1000, memory: 1000},
						allocatedResources:   &resources{cpu: 0, memory: 0},
					},
					allocations: make(map[string]string),
				},
			},
			expectedNodeClassResult: map[string]string{"node-to-keep-trucking": "test-class"},
			name:                    "handle node being tracked in class with other nodes",
		},
	}

	for _, tc := range testCases {
		tc.inputHandler.handleNodeUnavailableMessage(tc.inputNode)
		assert.Equal(t, tc.expectedNodePoolResult, tc.inputHandler.nodePool, tc.name)
		assert.Equal(t, tc.expectedNodeClassResult, tc.inputHandler.nodeClass, tc.name)
	}
}

func Test_updateHandler_checkNodeClass(t *testing.T) {
	testCases := []struct {
		inputNode          *api.Node
		expectedClassValue string
		name               string
	}{
		{
			inputNode:          &api.Node{NodeClass: ""},
			expectedClassValue: "chemtrail-default",
			name:               "empty node class",
		},
		{
			inputNode:          &api.Node{NodeClass: "high-memory"},
			expectedClassValue: "high-memory",
			name:               "operator configured node class",
		},
	}
	uh := &updateHandler{}

	for _, tc := range testCases {
		uh.checkNodeClass(tc.inputNode)
		assert.Equal(t, tc.expectedClassValue, tc.inputNode.NodeClass, tc.name)
	}
}

func Test_updateHandler_getNodeAllocatableResources(t *testing.T) {
	testCases := []struct {
		inputNode      *api.Node
		expectedOutput *resources
		name           string
	}{
		{
			inputNode: &api.Node{
				NodeResources: &api.NodeResources{
					Cpu:    api.NodeCpuResources{CpuShares: 5182},
					Memory: api.NodeMemoryResources{MemoryMB: 985},
				},
				ReservedResources: &api.NodeReservedResources{
					Cpu:    api.NodeReservedCpuResources{CpuShares: 0},
					Memory: api.NodeReservedMemoryResources{MemoryMB: 0},
				},
			},
			expectedOutput: &resources{cpu: 5182, memory: 985},
			name:           "0 reserved resources",
		},
		{
			inputNode: &api.Node{
				NodeResources: &api.NodeResources{
					Cpu:    api.NodeCpuResources{CpuShares: 5182},
					Memory: api.NodeMemoryResources{MemoryMB: 985},
				},
				ReservedResources: &api.NodeReservedResources{
					Cpu:    api.NodeReservedCpuResources{CpuShares: 10},
					Memory: api.NodeReservedMemoryResources{MemoryMB: 10},
				},
			},
			expectedOutput: &resources{cpu: 5172, memory: 975},
			name:           "reserved CPU and memory resources",
		},
		{
			inputNode: &api.Node{
				NodeResources: &api.NodeResources{
					Cpu:    api.NodeCpuResources{CpuShares: 5182},
					Memory: api.NodeMemoryResources{MemoryMB: 985},
				},
				ReservedResources: &api.NodeReservedResources{
					Cpu:    api.NodeReservedCpuResources{CpuShares: 10},
					Memory: api.NodeReservedMemoryResources{MemoryMB: 0},
				},
			},
			expectedOutput: &resources{cpu: 5172, memory: 985},
			name:           "reserved CPU but no reserved memory resources",
		},
		{
			inputNode: &api.Node{
				NodeResources: &api.NodeResources{
					Cpu:    api.NodeCpuResources{CpuShares: 5182},
					Memory: api.NodeMemoryResources{MemoryMB: 985},
				},
				ReservedResources: &api.NodeReservedResources{
					Cpu:    api.NodeReservedCpuResources{CpuShares: 0},
					Memory: api.NodeReservedMemoryResources{MemoryMB: 10},
				},
			},
			expectedOutput: &resources{cpu: 5182, memory: 975},
			name:           "no reserved CPU but reserved memory resources",
		},
		{
			inputNode: &api.Node{
				Resources: &api.Resources{
					CPU:      intToPointer(5182),
					MemoryMB: intToPointer(985),
				},
				Reserved: &api.Resources{
					CPU:      intToPointer(0),
					MemoryMB: intToPointer(0),
				},
			},
			expectedOutput: &resources{cpu: 5182, memory: 985},
			name:           "gh-26 older version of Nomad",
		},
	}
	uh := &updateHandler{}

	for _, tc := range testCases {
		actualOutput := uh.getNodeAllocatableResources(tc.inputNode)
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}
