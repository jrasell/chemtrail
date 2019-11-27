package policy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/rs/zerolog"
)

type Server struct {
	logger        zerolog.Logger
	policyBackend state.PolicyBackend
}

func NewServer(log zerolog.Logger, policyBackend state.PolicyBackend) *Server {
	return &Server{
		logger:        log.With().Str("component", "endpoint-policy").Logger(),
		policyBackend: policyBackend,
	}
}

func (s *Server) GetPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := s.policyBackend.GetPolicies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(policies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helper.WriteJSONResponse(w, bytes, http.StatusOK, s.logger)
}

func (s *Server) GetPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	class := vars["client-class"]

	policy, err := s.policyBackend.GetPolicy(class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if policy == nil {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(policy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helper.WriteJSONResponse(w, bytes, http.StatusOK, s.logger)
}

func (s *Server) PutPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	class := vars["client-class"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to read request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var p state.ClientScalingPolicy

	if err := json.Unmarshal(body, &p); err != nil {
		s.logger.Error().Err(err).Msg("failed to unmarshal request body into policy")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.Class = class

	if err := p.Validate(); err != nil {
		s.logger.Error().Err(err).Msg("failed to validate scale policy")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err = s.policyBackend.PutPolicy(&p); err != nil {
		s.logger.Error().Err(err).Msg("failed write scaling policy to storage backend")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	class := vars["client-class"]

	if err := s.policyBackend.DeletePolicy(class); err != nil {
		s.logger.Error().Err(err).Msg("failed delete scaling policy in storage backend")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
