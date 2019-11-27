package scale

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/jrasell/chemtrail/pkg/helper"
)

func (s *Server) GetScaleStatus(w http.ResponseWriter, r *http.Request) {
	events, err := s.ScaleBackend.GetScalingActivities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helper.WriteJSONResponse(w, bytes, http.StatusOK, s.Logger)
}

func (s *Server) GetScaleStatusInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	scaleID, err := uuid.FromString(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scaleState, err := s.ScaleBackend.GetScalingActivity(scaleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if scaleState == nil {
		http.Error(w, "ScaleID not found", http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(scaleState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helper.WriteJSONResponse(w, bytes, http.StatusOK, s.Logger)
}
