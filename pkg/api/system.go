package api

import metrics "github.com/armon/go-metrics"

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

// Metrics returns the currently stored telemetry information from a Chemtrail server.
func (s *System) Metrics() (*metrics.MetricsSummary, error) {
	var resp metrics.MetricsSummary
	err := s.client.get("/v1/system/metrics", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
