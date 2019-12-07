package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_StorageConfig(t *testing.T) {
	fakeCMD := &cobra.Command{}
	RegisterStorageConfig(fakeCMD)

	cfg := GetStorageConfig()
	assert.Equal(t, false, cfg.ConsulEnabled)
	assert.Equal(t, "chemtrail/", cfg.ConsulPath)
}
