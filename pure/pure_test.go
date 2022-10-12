package pure

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// PureTestSuite declares global variables, used as constants in all tests.
type PureTestSuite struct {
	suite.Suite
	configPrefix string
}

// SetupSuite sets the global state once, for all tests.
func (s *PureTestSuite) SetupSuite() {
	s.configPrefix = "configuration: "
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(PureTestSuite))
}

func (s *PureTestSuite) TestMyPureFunction() {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		expectedErr    error
	}{
		{
			name:           "Happy flow",
			input:          "my input\n",
			expectedOutput: "configuration: my input\n",
			expectedErr:    nil,
		},
		{
			name:           "Empty input",
			input:          "",
			expectedOutput: "",
			expectedErr:    ErrEmptyInput,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			output, err := functionUnderTest(s.configPrefix, testCase.input)
			if testCase.expectedErr != nil {
				s.ErrorIs(err, testCase.expectedErr)
				return
			}
			s.Require().NoError(err)
			s.Equal(output, testCase.expectedOutput)
		},
		)
	}
}
