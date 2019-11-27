package resource

import (
	"github.com/hashicorp/nomad/api"
)

func (n *updateHandler) runAllocUpdateHandler() {
	n.logger.Info().Msg("starting Chemtrail Nomad alloc update handler")

	for {
		select {
		case <-n.shutdownChan:
			n.logger.Info().Msg("shutting down Chemtrail alloc update handler")
			return
		case msg := <-n.allocUpdateChan:
			go n.handleAllocMessage(msg)
		}
	}
}

func (n *updateHandler) handleAllocMessage(msg interface{}) {
	alloc, ok := msg.(*api.Allocation)
	if !ok {
		n.logger.Error().Msg("received unexpected node alloc update message type")
		return
	}
	n.logger.Debug().
		Str("alloc-id", alloc.ID).
		Str("alloc-client-status", alloc.ClientStatus).
		Msg("received alloc update message to handle")

	n.nodeClassLock.RLock()
	class, ok := n.nodeClass[alloc.NodeID]
	n.nodeClassLock.RUnlock()

	// We should never get here, but if we do it indicates a big problem as we are not tracking the
	// node on which the allocation is running.
	if !ok {
		n.logger.Error().
			Str("alloc-id", alloc.ID).
			Str("alloc-client-status", alloc.ClientStatus).
			Str("node-id", alloc.NodeID).
			Msg("received alloc update running on a class we are not tracking")
		return
	}

	status, ok := n.nodePool[class].allocations[alloc.NodeID]
	if ok && status == alloc.ClientStatus {
		n.logger.Debug().
			Str("alloc-id", alloc.ID).
			Str("alloc-client-status", alloc.ClientStatus).
			Msg("alloc handler has previously processed allocation with the same status")
		return
	}

	switch alloc.ClientStatus {
	case "pending":
		// Pending is an intermediate stage of an allocations lifecycle, and currently this is
		// ignored until it reaches a different state.
		return
	case "running":
		n.handleAllocMessageRunning(class, alloc)
	default:
		n.handleAllocMessageTerminal(class, alloc)
	}
}

func (n *updateHandler) handleAllocMessageTerminal(class string, alloc *api.Allocation) {
	// Lock the worker pool so we can safely update a number of stats based on this allocation.
	n.nodePoolLock.Lock()
	defer n.nodePoolLock.Unlock()

	// Allocations which have a terminal status might not have been discovered before but we do
	// not want to continue any further as not to remove their resources from the stats where
	// they are not accounted for.
	if _, ok := n.nodePool[class].allocations[alloc.ID]; !ok {
		return
	}

	// Delete the allocation from our tracking.
	delete(n.nodePool[class].allocations, alloc.ID)

	n.nodePool[class].nodes[alloc.NodeID].resourceStats.allocatedResources.cpu -= float64(*alloc.Resources.CPU)
	n.nodePool[class].nodes[alloc.NodeID].resourceStats.allocatedResources.memory -= float64(*alloc.Resources.MemoryMB)

	// Update the class pools resource stats.
	n.nodePool[class].resourceStats.allocatedResources.cpu -= float64(*alloc.Resources.CPU)
	n.nodePool[class].resourceStats.allocatedResources.memory -= float64(*alloc.Resources.MemoryMB)
}

func (n *updateHandler) handleAllocMessageRunning(class string, alloc *api.Allocation) {
	// Lock the worker pool so we can safely update a number of stats based on this allocation.
	n.nodePoolLock.Lock()

	if status, ok := n.nodePool[class].allocations[alloc.ID]; ok && status == alloc.ClientStatus {
		n.nodePoolLock.Unlock()
		return
	}

	// Update the allocated resource stats.
	n.nodePool[class].nodes[alloc.NodeID].resourceStats.allocatedResources.cpu += float64(*alloc.Resources.CPU)
	n.nodePool[class].nodes[alloc.NodeID].resourceStats.allocatedResources.memory += float64(*alloc.Resources.MemoryMB)

	// Update the class pools resource stats.
	n.nodePool[class].resourceStats.allocatedResources.cpu += float64(*alloc.Resources.CPU)
	n.nodePool[class].resourceStats.allocatedResources.memory += float64(*alloc.Resources.MemoryMB)

	n.nodePool[class].allocations[alloc.ID] = alloc.ClientStatus

	// Our work here is done.
	n.nodePoolLock.Unlock()
}
