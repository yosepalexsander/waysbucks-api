package handler

import "net/http"

func GetUsers(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("get all users"))
}

func GetUser(w http.ResponseWriter, r *http.Request)  {
	
}

func Register(w http.ResponseWriter, r *http.Request)  {
	
}

func Login(w http.ResponseWriter, r *http.Request)  {
	
}