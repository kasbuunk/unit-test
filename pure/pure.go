package pure

import "errors"

var ErrEmptyInput = errors.New("input and config must not be empty")

func functionUnderTest(configVar, input string) (string, error) {
	if input == "" || configVar == "" {
		return "", ErrEmptyInput
	}
	output := configVar + input
	return output, nil
}
