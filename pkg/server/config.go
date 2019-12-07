package server

import (
	serverCfg "github.com/jrasell/chemtrail/pkg/config/server"
)

// Config contains all the required configuration to setup and run a Chemtrail server as defined by
// an operator.
type Config struct {
	Autoscale *serverCfg.AutoscalerConfig
	Provider  *serverCfg.ProviderConfig
	Server    *serverCfg.Config
	Storage   *serverCfg.StorageConfig
	TLS       *serverCfg.TLSConfig
	Telemetry *serverCfg.TelemetryConfig
}

// The system API endpoints.
const (
	routeGetSystemMetricsName    = "GetSystemMetrics"
	routeGetSystemMetricsPattern = "/v1/system/metrics"
	routeGetSystemHealthName     = "GetSystemHealth"
	routeGetSystemHealthPattern  = "/v1/system/health"
	telemetryInterval            = 10
)

// The policy API endpoints.
const (
	routeGetPoliciesName     = "GetPolicies"
	routeGetPoliciesPattern  = "/v1/policies"
	routeGetPolicyName       = "GetPolicy"
	routeGetPolicyPattern    = "/v1/policy/{client-class}"
	routePutPolicyName       = "PutPolicy"
	routePutPolicyPattern    = "/v1/policy/{client-class}"
	routeDeletePolicyName    = "DeletePolicy"
	routeDeletePolicyPattern = "/v1/policy/{client-class}"
)

// The scale API endpoints.
const (
	routeGetScaleStatusName        = "GetScaleStatus"
	routeGetScaleStatusPattern     = "/v1/scale/status"
	routeGetScaleStatusInfoName    = "GetScaleStatusInfo"
	routeGetScaleStatusInfoPattern = "/v1/scale/status/{id}"
	routePostScaleInName           = "PostScaleIn"
	routeScaleInPattern            = "/v1/scale/in/{client-class}"
	routeScaleOutName              = "PostScaleOut"
	routeScaleOutPattern           = "/v1/scale/out/{client-class}"
)
