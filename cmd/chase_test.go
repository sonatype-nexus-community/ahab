package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChaseCommandNoArgs(t *testing.T) {
	// pass a specific package manager to avoid test behavior changes on different OSs.
	_, err := executeCommand(rootCmd, chaseCmd.Use, "--os", "apk")
	assert.NotNil(t, err)
	assert.Equal(t, MsgMissingStdIn, err.Error())
}
