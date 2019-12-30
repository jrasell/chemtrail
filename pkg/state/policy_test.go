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
		{
			inputClientProvider: NoOpClientProvider,
			expectedOutput:      "no-op",
			name:                "no-op provider",
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
			inputClientProvider: NoOpClientProvider,
			expectedOutput:      nil,
			name:                "no-op provider",
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

func TestComparisonAction_String(t *testing.T) {
	testCases := []struct {
		inputComparisonAction ComparisonAction
		expectedOutput        string
		name                  string
	}{
		{
			inputComparisonAction: ActionScaleIn,
			expectedOutput:        "scale-in",
			name:                  "scale in comparison action",
		},
		{
			inputComparisonAction: ActionScaleOut,
			expectedOutput:        "scale-out",
			name:                  "scale out comparison action",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputComparisonAction.String()
		assert.Equal(t, tc.expectedOutput, actualOutput, tc.name)
	}
}

func TestComparisonAction_Validate(t *testing.T) {

	const fakeAction ComparisonAction = "fake"

	testCases := []struct {
		inputComparisonAction ComparisonAction
		expectedOutput        error
		name                  string
	}{
		{
			inputComparisonAction: ActionScaleIn,
			expectedOutput:        nil,
			name:                  "scale-in",
		},
		{
			inputComparisonAction: ActionScaleOut,
			expectedOutput:        nil,
			name:                  "scale-out",
		},
		{
			inputComparisonAction: fakeAction,
			expectedOutput:        errors.New("Action \"fake\" is not a valid option"),
			name:                  "invalid comparison action",
		},
	}

	for _, tc := range testCases {
		actualOutput := tc.inputComparisonAction.Validate()
		if tc.expectedOutput == nil {
			assert.Nil(t, actualOutput, tc.name)
		} else {
			assert.EqualError(t, actualOutput, tc.expectedOutput.Error(), tc.name)
		}
	}
}
