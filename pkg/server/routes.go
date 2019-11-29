package server

import (
	"net/http"

	policyV1 "github.com/jrasell/chemtrail/pkg/server/endpoints/v1/policy"
	scaleV1 "github.com/jrasell/chemtrail/pkg/server/endpoints/v1/scale"
	systemV1 "github.com/jrasell/chemtrail/pkg/server/endpoints/v1/system"
	"github.com/jrasell/chemtrail/pkg/server/router"
)

type routes struct {
	policy *policyV1.Server
	scale  *scaleV1.Server
	system *systemV1.Server
}

func (h *HTTPServer) setupRoutes() *router.RouteTable {
	h.logger.Debug().Msg("setting up HTTP server routes")
	return &router.RouteTable{h.setupSystemRoutes(), h.setupScaleRoutes(), h.setupPolicyRoutes()}
}

func (h *HTTPServer) setupSystemRoutes() []router.Route {
	h.logger.Debug().Msg("setting up HTTP server system routes")

	h.routes.system = systemV1.NewServer(h.logger, h.telemetry)

	return router.Routes{
		router.Route{
			Name:    routeGetSystemHealthName,
			Method:  http.MethodGet,
			Pattern: routeGetSystemHealthPattern,
			Handler: h.routes.system.GetHealth,
		},
		router.Route{
			Name:    routeGetSystemMetricsName,
			Method:  http.MethodGet,
			Pattern: routeGetSystemMetricsPattern,
			Handler: h.routes.system.GetMetrics,
		},
	}
}

func (h *HTTPServer) setupScaleRoutes() []router.Route {
	h.logger.Debug().Msg("setting up HTTP server scale routes")

	h.routes.scale = &scaleV1.Server{
		Logger:        h.logger.With().Str("component", "endpoint-scale").Logger(),
		Scale:         h.scaler,
		PolicyBackend: h.policyState,
		ScaleBackend:  h.scaleState,
	}

	return router.Routes{
		router.Route{
			Name:    routePostScaleInName,
			Method:  http.MethodPost,
			Pattern: routeScaleInPattern,
			Handler: h.routes.scale.PostScaleIn,
		},
		router.Route{
			Name:    routeScaleOutName,
			Method:  http.MethodPost,
			Pattern: routeScaleOutPattern,
			Handler: h.routes.scale.PostScaleOut,
		},
		router.Route{
			Name:    routeGetScaleStatusName,
			Method:  http.MethodGet,
			Pattern: routeGetScaleStatusPattern,
			Handler: h.routes.scale.GetScaleStatus,
		},
		router.Route{
			Name:    routeGetScaleStatusInfoName,
			Method:  http.MethodGet,
			Pattern: routeGetScaleStatusInfoPattern,
			Handler: h.routes.scale.GetScaleStatusInfo,
		},
	}
}

func (h *HTTPServer) setupPolicyRoutes() []router.Route {
	h.logger.Debug().Msg("setting up HTTP server policy routes")

	h.routes.policy = policyV1.NewServer(h.logger, h.policyState)

	return router.Routes{
		router.Route{
			Name:    routeGetPoliciesName,
			Method:  http.MethodGet,
			Pattern: routeGetPoliciesPattern,
			Handler: h.routes.policy.GetPolicies,
		},
		router.Route{
			Name:    routeGetPolicyName,
			Method:  http.MethodGet,
			Pattern: routeGetPolicyPattern,
			Handler: h.routes.policy.GetPolicy,
		},
		router.Route{
			Name:    routePutPolicyName,
			Method:  http.MethodPut,
			Pattern: routePutPolicyPattern,
			Handler: h.routes.policy.PutPolicy,
		},
		router.Route{
			Name:    routeDeletePolicyName,
			Method:  http.MethodDelete,
			Pattern: routeDeletePolicyPattern,
			Handler: h.routes.policy.DeletePolicy,
		},
	}
}
