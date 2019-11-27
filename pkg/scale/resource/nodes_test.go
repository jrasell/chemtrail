package resource

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
)

func Test_updateHandler_checkNodeClass(t *testing.T) {
	testCases := []struct {
		inputNode          *api.Node
		expectedClassValue string
		name               string
	}{
		{
			inputNode:          &api.Node{NodeClass: ""},
			expectedClassValue: "chemtrail-default",
			name:               "empty node class",
		},
		{
			inputNode:          &api.Node{NodeClass: "high-memory"},
			expectedClassValue: "high-memory",
			name:               "operator configured node class",
		},
	}
	uh := &updateHandler{}

	for _, tc := range testCases {
		uh.checkNodeClass(tc.inputNode)
		assert.Equal(t, tc.expectedClassValue, tc.inputNode.NodeClass, tc.name)
	}
}
