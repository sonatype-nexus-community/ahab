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
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sonatype-nexus-community/go-sona-types/configuration"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
	"github.com/spf13/viper"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose int
	logLady *logrus.Logger
)

var rootCmd = &cobra.Command{
	Use:   "ahab",
	Short: "ahab is a tool for scanning linux OS packages for vulnerabilities",
	Run: func(cmd *cobra.Command, args []string) {
		printHeader()
		_ = cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&ossIndexUser, flagNameOssiUsername, "u", "", "Specify your OSS Index Username")
	pf.StringVarP(&ossIndexToken, flagNameOssiToken, "t", "", "Specify your OSS Index API Token")
}

func initConfig() {
	var cfgFileToCheck string
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType(configuration.ConfigTypeYaml)
		cfgFileToCheck = cfgFile
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configPath := path.Join(home, types.OssIndexDirName)

		viper.AddConfigPath(configPath)
		viper.SetConfigType(configuration.ConfigTypeYaml)
		viper.SetConfigName(types.OssIndexConfigFileName)

		cfgFileToCheck = path.Join(configPath, types.OssIndexConfigFileName)
	}

	if fileExists(cfgFileToCheck) {
		// 'merge' OSSI config here, since IQ cmd also need OSSI config, and init order is not guaranteed
		if err := viper.MergeInConfig(); err != nil {
			panic(err)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func cleanUserName(origUsername string) string {
	runes := []rune(origUsername)
	cleanUsername := "***hidden***"
	if len(runes) > 0 {
		first := string(runes[0])
		last := string(runes[len(runes)-1])
		cleanUsername = first + "***hidden***" + last
	}
	return cleanUsername
}
