package cmd

import (
	"encoding/json"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestChaseCommandNoArgs(t *testing.T) {
	// pass a specific package manager to avoid test behavior changes on different OSs.
	_, err := executeCommand(rootCmd, chaseCmd.Use, "--os", "apk")
	assert.NotNil(t, err)
	assert.Equal(t, MsgMissingStdIn, err.Error())
}

func TestChaseCommandApkInvalidStdInText(t *testing.T) {
	oldStdIn, tmpFile := createFakeStdIn(t)
	defer func() {
		os.Stdin = oldStdIn
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/", r.URL.EscapedPath())

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()
	origOssIndexURL := ossIndexURL
	ossIndexURL = ts.URL
	defer func() {
		ossIndexURL = origOssIndexURL
	}()

	_, err := executeCommand(rootCmd, chaseCmd.Use, "--os", "apk")
	assert.NotNil(t, err)
	assert.Equal(t, "An error occurred: [400 Bad Request] error accessing OSS Index", err.Error())
}

func TestChaseCommandBadUserAndToken(t *testing.T) {
	oldStdIn, tmpFile := createFakeStdInWithString(t, "alpine-baselayout-3.1.2-r0 - Alpine base dir structure and init scripts")
	defer func() {
		os.Stdin = oldStdIn
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/", r.URL.EscapedPath())

		assert.Equal(t, "Basic YmFkdXNlcjpiYWR0b2tlbg==", r.Header.Get("Authorization"), r)

		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()
	origOssIndexURL := ossIndexURL
	ossIndexURL = ts.URL
	defer func() {
		ossIndexURL = origOssIndexURL
	}()

	// pass a specific package manager to avoid test behavior changes on different OSs.
	_, err := executeCommand(rootCmd, chaseCmd.Use, "-u", "baduser", "-t", "badtoken", "--os", "apk")
	assert.NotNil(t, err)
	assert.Equal(t, "An error occurred: [401 Unauthorized] error accessing OSS Index", err.Error())
}

func TestChaseCommandEmptyUserAndToken(t *testing.T) {
	oldStdIn, tmpFile := createFakeStdInWithString(t, "alpine-baselayout-3.1.2-r0 - Alpine base dir structure and init scripts")
	defer func() {
		os.Stdin = oldStdIn
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/", r.URL.EscapedPath())

		assert.Equal(t, "", r.Header.Get("Authorization"), r)

		w.WriteHeader(http.StatusOK)
		coordinates := []types.Coordinate{expectedCoordinate}
		jsonCoordinates, _ := json.Marshal(coordinates)
		_, _ = w.Write(jsonCoordinates)
	}))
	defer ts.Close()
	origOssIndexURL := ossIndexURL
	ossIndexURL = ts.URL
	defer func() {
		ossIndexURL = origOssIndexURL
	}()

	// pass a specific package manager to avoid test behavior changes on different OSs.
	_, err := executeCommand(rootCmd, chaseCmd.Use, "-u", "", "-t", "", "--os", "apk")
	assert.Nil(t, err)
}

var expectedCoordinate types.Coordinate

func TestChaseCommandViperUserAndToken(t *testing.T) {
	oldStdIn, tmpFile := createFakeStdInWithString(t, "alpine-baselayout-3.1.2-r0 - Alpine base dir structure and init scripts")
	defer func() {
		os.Stdin = oldStdIn
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/", r.URL.EscapedPath())

		assert.Equal(t, "Basic b3NzaVVzZXJuYW1lVmFsdWU6b3NzaVRva2VuVmFsdWU=", r.Header.Get("Authorization"), r)

		w.WriteHeader(http.StatusOK)
		coordinates := []types.Coordinate{expectedCoordinate}
		jsonCoordinates, _ := json.Marshal(coordinates)
		_, _ = w.Write(jsonCoordinates)
	}))
	defer ts.Close()
	origOssIndexURL := ossIndexURL
	ossIndexURL = ts.URL
	defer func() {
		ossIndexURL = origOssIndexURL
	}()

	viper.Reset()
	defer viper.Reset()

	tempDir := setupConfig(t)
	defer resetConfig(t, tempDir)

	setupTestOSSIConfigFileValues(t, tempDir)
	defer func() {
		resetOSSIConfigFile()
	}()

	// pass a specific package manager to avoid test behavior changes on different OSs.
	_, err := executeCommand(rootCmd, chaseCmd.Use, "--os", "apk")
	assert.Nil(t, err)
}
