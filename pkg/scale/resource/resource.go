package resource

import (
	"math"

	"github.com/jrasell/chemtrail/pkg/client"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// classInfo is used to track statistics relating to a pool of nodes defined via a Nomad class.
type classInfo struct {
	class         string
	nodes         map[string]*nodeInfo
	resourceStats *resourceStats
	allocations   map[string]string
}

// nodeInfo represents an individual node within the Nomad cluster.
type nodeInfo struct {
	ID            string
	class         string
	status        string
	eligibility   string
	resourceStats *resourceStats
}

// resourceStats represents the currently tracked CPU and memory stats for the component. This is
// designed so that it can be used to perform the basic calculations needed to measure how consumed
// the component is.
type resourceStats struct {
	allocatableResources *resources
	allocatedResources   *resources
}

// resources is the basic tracked resources which Chemtrail can scale on.
type resources struct {
	cpu    float64
	memory float64
}

// AllocatedStats is used to return information about the current allocated resources within a
// class as a percentage.
type AllocatedStats struct {

	// CPU is the currently allocated CPU as a percentage of the overall allocatable CPU resource
	// within a class.
	CPU float64

	// Memory is the currently allocated memory as a percentage of the overall allocatable memory
	// resource within a class.
	Memory float64
}

type handler struct {
	logger      zerolog.Logger
	nodeManager *updateHandler
}

// GetLeastAllocatedNodeInClass satisfies the GetLeastAllocatedNodeInClass function on the Handler
// interface.
func (h *handler) GetLeastAllocatedNodeInClass(class string) *nodeInfo {
	classInfo, ok := h.nodeManager.nodePool[class]
	if !ok {
		return nil
	}

	var (
		lowestPercentage float64
		lowestNode       *nodeInfo
	)

	for _, node := range classInfo.nodes {

		// In its current form, we need to protect the node Chemtrail is running on from scaling.
		// This will change in the future when HA features come in, but it easier now to just skip
		// the assessment of the Chemtrail node.
		if node.ID == h.nodeManager.nomad.NodeID {
			continue
		}

		stats := h.calculateAllocatedPercentageStats(node.resourceStats)

		// We need to set an initial benchmark for comparison.
		if lowestNode == nil && lowestPercentage == 0 {
			lowestPercentage = stats.Memory
			lowestNode = node
		}

		if stats.Memory < lowestPercentage {
			lowestPercentage = stats.Memory
			lowestNode = node
		}

		if stats.CPU < lowestPercentage {
			lowestPercentage = stats.Memory
			lowestNode = node
		}
	}

	return lowestNode
}

// GetClassResourceAllocation satisfies the GetClassResourceAllocation function on the Handler interface.
func (h *handler) GetClassResourceAllocation(class string) (*AllocatedStats, error) {
	h.nodeManager.nodePoolLock.RLock()
	defer h.nodeManager.nodePoolLock.RUnlock()

	// Check that we have nodes within our state, otherwise we have nothing to calculate.
	nodes, ok := h.nodeManager.nodePool[class]
	if !ok {
		return nil, errors.New("no nodes of class found")
	}
	return h.calculateAllocatedPercentageStats(nodes.resourceStats), nil
}

// StopUpdateHandlers satisfies the StopUpdateHandlers function on the Handler interface.
func (h *handler) StopUpdateHandlers() { close(h.nodeManager.shutdownChan) }

// GetAllocUpdateChan satisfies the GetAllocUpdateChan function on the Handler interface.
func (h *handler) GetAllocUpdateChan() chan interface{} { return h.nodeManager.allocUpdateChan }

// RunAllocUpdateHandler satisfies the RunAllocUpdateHandler function on the Handler interface.
func (h *handler) RunAllocUpdateHandler() { go h.nodeManager.runAllocUpdateHandler() }

// GetNodeUpdateChan satisfies the GetNodeUpdateChan function on the Handler interface.
func (h *handler) GetNodeUpdateChan() chan interface{} { return h.nodeManager.nodeUpdateChan }

// RunNodeUpdateHandler satisfies the RunNodeUpdateHandler function on the Handler interface.
func (h *handler) RunNodeUpdateHandler() { go h.nodeManager.runNodeUpdateHandler() }

// GetNodesOfClass satisfies the GetNodesOfClass function on the Handler interface.
func (h *handler) GetNodesOfClass(class string) map[string]*nodeInfo {
	if classInfo, ok := h.nodeManager.nodePool[class]; ok {
		return classInfo.nodes
	}
	return nil
}

// NewHandler creates a new resource handler for interactions with the resource state stored within
// Chemtrail based off updates from Nomad via watchers.
func NewHandler(logger zerolog.Logger, nomad *client.Nomad) Handler {
	return &handler{
		logger: logger,
		nodeManager: &updateHandler{
			logger:          logger,
			nomad:           nomad,
			nodePool:        make(map[string]*classInfo),
			nodeClass:       make(map[string]string),
			nodeUpdateChan:  make(chan interface{}),
			allocUpdateChan: make(chan interface{}),
			shutdownChan:    make(chan struct{}),
		},
	}
}

// calculateAllocatedPercentageStats is used to calculate the percentage of resources allocated out
// of the total allocatable resources.
func (h *handler) calculateAllocatedPercentageStats(stats *resourceStats) *AllocatedStats {
	cpu := (stats.allocatedResources.cpu * float64(100)) / stats.allocatableResources.cpu
	mem := (stats.allocatedResources.memory * float64(100)) / stats.allocatableResources.memory

	return &AllocatedStats{
		CPU:    math.Round(cpu),
		Memory: math.Round(mem),
	}
}
