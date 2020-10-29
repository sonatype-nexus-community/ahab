//
// Copyright 2018-present Sonatype Inc.
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
	"bytes"
	"github.com/sonatype-nexus-community/go-sona-types/configuration"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()

	return buf.String(), err
}

func TestRootCommandNoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "")
	assert.Nil(t, err)
}

func TestRootCommandUnknownCommand(t *testing.T) {
	output, err := executeCommand(rootCmd, "one", "two")
	assert.Contains(t, output, "Error: unknown command \"one\" for \"ahab\"")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unknown command \"one\" for \"ahab\"")
}

func setupConfig(t *testing.T) (tempDir string) {
	tempDir, err := ioutil.TempDir("", "config-test")
	assert.NoError(t, err)
	return tempDir
}

func resetConfig(t *testing.T, tempDir string) {
	var err error
	assert.NoError(t, err)
	_ = os.RemoveAll(tempDir)
}

func TestInitConfig(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	tempDir := setupConfig(t)
	defer resetConfig(t, tempDir)

	setupTestOSSIConfigFileValues(t, tempDir)
	defer func() {
		resetOSSIConfigFile()
	}()

	initConfig()

	assert.Equal(t, "ossiUsernameValue", viper.GetString(configuration.ViperKeyUsername))
	assert.Equal(t, "ossiTokenValue", viper.GetString(configuration.ViperKeyToken))
}

func TestInitConfigWithNoConfigFile(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	tempDir := setupConfig(t)
	defer resetConfig(t, tempDir)

	setupTestOSSIConfigFileValues(t, tempDir)
	defer func() {
		resetOSSIConfigFile()
	}()
	// delete the config file
	assert.NoError(t, os.Remove(cfgFile))

	initConfig()

	assert.Equal(t, "", viper.GetString(configuration.ViperKeyUsername))
	assert.Equal(t, "", viper.GetString(configuration.ViperKeyToken))
}

func setupTestOSSIConfigFile(t *testing.T, tempDir string) {
	cfgDir := path.Join(tempDir, types.OssIndexDirName)
	assert.Nil(t, os.Mkdir(cfgDir, 0700))

	cfgFile = path.Join(tempDir, types.OssIndexDirName, types.OssIndexConfigFileName)
}

func resetOSSIConfigFile() {
	cfgFile = ""
}

func setupTestOSSIConfigFileValues(t *testing.T, tempDir string) {
	setupTestOSSIConfigFile(t, tempDir)

	const credentials = configuration.ViperKeyUsername + ": ossiUsernameValue\n" +
		configuration.ViperKeyToken + ": ossiTokenValue"
	assert.Nil(t, ioutil.WriteFile(cfgFile, []byte(credentials), 0644))
}

func createFakeStdIn(t *testing.T) (oldStdIn *os.File, tmpFile *os.File) {
	return createFakeStdInWithString(t, "Testing")
}
func createFakeStdInWithString(t *testing.T, inputString string) (oldStdIn *os.File, tmpFile *os.File) {
	content := []byte(inputString)
	tmpFile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		t.Error(err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		t.Error(err)
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Error(err)
	}

	oldStdIn = os.Stdin

	os.Stdin = tmpFile
	return oldStdIn, tmpFile
}
