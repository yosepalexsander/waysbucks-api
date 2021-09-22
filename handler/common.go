package handler

import (
	"encoding/json"
	"net/http"
)


func internalServerError(w http.ResponseWriter)  {
	w.WriteHeader(http.StatusInternalServerError)
	resp, _ := json.Marshal(CommonResponse{
		Message: "server error",
	})
	w.Write(resp)
}

func forbidden(w http.ResponseWriter)  {
	w.WriteHeader(http.StatusForbidden)
	resp, _ := json.Marshal(CommonResponse{
		Message: "access denied",
	})
	w.Write(resp)
}

func notFound(w http.ResponseWriter)  {
	resp, _ := json.Marshal(CommonResponse{
		Message: "resource not found",
	})
	w.WriteHeader(http.StatusNotFound)
	w.Write(resp)
}

func responseOK(w http.ResponseWriter, resp []byte)  {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func badRequest(w http.ResponseWriter, msg string)  {
	resp, _ := json.Marshal(CommonResponse{
		Message: msg,
	})
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}