package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

const baseUrl = "http://localhost:8080/api/v1"

func TestLogin(t *testing.T) {
	t.Parallel()
  client := &http.Client{
		Timeout: 15 * time.Second,
	}

	reqStruct := []Login_Req{
		{
			Email: "test14@gmail.com",
			Password: "12345678",
		},
		{
			Email: "test13@gmail.com",
			Password: "12345678",
		},
	}

	for _, req := range reqStruct {
		go func (req Login_Req)  {
			reqBody, _ := json.Marshal(req) 
			requestReader := bytes.NewReader(reqBody)
			request, err := http.NewRequest("POST", baseUrl+"/login", requestReader)
			if err != nil {
				t.Errorf("failed to create new request")
			}
			response, err := client.Do(request)
			if err != nil {
				t.Errorf(err.Error())
			}
		
			if response.StatusCode != 200 {
				t.Errorf("response status code is not 200")
			}
			defer response.Body.Close()
		}(req)
	}
}