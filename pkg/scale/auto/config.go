package auto

import (
	"github.com/jrasell/chemtrail/pkg/client"
	"github.com/jrasell/chemtrail/pkg/scale"
	"github.com/jrasell/chemtrail/pkg/scale/resource"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/rs/zerolog"
)

type Config struct {
	Nomad    *client.Nomad
	Logger   zerolog.Logger
	Policy   state.PolicyBackend
	Resource resource.Handler
	Scale    scale.Scale
	Interval int
	Threads  int
}
