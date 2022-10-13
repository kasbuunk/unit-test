package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	handler              http.Handler
	acceptedHeader       string
	expectedResponseBody string
}

func (s *MiddlewareTestSuite) SetupSuite() {
	s.expectedResponseBody = "Middleware succeeded!"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprint(w, s.expectedResponseBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.FailNow("writing response")
			return
		}
	})
	s.acceptedHeader = "X-Header"
	s.handler = CheckHeaderExists(testHandler, s.acceptedHeader)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) TestMiddleware() {
	testCases := []struct {
		name                   string
		header                 string
		expectedResponseStatus int
	}{
		{
			name:                   "Happy flow",
			header:                 s.acceptedHeader,
			expectedResponseStatus: http.StatusOK,
		},
		{
			name:                   "Crappy flow",
			header:                 "",
			expectedResponseStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
			req.Header.Add(testCase.header, "header-content")

			recorder := httptest.NewRecorder()
			s.handler.ServeHTTP(recorder, req)
			res := recorder.Result()
			s.Require().Equal(
				testCase.expectedResponseStatus,
				res.StatusCode,
				"Unexpected http response status code.",
			)

			if testCase.expectedResponseStatus != http.StatusOK {
				return
			}

			defer res.Body.Close()
			responseBodyContent, err := io.ReadAll(res.Body)
			s.Require().NoError(err)
			s.Equal(s.expectedResponseBody, string(responseBodyContent))
		},
		)
	}
}
