package api

import (
	"github.com/gofrs/uuid"
	"github.com/jrasell/chemtrail/pkg/state"
)

type Scale struct {
	client *Client
}

func (c *Client) Scale() *Scale {
	return &Scale{client: c}
}

type ScaleResp struct {
	ID string
}

func (s *Scale) StatusList() (*map[uuid.UUID]state.ScalingActivity, error) {
	var resp map[uuid.UUID]state.ScalingActivity
	err := s.client.get("/v1/scale/status", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Scale) StatusInfo(id string) (*state.ScalingActivity, error) {
	var resp state.ScalingActivity
	err := s.client.get("/v1/scale/status/"+id, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Scale) In(class string) (*ScaleResp, error) {
	var resp ScaleResp
	err := s.client.put("/v1/scale/in/"+class, nil, &resp, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Scale) Out(class string) (*ScaleResp, error) {
	var resp ScaleResp
	err := s.client.put("/v1/scale/out/"+class, nil, &resp, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
