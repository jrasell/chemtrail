package helper

import (
	"net/http"

	"github.com/rs/zerolog"
)

const (
	headerKeyContentType       = "Content-Type"
	headerValueContentTypeJSON = "application/json; charset=utf-8"
)

func WriteJSONResponse(w http.ResponseWriter, bytes []byte, statusCode int, log zerolog.Logger) {
	w.Header().Set(headerKeyContentType, headerValueContentTypeJSON)
	w.WriteHeader(statusCode)
	if _, err := w.Write(bytes); err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
	}
}
