package e2e

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

type apiRequest struct {
	Request  *http.Request
	Recorder *httptest.ResponseRecorder
}

func NewAPIRequest(method, url string, data []byte) *apiRequest {
	instance := apiRequest{}
	instance.Request = httptest.NewRequest(method, url, bytes.NewBuffer(data))
	return &instance
}

func (a *apiRequest) GetRecorder(requestHandler http.Handler) *httptest.ResponseRecorder {
	if a == nil {
		return nil
	}
	a.Recorder = httptest.NewRecorder()
	requestHandler.ServeHTTP(a.Recorder, a.Request)
	return a.Recorder
}
