package api

type System struct {
	client *Client
}

func (c *Client) System() *System {
	return &System{client: c}
}

type HealthResp struct {
	Status string
}

func (s *System) Health() (*HealthResp, error) {
	var resp HealthResp
	err := s.client.get("/v1/system/health", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
