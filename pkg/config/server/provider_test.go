package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ProviderConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterProviderConfig(fakeCMD)

	cfg := GetProviderConfig()
	assert.Equal(t, false, cfg.AWSASG)
}
