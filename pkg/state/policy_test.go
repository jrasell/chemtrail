package state

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestClientProvider_String(t *testing.T) {
	testCases := []struct {
		inputClientProvider ClientProvider
		expectedOutput      string
		name                string
	}{
		{
			inputClientProvider: AWSAutoScaling,
			expectedOutput:      "aws-autoscaling",
			name:                "AWS autoscaling provider",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputClientProvider.String()
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func TestClientProvider_Validate(t *testing.T) {

	const fakeProvider ClientProvider = "fake"

	testCases := []struct {
		inputClientProvider ClientProvider
		expectedOutput      error
		name                string
	}{
		{
			inputClientProvider: AWSAutoScaling,
			expectedOutput:      nil,
			name:                "AWS autoscaling provider",
		},
		{
			inputClientProvider: fakeProvider,
			expectedOutput:      errors.New("unsupported client provider \"fake\""),
			name:                "invalid scaling provider",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputClientProvider.Validate()
		if tc.expectedOutput == nil {
			assert.Nil(t, actualOutput, tc.name)
		} else {
			assert.EqualError(t, actualOutput, tc.expectedOutput.Error(), tc.name)
		}
	}
}
