package client

import (
	"os"

	"github.com/hashicorp/nomad/api"
)

// Nomad is a wrapper around the Nomad client which includes the current nodeID if found. This
// allows the Chemtrail server to protect the node it is running on from scale in activities which
// would cause undesirable situations.
type Nomad struct {
	Client *api.Client
	NodeID string
}

// NewNomadClient builds the reusable Nomad client.
func NewNomadClient() (*Nomad, error) {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	id, err := chemtrailNodeID(c)
	if err != nil {
		return nil, err
	}

	return &Nomad{
		Client: c,
		NodeID: id,
	}, nil
}

// chemtrailNodeID attempts to find the Nomad nodeID Chemtrail is running.
func chemtrailNodeID(client *api.Client) (string, error) {
	if envVar := os.Getenv("NOMAD_ALLOC_ID"); envVar == "" {
		return envVar, nil
	}

	self, err := client.Agent().Self()
	if err != nil {
		return "", err
	}
	return self.Stats["client"]["node_id"], nil
}
