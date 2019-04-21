package writer

import (
	"net/http"

	"github.com/mailru/easyjson"
)

func WriteResponseJSON(w http.ResponseWriter, status int, data easyjson.Marshaler) {
	w.WriteHeader(status)
	easyjson.MarshalToHTTPResponseWriter(data, w)
}
