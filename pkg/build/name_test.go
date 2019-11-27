package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProgramName(t *testing.T) {
	programName = "ChemtrailProgramNameTest"
	returnProgramName := ProgramName()
	assert.Equal(t, programName, returnProgramName)
}

func Test_SetProgramName(t *testing.T) {
	SetProgramName("ChemtrailSetProgramNameTest")
	assert.Equal(t, programName, "ChemtrailSetProgramNameTest")
}
