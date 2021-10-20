package handler

import (
	"encoding/json"
	"net/http"
)

type commonResponse struct {
	Message string `json:"message"`
}

func internalServerError(w http.ResponseWriter) {
	resp, _ := json.Marshal(commonResponse{
		Message: "server error",
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resp)
}

func forbidden(w http.ResponseWriter) {
	resp, _ := json.Marshal(commonResponse{
		Message: "access denied",
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write(resp)
}

func notFound(w http.ResponseWriter) {
	resp, _ := json.Marshal(commonResponse{
		Message: "resource not found",
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write(resp)
}

func responseOK(w http.ResponseWriter, resp []byte) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func badRequest(w http.ResponseWriter, msg string) {
	resp, _ := json.Marshal(commonResponse{
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
