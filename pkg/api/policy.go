package api

type Policy struct {
	client *Client
}

func (c *Client) Policy() *Policy {
	return &Policy{client: c}
}

type ScalingPolicy struct {
	Enabled        bool
	Class          string
	MinCount       int
	MaxCount       int
	ScaleOutCount  int
	ScaleInCount   int
	Provider       string
	ProviderConfig map[string]string
	Checks         map[string]Check
}

type Check struct {
	Enabled              bool
	Resource             string
	ComparisonOperator   string
	ComparisonPercentage float64
	Action               string
}

func (p *Policy) Delete(class string) error {
	return p.client.delete("/v1/policy/"+class, nil)
}

func (p *Policy) Write(class string, policy *ScalingPolicy) error {
	return p.client.put("/v1/policy/"+class, policy, nil, nil)
}

func (p *Policy) List() (*map[string]ScalingPolicy, error) {
	var resp map[string]ScalingPolicy
	err := p.client.get("/v1/policies", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *Policy) Info(class string) (*ScalingPolicy, error) {
	var resp ScalingPolicy
	err := p.client.get("/v1/policy/"+class, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
