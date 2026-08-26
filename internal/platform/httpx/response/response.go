package response

import (
	"encoding/json"
	"net/http"
)

const InternalErrorMsg = "oops! internal error"

const internalErrorBody = `{"error":"` + InternalErrorMsg + `"}`

type errBody struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		body, status = []byte(internalErrorBody), http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, errBody{Error: msg})
}

func ErrorWithFields(w http.ResponseWriter, status int, msg string, fields map[string]string) {
	JSON(w, status, errBody{Error: msg, Fields: fields})
}

func ErrorInternal(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, InternalErrorMsg)
}
