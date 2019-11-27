package awsasg

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/jrasell/chemtrail/pkg/scale/provider"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ provider.ClientProvider = (*ClientProvider)(nil)

type ClientProvider struct {
	log       zerolog.Logger
	asgClient *autoscaling.Client
	ec2Client *ec2.Client
	eventChan chan *state.EventMessage
}

const (
	configKeyASGName = "asg-name"
)

func NewAWSASGProvider(log zerolog.Logger, eventChan chan *state.EventMessage) provider.ClientProvider {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil
	}

	return &ClientProvider{
		log:       log.With().Str("provider", state.AWSAutoScaling.String()).Logger(),
		asgClient: autoscaling.New(cfg),
		ec2Client: ec2.New(cfg),
		eventChan: eventChan,
	}
}

// Name satisfies the provider.ClientProvider Name interface function.
func (a *ClientProvider) Name() string { return state.AWSAutoScaling.String() }

// ScaleOut satisfies the provider.ClientProvider ScaleOut interface function.
func (a *ClientProvider) ScaleOut(msg *state.ScalingRequest) error {
	asgName, err := a.getProviderConfigValue(msg, configKeyASGName)
	if err != nil {
		return err
	}

	asg, err := a.describeAutoScalingGroup(asgName)
	a.handleEvent(eventTypeDesc, err, nil, msg.ID)
	if err != nil {
		return err
	}

	input := autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(asgName),
		AvailabilityZones:    asg.AvailabilityZones,
		DesiredCapacity:      aws.Int64(*asg.DesiredCapacity + int64(msg.Policy.ScaleOutCount)),
	}

	_, err = a.asgClient.UpdateAutoScalingGroupRequest(&input).Send(context.Background())
	a.handleEvent(eventTypeUpdate, err, nil, msg.ID)
	return err
}

// ScaleIn satisfies the provider.ClientProvider ScaleIn interface function.
func (a *ClientProvider) ScaleIn(msg *state.ScalingRequest, id string) error {
	asgName, err := a.getProviderConfigValue(msg, configKeyASGName)
	if err != nil {
		return err
	}

	asgInput := autoscaling.DetachInstancesInput{
		AutoScalingGroupName:           aws.String(asgName),
		InstanceIds:                    []string{id},
		ShouldDecrementDesiredCapacity: aws.Bool(true),
	}

	_, err = a.asgClient.DetachInstancesRequest(&asgInput).Send(context.Background())
	a.handleEvent(eventTypeUpdate, err, aws.String(id), msg.ID)
	if err != nil {
		return err
	}

	ec2Input := ec2.TerminateInstancesInput{DryRun: aws.Bool(false), InstanceIds: []string{id}}

	_, err = a.ec2Client.TerminateInstancesRequest(&ec2Input).Send(context.Background())
	a.handleEvent(eventTypeTerminate, err, aws.String(id), msg.ID)
	return err
}

func (a *ClientProvider) describeAutoScalingGroup(name string) (*autoscaling.AutoScalingGroup, error) {
	input := autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []string{name}}

	resp, err := a.asgClient.DescribeAutoScalingGroupsRequest(&input).Send(context.Background())
	if err != nil {
		return nil, err
	}

	asgs := len(resp.AutoScalingGroups)
	if asgs != 1 {
		return nil, errors.Errorf("described %v AutoScaling Groups, expected 1", asgs)
	}

	return &resp.AutoScalingGroups[0], nil
}

func (a *ClientProvider) getProviderConfigValue(msg *state.ScalingRequest, key string) (string, error) {
	v, ok := msg.Policy.ProviderConfig[key]
	if !ok {
		return "", errors.Errorf("required provider config key %s not found", key)
	}
	return v, nil
}
