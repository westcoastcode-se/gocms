package cms

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int
	Message string
}

func returnSuccess(rw http.ResponseWriter, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(jsonResponse)
	return
}

func returnErrorResponse(rw http.ResponseWriter, code int, message string) {
	response := &ErrorResponse{
		Code:    code,
		Message: message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	_, _ = rw.Write(jsonResponse)
	return
}

func returnForbidden(rw http.ResponseWriter) {
	returnErrorResponse(rw, http.StatusForbidden, "Forbidden")
}

func returnNotFound(rw http.ResponseWriter) {
	returnErrorResponse(rw, http.StatusNotFound, "Not Found")
}

func returnMethodNotAllowed(rw http.ResponseWriter) {
	returnErrorResponse(rw, http.StatusMethodNotAllowed, "Method Not Allowed")
}

