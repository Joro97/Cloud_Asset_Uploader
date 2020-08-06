package responses

import (
	"encoding/json"
	"net/http"
)

// WriteBadRequest writes StatusBadRequest response and the given message to the given ResponseWriter.
func WriteBadRequest(wr http.ResponseWriter, message string) {
	wr.WriteHeader(http.StatusBadRequest)
	resp, _ := json.Marshal(message)
	wr.Write(resp)
}

// WriteOkResponse writes StatusOK response and the passed data to the given ResponseWriter.
func WriteOkResponse(wr http.ResponseWriter, data interface{}) {
	wr.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(data)
	wr.Write(resp)
}

// WriteInternalServerErrorResponse writes StatusInternalServerError response and the given message to the given ResponseWriter.
func WriteInternalServerErrorResponse(wr http.ResponseWriter, message string) {
	wr.WriteHeader(http.StatusInternalServerError)
	resp, _ := json.Marshal(message)
	wr.Write(resp)
}

// WriteResourceNotFoundResponse writes StatusNotFound response and the given message to the given ResponseWriter.
func WriteResourceNotFoundResponse(wr http.ResponseWriter, message string) {
	wr.WriteHeader(http.StatusNotFound)
	resp, _ := json.Marshal(message)
	wr.Write(resp)
}

