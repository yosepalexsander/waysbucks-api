package handler

import (
	"encoding/json"
	"net/http"
)

type commonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func internalServerError(w http.ResponseWriter) {
	resp, _ := json.Marshal(commonResponse{
		Error:   true,
		Message: "server error",
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resp)
}

func forbidden(w http.ResponseWriter) {
	resp, _ := json.Marshal(commonResponse{
		Error:   true,
		Message: "access denied",
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write(resp)
}

func notFound(w http.ResponseWriter) {
	resp, _ := json.Marshal(commonResponse{
		Error:   true,
		Message: "resource not found",
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write(resp)
}

func responseOK(w http.ResponseWriter, resp []byte) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func badRequest(w http.ResponseWriter, msg string) {
	resp, _ := json.Marshal(commonResponse{
		Error:   true,
		Message: msg,
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}

func serviceUnavailable(w http.ResponseWriter, msg string) {
	resp, _ := json.Marshal(commonResponse{
		Message: msg,
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write(resp)
}
