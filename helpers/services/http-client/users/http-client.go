package http_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	shared_dtos "image-reports/shared/dtos"
)

type HttpClient interface {
	CheckCredentials(credentials shared_dtos.UserCredentials) (shared_dtos.UserResponse, error)
}

type httpClient struct {
}

func NewHttpClient() HttpClient {
	return &httpClient{}
}

func (client *httpClient) CheckCredentials(credentials shared_dtos.UserCredentials) (shared_dtos.UserResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(credentials); err != nil {
		return shared_dtos.UserResponse{}, err
	}

	resp, err := http.Post("http://localhost:8080/v1/auth", "text/json", &buf)
	if err != nil {
		return shared_dtos.UserResponse{}, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		errRes := shared_dtos.ErrorResponse{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&errRes); err != nil {
			return shared_dtos.UserResponse{}, err
		}
		return shared_dtos.UserResponse{}, errors.New(errRes.Error)
	}

	usr := shared_dtos.UserResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&usr); err != nil {
		return usr, err
	}

	return usr, nil
}
