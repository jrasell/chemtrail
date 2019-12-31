package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapStringsToSliceString(t *testing.T) {
	input := map[string]string{"foo": "bar", "this": "that", "life": "death"}
	expectedOut := []string{"foo:bar", "this:that", "life:death"}
	actualOut := MapStringsToSliceString(input, ":")
	assert.Equal(t, expectedOut, actualOut)
}
