package resource

import (
	"sync"

	"github.com/jrasell/chemtrail/pkg/client"
	"github.com/rs/zerolog"
)

// Handler is the interface which governs how the Chemtrail resource manager state is interacted
// with. In order to lessen the load on the Nomad API, resource tracking is done via watchers,
// rather than API calls at the time they are needed. This means Chemtrail can make faster
// decisions and responses to requests, at a much lower impact to the Nomad cluster servers.
type Handler interface {

	// GetNodeUpdateChan returns the channel where the node watcher should send updates regarding
	// the Nomad client cluster pool.
	GetNodeUpdateChan() chan interface{}

	// RunNodeUpdateHandler triggers the process which handles node updates and listens on the
	// channel as returned via GetNodeUpdateChan. When implementing this interface, the process
	// listening on the channel should expect a type of github.com/hashicorp/nomad/api/.(*Node).
	RunNodeUpdateHandler()

	// GetAllocUpdateChan returns the channel where the alloc watcher should send updates regarding
	// the Nomad cluster allocations.
	GetAllocUpdateChan() chan interface{}

	// RunAllocUpdateHandler triggers the process which handles node updates and listens on the
	// channel as returned via GetAllocUpdateChan. When implementing this interface, the process
	// listening on the channel should expect a type of
	// github.com/hashicorp/nomad/api/.(*Allocation).
	RunAllocUpdateHandler()

	// StopUpdateHandlers is used to stop all the running update handlers within the resource
	// process.
	StopUpdateHandlers()

	// GetNodesOfClass returns the currently stored node mapping relating to a particular class as
	// requested. The key of the map is the Nomad NodeID as specified by
	// github.com/hashicorp/nomad/api/.(*Node.ID).
	GetNodesOfClass(class string) map[string]*nodeInfo

	// GetClassResourceAllocation is used to perform allocation calculations for the class in
	// question. The function will use the stored class statistics to calculate the percentage of
	// resources currently allocated. The calculation uses allocated rather than actually used as
	// Nomad does not oversubscribe, and jobs will fail to run if there are not enough allocatable
	// resources.
	GetClassResourceAllocation(class string) (*AllocatedStats, error)

	// GetLeastAllocatedNodeInClass is used to find the node in the class pool which is the least
	// allocated. This is the current default and hardcoded mode for scaling in as it reduces the
	// amount of resources that need to be migrated across the cluster.
	GetLeastAllocatedNodeInClass(class string) *nodeInfo
}

type updateHandler struct {
	logger zerolog.Logger
	nomad  *client.Nomad

	// nodePool is the state for the Chemtrail node resources. The map is keyed by the Nomad client
	// class parameter.
	nodePool     map[string]*classInfo
	nodePoolLock sync.RWMutex

	// nodeClass keeps a tracking of the nodes and the class they are configured with. This is
	// useful as allocations do not contain this data but do contain the nodeID. Therefore we can
	// use this map for a quick lookup to understand which class pool the allocation resources are
	// associated with.
	nodeClass     map[string]string
	nodeClassLock sync.RWMutex

	// nodeUpdateChan is where updates from the node watcher should be sent for processing.
	nodeUpdateChan chan interface{}

	// allocUpdateChan is where updates from the alloc watcher should be sent for processing.
	allocUpdateChan chan interface{}

	// shutdownChan is used to coordinate the shutdown of the resource processes in a clean manner.
	shutdownChan chan struct{}
}
