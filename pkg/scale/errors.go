package scale

import "errors"

const (
	scalingPreconditionCheckFailedMsg = "scaling activity failed precondition check"
)

var (
	errScalingPolicyDisabled      = errors.New("scaling policy is currently disabled")
	errScalingProviderNotFound    = errors.New("scaling provider not found in configuration")
	errNoNodesFoundInClass        = errors.New("no Nomad nodes found of client class")
	errScalingInCountCheckFailed  = errors.New("scaling in activity would break policy minimum threshold")
	errScalingOutCountCheckFailed = errors.New("scaling out activity would break policy maximum threshold")
)
