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

	reqStruct := []struct {
		Email    string
		Password string
	}{
		{
			Email:    "user2@gmail.com",
			Password: "12345678",
		},
		{
			Email:    "user3@gmail.com",
			Password: "12345678",
		},
	}

	for _, req := range reqStruct {
		reqBody, _ := json.Marshal(req)
		requestReader := bytes.NewReader(reqBody)
		request, err := http.NewRequest("POST", baseUrl+"/auth/login", requestReader)
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
		response.Body.Close()
	}
}

func TestRegisterWithInvalidBody(t *testing.T) {
	t.Parallel()
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	reqStruct := []struct {
		Email    string
		Password string
	}{
		{
			Email:    "test14gmail.com",
			Password: "123456",
		},
		{
			Email:    "test13gmail.com",
			Password: "123456",
		},
	}
	for _, req := range reqStruct {
		reqBody, _ := json.Marshal(req)
		requestReader := bytes.NewReader(reqBody)
		request, err := http.NewRequest("POST", baseUrl+"/auth/register", requestReader)
		if err != nil {
			t.Errorf("failed to create new request")
		}
		response, err := client.Do(request)
		if err != nil {
			t.Errorf(err.Error())
		}

		if response.StatusCode != 400 {
			t.Errorf("response status code is not 400")
		}
		response.Body.Close()
	}
}

func TestVerifyGoogleTokenID(t *testing.T) {
	_, err := VerifyTokenID("eyJhbGciOiJSUzI1NiIsImtpZCI6IjM4ZjM4ODM0NjhmYzY1OWFiYjQ0NzVmMzYzMTNkMjI1ODVjMmQ3Y2EiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJuYmYiOjE2NTM4NzYxNzUsImF1ZCI6IjM1MTE0OTEyNTczNi1icDFpNWdiNm4wcm85c3JkYTBhMDdxZm01bDZoYmtpZy5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjEwNjI0NTkyNTYzMDk2NDA1NDIxMiIsImVtYWlsIjoia3VyYW1hbmF0c3UwM0BnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXpwIjoiMzUxMTQ5MTI1NzM2LWJwMWk1Z2I2bjBybzlzcmRhMGEwN3FmbTVsNmhia2lnLmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwibmFtZSI6Ikt1cmFtYSBOYXRzdSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQVRYQUp3NVpvYUQ1ZVBlSXR2WjNHWGxzMy1HZG1hTFJOVjdsbWVpOElhdz1zOTYtYyIsImdpdmVuX25hbWUiOiJLdXJhbWEiLCJmYW1pbHlfbmFtZSI6Ik5hdHN1IiwiaWF0IjoxNjUzODc2NDc1LCJleHAiOjE2NTM4ODAwNzUsImp0aSI6ImNhMjgxNmU0NjM0NDViYmJmNWMyYjBkYWIyMmRlNmNiZTBhNzJhMjgifQ.UubNJq_YGCSV1-PjyPuvK226t5Fp599B1J1lqbdT6qBZnCb0fgAxCTfRj7o_sEfXRZsVrwgzmyhMesGDLNxox90x6eVg5ba_zqEuJLGgJncvqZWxoEkuQrbhq2onk4b14ilVSqtr1vibNojaYwD7vXqlPi4mxWSzpfuRwgptrrCa5WXwFJzuRViraPAkQPkuxUhqymzrmhlfSwLKiH-YNYCJncRK3u4ByNUusQzbB7DmqTGdIvLjg6pNpzmEbDeDBprtprR3UWuxBFcRJoeJCfmhSsQ6cnF5BwLA-oPwi24TBtzU23aQpNEKASI_5BtPdUgw5mPmeuYQ3qRkqCJPZw")
	if err != nil {
		t.Errorf("Verify token failed with error: %s", err.Error())
	}
}
