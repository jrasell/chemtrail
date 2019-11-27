package system

import (
	"encoding/json"
	"net/http"

	metrics "github.com/armon/go-metrics"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/rs/zerolog"
)

type Server struct {
	logger    zerolog.Logger
	telemetry *metrics.InmemSink
}

func NewServer(logger zerolog.Logger, telemetry *metrics.InmemSink) *Server {
	return &Server{
		logger:    logger.With().Str("component", "endpoint-system").Logger(),
		telemetry: telemetry,
	}
}

func (s *Server) GetHealth(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSONResponse(w, []byte("{\"status\":\"ok\"}"), http.StatusOK, s.logger)
}

func (s *Server) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metricData, err := s.telemetry.DisplayMetrics(w, r)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get latest telemetry data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(metricData)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal HTTP response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.WriteJSONResponse(w, out, http.StatusOK, s.logger)
}
