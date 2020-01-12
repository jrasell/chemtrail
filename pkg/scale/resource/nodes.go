package resource

import (
	"github.com/hashicorp/nomad/api"
)

func (n *updateHandler) runNodeUpdateHandler() {
	n.logger.Info().Msg("starting Chemtrail Nomad node update handler")

	for {
		select {
		case <-n.shutdownChan:
			n.logger.Info().Msg("shutting down Chemtrail node update handler")
			return
		case msg := <-n.nodeUpdateChan:
			go n.handleClientMessage(msg)
		}
	}
}

func (n *updateHandler) handleClientMessage(msg interface{}) {
	node, ok := msg.(*api.Node)
	if !ok {
		n.logger.Error().Msg("received unexpected node update message type")
		return
	}
	n.logger.Debug().
		Str("node-id", node.ID).
		Str("node-status", node.Status).
		Str("node-eligibility", node.SchedulingEligibility).
		Msg("received node update message to handle")

	// Perform our node class check before we handle the actual message.
	n.checkNodeClass(node)

	switch node.Status {
	case "initializing":
		// If the client is starting up, there is no need to process this update. We will catch the
		// node joining the cluster when it becomes ready and process the information then.
	case "down":
		n.handleNodeUnavailableMessage(node)
	case "ready":
		if node.SchedulingEligibility == "eligible" {
			n.handleNodeAvailableMessage(node)
		} else if node.SchedulingEligibility == "ineligible" {
			n.handleNodeUnavailableMessage(node)
		}
	}
}

func (n *updateHandler) handleNodeAvailableMessage(node *api.Node) {
	// Ensure we have the node class map updated.
	n.nodeClassLock.Lock()
	n.nodeClass[node.ID] = node.NodeClass
	n.nodeClassLock.Unlock()

	n.nodePoolLock.Lock()
	defer n.nodePoolLock.Unlock()

	// Attempt to read the node class out of the state.
	stored, ok := n.nodePool[node.NodeClass]
	if !ok {
		n.nodePool[node.NodeClass] = &classInfo{
			allocations: make(map[string]string),
			class:       node.NodeClass,
			nodes:       make(map[string]*nodeInfo),
			resourceStats: &resourceStats{
				allocatedResources:   &resources{},
				allocatableResources: &resources{},
			},
		}
	}

	// Ensure we do not process a node which we are already tracking in the correct state. If we
	// skip this check, we can overwrite our resource stats.
	if stored != nil {
		if storedNode, ok := stored.nodes[node.ID]; ok {
			if node.Status == storedNode.status {
				n.logger.Debug().
					Str("node-id", node.ID).
					Str("node-status", storedNode.status).
					Msg("node has already been processed with current status")
				return
			}
		}
	}

	// Build the required information of the node.
	info := nodeInfo{
		ID:          node.ID,
		status:      node.Status,
		class:       node.NodeClass,
		eligibility: node.SchedulingEligibility,
		resourceStats: &resourceStats{
			allocatedResources:   &resources{},
			allocatableResources: n.getNodeAllocatableResources(node),
		},
	}
	n.nodePool[node.NodeClass].nodes[node.ID] = &info

	n.logger.Info().
		Str("node-id", info.ID).
		Str("node-class", info.class).
		Float64("node-allocatable-cpu", info.resourceStats.allocatableResources.cpu).
		Float64("node-allocatable-memory", info.resourceStats.allocatableResources.memory).
		Msg("added node to Chemtrail internal state")

	// Update the node class pool high level resource tracking stats.
	n.nodePool[node.NodeClass].resourceStats.allocatableResources.cpu += info.resourceStats.allocatableResources.cpu
	n.nodePool[node.NodeClass].resourceStats.allocatableResources.memory += info.resourceStats.allocatableResources.memory
}

func (n *updateHandler) handleNodeUnavailableMessage(node *api.Node) {

	// If the node class is not being tracked, we do not need to do anything as this handled node
	// needs to be removed from tracking anyway.
	if _, ok := n.nodePool[node.NodeClass]; !ok {
		return
	}

	if info, ok := n.nodePool[node.NodeClass].nodes[node.ID]; ok {

		n.nodePoolLock.Lock()

		delete(n.nodePool[node.NodeClass].nodes, node.ID)

		n.nodePool[node.NodeClass].resourceStats.allocatableResources.cpu -= info.resourceStats.allocatableResources.cpu
		n.nodePool[node.NodeClass].resourceStats.allocatableResources.memory -= info.resourceStats.allocatableResources.memory

		n.nodePoolLock.Unlock()

		n.nodeClassLock.Lock()
		delete(n.nodeClass, node.ID)
		n.nodeClassLock.Unlock()
	}
}

// checkNodeClass is used to check whether the received node has its Class set. If not, we set the
// Chemtrail default.
func (n *updateHandler) checkNodeClass(node *api.Node) {
	if node.NodeClass == "" {
		n.logger.Debug().
			Str("node-id", node.ID).
			Msg("node has empty class parameter, using Chemtrail default")
		node.NodeClass = "chemtrail-default"
	}
}

// getNodeAllocatableResources takes the desired node and calculates the amount of allocatable
// resources. Older versions of Nomad do not include the NodeResources and ReservedResources fields
// so this must be checked and the older fields used if possible. GH-26 contains additional
// details.
func (n *updateHandler) getNodeAllocatableResources(node *api.Node) *resources {
	r := resources{}

	if node.NodeResources != nil && node.ReservedResources != nil {
		r.cpu = float64(node.NodeResources.Cpu.CpuShares) - float64(node.ReservedResources.Cpu.CpuShares)
		r.memory = float64(node.NodeResources.Memory.MemoryMB) - float64(node.ReservedResources.Memory.MemoryMB)
	}
	if node.Resources != nil && node.Reserved != nil {
		r.cpu = float64(*node.Resources.CPU) - float64(*node.Reserved.CPU)
		r.memory = float64(*node.Resources.MemoryMB) - float64(*node.Reserved.MemoryMB)
	}
	return &r
}
