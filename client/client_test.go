package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

// UserTestSuite declares what global state it needs.
type UserTestSuite struct {
	suite.Suite
	serverEndpoint string
}

// SetupTest initializes global state.
func (s *UserTestSuite) SetupTest() {
	s.serverEndpoint = "https://example.com/endpoint"
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestBuildRequest() {
	testCases := []struct {
		name        string
		inputParam  string
		expectedURL string
		expectedErr error
	}{
		{
			name:        "Valid get request",
			inputParam:  "my input param",
			expectedURL: "https://example.com/endpoint?input=my+input+param",
			expectedErr: nil,
		},
		{
			name:        "Empty param",
			inputParam:  "",
			expectedURL: "",
			expectedErr: errEmptyParam,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			req, err := buildRequest(context.Background(), s.serverEndpoint, testCase.inputParam)
			if testCase.expectedErr != nil {
				s.Equal(testCase.expectedErr, err)
			} else {
				s.Equal(testCase.expectedURL, req.URL.String())
			}
		})
	}
}

func (s *UserTestSuite) TestUnmarshalResponseBody() {
	responseBodyUsers := []byte(`[{"name":"Klaas Kippengaas","Email":"klaas@kippengaas.nl"},{"name":"Henk de Vries","email":"henkdevries@hotmail.com"}]`)
	expectedUsers := Users{
		{
			Name:  "Klaas Kippengaas",
			Email: "klaas@kippengaas.nl",
		},
		{
			Name:  "Henk de Vries",
			Email: "henkdevries@hotmail.com",
		},
	}

	r := bytes.NewReader(responseBodyUsers)
	users, err := unmarshalResponseBody(r)
	s.Require().NoError(err)
	s.Equal(expectedUsers, users)
}
