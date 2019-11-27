package state

import "github.com/pkg/errors"

// PolicyBackend is the interface which must be satisfied in order to store Chemtrail scaling
// policies.
type PolicyBackend interface {

	// GetPolicies returns all the currently stored policies from the backend. The map is keyed by
	// the class identifier.
	GetPolicies() (map[string]*ClientScalingPolicy, error)

	// GetPolicy takes an input class identifier and returns the scaling policy if it is found
	// within the backend.
	GetPolicy(class string) (*ClientScalingPolicy, error)

	// PutPolicy is used to write a new or update an existing class policy. Updates should
	// overwrite any stored information for the class.
	PutPolicy(policy *ClientScalingPolicy) error

	// DeletePolicy is used to delete the class scaling policy if it exists within the backend
	// storage. If the specified class does not have a policy within the backend, the call will be
	// no-op.
	DeletePolicy(class string) error
}

// ClientScalingPolicy represents a Nomad client class scaling document. It is used to define all
// parameters and configuration options needed to perform scaling evaluations and requests.
type ClientScalingPolicy struct {
	Enabled        bool                    `json:"Enabled"`
	Class          string                  `json:"Class"`
	MinCount       int                     `json:"MinCount"`
	MaxCount       int                     `json:"MaxCount"`
	ScaleOutCount  int                     `json:"ScaleOutCount"`
	ScaleInCount   int                     `json:"ScaleInCount"`
	Provider       ClientProvider          `json:"Provider"`
	ProviderConfig map[string]string       `json:"ProviderConfig"`
	Checks         map[string]*PolicyCheck `json:"Checks"`
}

// PolicyCheck is an individual check to be performed as part of an autoscaling evaluation of the
// Nomad client class.
type PolicyCheck struct {

	// Enabled is a boolean flag to identify whether this specific check should be actively run or
	// not.
	Enabled bool `json:"Enabled"`

	// ScaleResource identifies the Nomad resource which will be checked within this policy check
	// evaluation.
	Resource ScaleResource `json:"ScaleResource"`

	// ComparisonOperator determines how the value if compared to the theshold.
	ComparisonOperator ComparisonOperator `json:"ComparisonOperator"`

	// ComparisonPercentage is the float64 value that will be compared to the result of the metric
	// query.
	ComparisonPercentage float64 `json:"ComparisonPercentage"`

	// Action is the scaling action that should be taken if the queried metric fails the comparison
	// check.
	Action ComparisonAction `json:"Action"`
}

// Validate can be used to validate the contents of a scaling policy. This ensures it contains the
// minimum viable config available for both Chemtrail, and the backend provider it is intended to
// be used with.
func (c ClientScalingPolicy) Validate() error {

	// Validate the provider setting. This is required within all scaling policies.
	if err := c.Provider.Validate(); err != nil {
		return err
	}

	// Currently in its MVP stage, Chemtrail can only handle scaling in with a count of 1. This
	// will change in the future so the parameter is there but its currently protected.
	if c.ScaleInCount > 1 {
		return errors.New("currently Chemtrail can only handle ScaleInCount of 1")
	}

	// Depending on the provider, we will have different base requirements for the config.
	switch c.Provider {
	case AWSAutoScaling:
		if _, ok := c.ProviderConfig["asg-name"]; !ok {
			return errors.New("provider config must include \"asg-name\" parameter")
		}
	}

	// Iterate over the checks and validate the required components. The first error is returned,
	// rather than collecting.
	for name, check := range c.Checks {
		if err := check.Resource.Validate(); err != nil {
			return errors.Wrap(err, "failed to validate check"+name)
		}

		if err := check.ComparisonOperator.Validate(); err != nil {
			return errors.Wrap(err, "failed to validate check"+name)
		}

		if err := check.Action.Validate(); err != nil {
			return errors.Wrap(err, "failed to validate check"+name)
		}
	}
	return nil
}

// ClientProvider is an identifier to the backend which provides the Nomad client workers. This is
// used to ensure the correct APIs are called when wanting to perform scaling activities.
type ClientProvider string

// String returns the string representation of the client provider.
func (c ClientProvider) String() string { return string(c) }

// Validate can be used to ensure the passed ClientProvider is valid and able to be handled by the
// current Chemtrail version.
func (c ClientProvider) Validate() error {
	switch c {
	case AWSAutoScaling:
		return nil
	default:
		return errors.Errorf("unsupported client provider %s", c.String())
	}
}

const (
	// AWSAutoScaling uses AWS AutoScaling groups to provide the client workers.
	AWSAutoScaling ClientProvider = "aws-autoscaling"
)

// ComparisonOperator is the operator used when evaluating a metric value against a threshold.
type ComparisonOperator string

// String returns the string form of the ComparisonOperator.
func (co ComparisonOperator) String() string { return string(co) }

// Validate checks the ComparisonOperator is a valid and that it can be handled within the
// autoscaler.
func (co ComparisonOperator) Validate() error {
	switch co {
	case ComparisonGreaterThan, ComparisonLessThan:
		return nil
	default:
		return errors.Errorf("ComparisonOperator %s is not a valid option", co.String())
	}
}

const (
	ComparisonGreaterThan ComparisonOperator = "greater-than"
	ComparisonLessThan    ComparisonOperator = "less-than"
)

// ComparisonAction is the action to take if the metric breaks the threshold.
type ComparisonAction string

// String returns the string form of the ComparisonAction.
func (ca ComparisonAction) String() string { return string(ca) }

// Validate checks the ComparisonAction is a valid and that it can be handled within the
// autoscaler.
func (ca ComparisonAction) Validate() error {
	switch ca {
	case ActionScaleIn, ActionScaleOut:
		return nil
	default:
		return errors.Errorf("Action %s is not a valid option", ca.String())
	}
}

const (
	// ActionScaleIn performs a scale in operation.
	ActionScaleIn ComparisonAction = "scale-in"

	// ActionScaleOut performs a scale out operation.
	ActionScaleOut ComparisonAction = "scale-out"
)

// ScaleResource is the Nomad resource to evaluate within the defined check.
type ScaleResource string

// String returns the string form of the ScaleResource.
func (r ScaleResource) String() string { return string(r) }

// Validate checks the ScaleResource is a valid and that it can be handled within the autoscaler.
func (r ScaleResource) Validate() error {
	switch r {
	case ScaleResourceCPU, ScaleResourceMemory:
		return nil
	default:
		return errors.Errorf("ScaleResource %s is not a valid option", r.String())
	}
}

const (
	// ScaleResourceCPU represents the CPU resource stanza parameter in a Nomad job as specified:
	// https://www.nomadproject.io/docs/job-specification/resources.html#cpu
	ScaleResourceCPU ScaleResource = "cpu"

	// ScaleResourceMemory represents the memory resource stanza parameter in a Nomad job as
	// specified: https://www.nomadproject.io/docs/job-specification/resources.html#memory
	ScaleResourceMemory ScaleResource = "memory"
)
