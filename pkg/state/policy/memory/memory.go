package memory

import (
	"sync"

	"github.com/jrasell/chemtrail/pkg/state"
)

type PolicyBackend struct {
	policies map[string]*state.ClientScalingPolicy
	sync.RWMutex
}

func NewPolicyBackend() state.PolicyBackend {
	return &PolicyBackend{policies: make(map[string]*state.ClientScalingPolicy)}
}

// GetPolicies satisfies the GetPolicies function on the state.PolicyBackend interface.
func (p *PolicyBackend) GetPolicies() (map[string]*state.ClientScalingPolicy, error) {
	p.RLock()
	val := p.policies
	p.RUnlock()
	return val, nil
}

// GetPolicy satisfies the GetPolicy function on the state.PolicyBackend interface.
func (p *PolicyBackend) GetPolicy(class string) (*state.ClientScalingPolicy, error) {
	p.RLock()
	defer p.RUnlock()

	if val, ok := p.policies[class]; ok {
		return val, nil
	}
	return nil, nil
}

// PutPolicy satisfies the PutPolicy function on the state.PolicyBackend interface.
func (p *PolicyBackend) PutPolicy(policy *state.ClientScalingPolicy) error {
	p.Lock()
	p.policies[policy.Class] = policy
	p.Unlock()
	return nil
}

// DeletePolicy satisfies the DeletePolicy function on the state.PolicyBackend interface.
func (p *PolicyBackend) DeletePolicy(class string) error {
	p.Lock()
	delete(p.policies, class)
	p.Unlock()
	return nil
}
