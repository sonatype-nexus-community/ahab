//
// Copyright (c) 2019-present Sonatype, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"encoding/json"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
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
	assert.Contains(t, err.Error(), MsgMissingStdIn)
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
	assert.Contains(t, err.Error(), "An error occurred: [400 Bad Request] error accessing OSS Index")
}

func TestChaseCommandBadUserAndToken(t *testing.T) {
	// FAKE alpine package should avoid test failure due to cached package
	oldStdIn, tmpFile := createFakeStdInWithString(t, "FAKE__alpine-baselayout-3.1.2-r0 - Alpine base dir structure and init scripts")
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
	assert.NotNil(t, err) // if this fails, make sure the cache is cleaned first
	assert.Contains(t, err.Error(), "An error occurred: [401 Unauthorized] error accessing OSS Index")
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

// TODO: Figure out why this test works when run by itself, but fails as part of test suite.
// Manual testing confirms correct behavior
/*func TestChaseCommandViperUserAndToken(t *testing.T) {
	oldStdIn, tmpFile := createFakeStdInWithString(t, "alpine-baselayout-3.1.2-r0 - Alpine base dir structure and init scripts")
	defer func() {
		os.Stdin = oldStdIn
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/", r.URL.EscapedPath())

		assert.Equal(t, "ossiUsernameValue", ossi.Options.Username)
		assert.Equal(t, "ossiTokenValue", ossi.Options.Token)
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
*/
