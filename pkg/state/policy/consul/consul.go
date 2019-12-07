package consul

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/pkg/errors"
)

// baseKVPath is suffixed to the operator supplied Consul path as the location where Chemtrail
// policies are stored.
const baseKVPath = "policies/"

// PolicyBackend is the Consul implementation of the state.PolicyBackend interface.
type PolicyBackend struct {
	path string
	kv   *api.KV
}

// NewPolicyBackend returns the Consul implementation of the state.PolicyBackend interface.
func NewPolicyBackend(path string, client *api.Client) state.PolicyBackend {
	return &PolicyBackend{
		path: path + baseKVPath,
		kv:   client.KV(),
	}
}

// GetPolicies satisfies the GetPolicies function on the state.PolicyBackend interface.
func (p PolicyBackend) GetPolicies() (map[string]*state.ClientScalingPolicy, error) {
	kv, _, err := p.kv.List(p.path, nil)
	if err != nil {
		return nil, err
	}

	out := make(map[string]*state.ClientScalingPolicy)

	// If there are no KV entries, return the empty map.
	if kv == nil {
		return out, nil
	}

	for i := range kv {
		p := &state.ClientScalingPolicy{}

		if err := json.Unmarshal(kv[i].Value, p); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
		}
		keySplit := strings.Split(kv[i].Key, "/")

		out[keySplit[len(keySplit)-1]] = p
	}

	return out, nil
}

// GetPolicy satisfies the GetPolicy function on the state.PolicyBackend interface.
func (p PolicyBackend) GetPolicy(class string) (*state.ClientScalingPolicy, error) {
	kv, _, err := p.kv.Get(p.path+class, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil {
		return nil, nil
	}

	out := &state.ClientScalingPolicy{}

	if err := json.Unmarshal(kv.Value, out); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal Consul KV value")
	}

	return out, nil
}

// PutPolicy satisfies the PutPolicy function on the state.PolicyBackend interface.
func (p PolicyBackend) PutPolicy(policy *state.ClientScalingPolicy) error {
	marshal, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	pair := &api.KVPair{
		Key:   p.path + policy.Class,
		Value: marshal,
	}

	_, err = p.kv.Put(pair, nil)
	return err
}

// DeletePolicy satisfies the DeletePolicy function on the state.PolicyBackend interface.
func (p PolicyBackend) DeletePolicy(class string) error {
	_, err := p.kv.Delete(p.path+class, nil)
	return err
}
