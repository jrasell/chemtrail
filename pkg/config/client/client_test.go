package client

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ChemtrailConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterConfig(fakeCMD)
	assert.Equal(t, configKeyChemtrailAddrDefault, GetConfig().Addr)
}
