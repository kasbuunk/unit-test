// Package client is used to test clients.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

var (
	errEmptyParam  = errors.New("input param must not be empty")
	errStatusNotOK = errors.New("status code is not ok")
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Users []User

// RequestUsers sends the request built by buildRequest.
func RequestUsers(ctx context.Context) ([]User, error) {
	req, err := buildRequest(ctx, "https://restserverwithoutaclient/myendpoint", "param")
	if err != nil {
		return nil, errors.Wrap(err, "constructing get request")
	}

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, errStatusNotOK
	}
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}

	users, err := unmarshalResponseBody(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response body")
	}
	return users, nil
}

func buildRequest(ctx context.Context, serverEndpoint, inputParam string) (*http.Request, error) {
	if inputParam == "" {
		return nil, errEmptyParam
	}
	reqURL := fmt.Sprintf("%v?input=%v", serverEndpoint, url.QueryEscape(inputParam))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "constructing clients get request")
	}
	return req, nil
}

func unmarshalResponseBody(body io.Reader) (Users, error) {
	var dest Users

	users, err := io.ReadAll(body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}
	err = json.Unmarshal(users, &dest)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response body")
	}
	return dest, nil
}
