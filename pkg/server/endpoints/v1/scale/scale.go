package scale

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/jrasell/chemtrail/pkg/helper"
	"github.com/jrasell/chemtrail/pkg/scale"
	"github.com/jrasell/chemtrail/pkg/state"
	"github.com/rs/zerolog"
)

type Server struct {
	Logger        zerolog.Logger
	Scale         scale.Scale
	PolicyBackend state.PolicyBackend
	ScaleBackend  state.ScaleBackend
}

func (s *Server) PostScaleIn(w http.ResponseWriter, r *http.Request) {

	msg, err := s.prepareScaleMessage(mux.Vars(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if msg == nil {
		http.Error(w, "no scaling policy held for targeted class", http.StatusUnprocessableEntity)
		return
	}
	msg.Direction = state.ScaleDirectionIn

	code, err := s.Scale.OKToScale(msg)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	go s.Scale.InvokeScaling(msg)

	helper.WriteJSONResponse(w, []byte(fmt.Sprintf("{\"ID\":\"%s\"}", msg.ID)), http.StatusOK, s.Logger)
}

func (s *Server) PostScaleOut(w http.ResponseWriter, r *http.Request) {
	msg, err := s.prepareScaleMessage(mux.Vars(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if msg == nil {
		http.Error(w, "no scaling policy held for targeted class", http.StatusUnprocessableEntity)
		return
	}
	msg.Direction = state.ScaleDirectionOut

	code, err := s.Scale.OKToScale(msg)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	go s.Scale.InvokeScaling(msg)

	helper.WriteJSONResponse(w, []byte(fmt.Sprintf("{\"ID\":\"%s\"}", msg.ID)), http.StatusOK, s.Logger)
}

func (s *Server) prepareScaleMessage(vars map[string]string) (*state.ScalingRequest, error) {
	targetClass := vars["client-class"]

	classPolicy, err := s.PolicyBackend.GetPolicy(targetClass)
	if err != nil {
		return nil, err
	}

	if classPolicy == nil {
		return nil, nil
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	msg := state.ScalingRequest{
		ID:     id,
		Policy: classPolicy,
	}
	return &msg, nil
}
