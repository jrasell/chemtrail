package provider

import "github.com/jrasell/chemtrail/pkg/state"

// ClientProvider is the interface that needs to be implemented by providers which are responsible
// scaling Nomad client machines.
type ClientProvider interface {

	// Name returns the human readable name for the provider.
	Name() string

	// ScaleIn will trigger a scaling in event of the provider. When implementing this function, it
	// should handle only provider interactions as well as manging activity update event based on
	// what occurs. When calling this function, all safety checks should have been completed to
	// ensure no policy parameters are violated. The passed target should be used to identify the
	// node in a way that the provider can understand.
	ScaleIn(req *state.ScalingRequest, target string) error

	// ScaleOut will trigger a scaling out event of the provider. When implementing this function,
	// it should handle only provider interactions as well as manging activity update event based
	// on what occurs. When calling this function, all safety checks should have been completed to
	// ensure no policy parameters are violated.
	ScaleOut(req *state.ScalingRequest) error
}
