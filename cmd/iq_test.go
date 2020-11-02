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
	"fmt"
	"github.com/sonatype-nexus-community/go-sona-types/configuration"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestIqApplicationFlagMissing(t *testing.T) {
	_, err := executeCommand(rootCmd, iqCmd.Use)
	//output, err := executeCommand(rootCmd, iqCmd.Use)
	//assert.Contains(t, output, "Error: \""+flagNameIqApplication+"\" not set, see usage for more information")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "\""+flagNameIqApplication+"\" not set, see usage for more information")
}

func TestIqHelp(t *testing.T) {
	output, err := executeCommand(rootCmd, iqCmd.Use, "--help")
	assert.Contains(t, output, "\tyum list installed | ./ahab iq --"+flagNameIqApplication+" testapp")
	assert.Nil(t, err)
}

func setupIQConfigFile(t *testing.T, tempDir string) {
	cfgDirIQ := path.Join(tempDir, types.IQServerDirName)
	assert.Nil(t, os.Mkdir(cfgDirIQ, 0700))

	cfgFileIQ = path.Join(tempDir, types.IQServerDirName, types.IQServerConfigFileName)
}
func resetIQConfigFile() {
	cfgFileIQ = ""
}

func TestInitIQConfig(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	tempDir := setupConfig(t)
	defer resetConfig(t, tempDir)

	setupTestOSSIConfigFileValues(t, tempDir)
	defer func() {
		resetOSSIConfigFile()
	}()

	setupIQConfigFile(t, tempDir)
	defer func() {
		resetIQConfigFile()
	}()

	credentials := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n",
		configuration.ViperKeyIQUsername, "iqUsernameValue",
		configuration.ViperKeyIQToken, "iqTokenValue",
		configuration.ViperKeyIQServer, "iqServerValue")
	assert.Nil(t, ioutil.WriteFile(cfgFileIQ, []byte(credentials), 0644))

	// init order is not guaranteed
	initIQConfig()
	initConfig()

	// verify the OSSI stuff, since we will call both OSSI and IQ
	assert.Equal(t, "ossiUsernameValue", viper.GetString(configuration.ViperKeyUsername))
	assert.Equal(t, "ossiTokenValue", viper.GetString(configuration.ViperKeyToken))
	// verify the IQ stuff
	assert.Equal(t, "iqUsernameValue", viper.GetString(configuration.ViperKeyIQUsername))
	assert.Equal(t, "iqTokenValue", viper.GetString(configuration.ViperKeyIQToken))
	assert.Equal(t, "iqServerValue", viper.GetString(configuration.ViperKeyIQServer))
}

func TestInitIQConfigWithNoConfigFile(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	tempDir := setupConfig(t)
	defer resetConfig(t, tempDir)

	setupTestOSSIConfigFileValues(t, tempDir)
	defer func() {
		resetOSSIConfigFile()
	}()

	setupIQConfigFile(t, tempDir)
	defer func() {
		resetIQConfigFile()
	}()
	credentials := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n",
		configuration.ViperKeyIQUsername, "iqUsernameValue",
		configuration.ViperKeyIQToken, "iqTokenValue",
		configuration.ViperKeyIQServer, "iqServerValue")
	assert.Nil(t, ioutil.WriteFile(cfgFileIQ, []byte(credentials), 0644))

	// delete the config files
	assert.NoError(t, os.Remove(cfgFile))
	assert.NoError(t, os.Remove(cfgFileIQ))

	// init order is not guaranteed
	initIQConfig()
	initConfig()

	// verify the OSSI stuff, since we will call both OSSI and IQ
	assert.Equal(t, "", viper.GetString(configuration.ViperKeyUsername))
	assert.Equal(t, "", viper.GetString(configuration.ViperKeyToken))
	// verify the IQ stuff
	assert.Equal(t, "", viper.GetString(configuration.ViperKeyIQUsername))
	assert.Equal(t, "", viper.GetString(configuration.ViperKeyIQToken))
	assert.Equal(t, "", viper.GetString(configuration.ViperKeyIQServer))
}
