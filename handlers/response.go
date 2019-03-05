package handlers

import (
	"net/http"

	"github.com/mailru/easyjson"
)

func writeResponseJSON(w http.ResponseWriter, status int, data easyjson.Marshaler) {
	w.WriteHeader(status)
	easyjson.MarshalToHTTPResponseWriter(data, w)
}
