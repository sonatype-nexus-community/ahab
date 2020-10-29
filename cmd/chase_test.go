package cmd

import (
	"github.com/sonatype-nexus-community/ahab/packages"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChaseCommandNoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, chaseCmd.Use)
	assert.NotNil(t, err)
	assert.Equal(t, packages.SupportedPackageManagers, err.Error())
}
