//
// Copyright 2019-Present Sonatype Inc.
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
	"github.com/sonatype-nexus-community/ahab/internal/customerrors"
	"github.com/sonatype-nexus-community/go-sona-types/configuration"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"

	"github.com/sonatype-nexus-community/ahab/buildversion"
	"github.com/sonatype-nexus-community/ahab/packages"
	"github.com/sonatype-nexus-community/go-sona-types/iq"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	cfgFileIQ   string
	iqUsername  string
	iqToken     string
	iqHost      string
	stage       string
	application string
	maxRetries  int
	lifecycle   *iq.Server
)

const (
	flagNameIqUsername    = "iq-username"
	flagNameIqToken       = "iq-token"
	flagNameIqStage       = "iq-stage"
	flagNameIqApplication = "iq-application"
	flagNameIqServerUrl   = "iq-server-url"
)

func init() {
	cobra.OnInitialize(initIQConfig)

	pf := iqCmd.PersistentFlags()
	pf.StringVar(&packageManager, "os", "", "Specify a value for the operating system type you want to scan (alpine, debian, fedora). Useful if autodetection fails and/or you want to explicitly set it.")
	pf.StringVar(&packageManager, "package-manager", "", "Specify package manager type you want to scan (apk, dnf, dpkg or yum). Useful if autodetection fails and/or you want to explicitly set it.")
	pf.StringVarP(&iqUsername, flagNameIqUsername, "l", "admin", "Specify Nexus IQ Username for request")
	pf.StringVarP(&iqToken, flagNameIqToken, "k", "admin123", "Specify Nexus IQ Token/Password for request")
	pf.StringVarP(&iqHost, flagNameIqServerUrl, "x", "http://localhost:8070", "Specify Nexus IQ Server URL")
	pf.BoolVar(&quiet, "quiet", true, "Quiet removes the header from being printed")

	pf.StringVarP(&application, flagNameIqApplication, "a", "", "Specify public application ID for request (required)")
	if err := iqCmd.MarkPersistentFlagRequired(flagNameIqApplication); err != nil {
		panic(err)
	}

	pf.StringVarP(&stage, flagNameIqStage, "s", "develop", "Specify stage for application")
	pf.IntVar(&maxRetries, "max-retries", 300, "Specify maximum number of tries to poll Nexus IQ Server")
	pf.CountVarP(&verbose, "", "v", "Set log level, higher is more verbose")

	iqCmd.Flag("os").Deprecated = "use package-manager"

	rootCmd.AddCommand(iqCmd)
}

var iqCmd = &cobra.Command{
	Use:   "iq",
	Short: "iq is used for auditing your projects with Nexus IQ Server",
	Example: `
	dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab iq --` + flagNameIqApplication + ` testapp
	yum list installed | ./ahab iq --` + flagNameIqApplication + ` testapp
	dnf list installed | ./ahab iq --` + flagNameIqApplication + ` testapp
	apk info -vv | sort | ./ahab iq	--` + flagNameIqApplication + ` testapp
	`,
	PreRun: func(cmd *cobra.Command, args []string) { bindViperIq(cmd) },
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if r := recover(); r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("pkg: %v", r)
				}
				err = customerrors.ErrorShowLogPath{Err: err}
			}
		}()

		if !quiet {
			printHeader()
		}

		logLady, err = getLogger(verbose)
		if err != nil {
			panic(err)
		}

		lifecycle, err = iq.New(logLady,
			iq.Options{
				User:          viper.GetString(configuration.ViperKeyIQUsername),
				Token:         viper.GetString(configuration.ViperKeyIQToken),
				Server:        viper.GetString(configuration.ViperKeyIQServer),
				Application:   application,
				Stage:         stage,
				Tool:          "ahab-client",
				Version:       buildversion.BuildVersion,
				OSSIndexToken: viper.GetString(configuration.ViperKeyToken),
				OSSIndexUser:  viper.GetString(configuration.ViperKeyUsername),
				DBCacheName:   "ahab-cache",
				MaxRetries:    maxRetries,
			})
		if err != nil {
			panic(err)
		}

		logLady.WithField("lifecycle", iq.Options{
			User:          cleanUserName(lifecycle.Options.User),
			Token:         "***hidden***",
			Server:        lifecycle.Options.Server,
			Application:   lifecycle.Options.Application,
			Stage:         lifecycle.Options.Stage,
			Tool:          lifecycle.Options.Tool,
			Version:       lifecycle.Options.Version,
			OSSIndexUser:  cleanUserName(lifecycle.Options.OSSIndexUser),
			OSSIndexToken: "***hidden***",
			DBCacheName:   lifecycle.Options.DBCacheName,
			MaxRetries:    lifecycle.Options.MaxRetries,
		}).Debug("Created iq server")

		if packageManager == "" {
			logLady.Trace("Attempting to detect package manager for you")
			manager, err := packages.DetectPackageManager(logLady)
			if err != nil {
				logLady.Error(err)
				panic(err)
			}
			packageManager = manager
		}

		pkgs, err := parseStdIn(&packageManager)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

		purls := pkgs.ExtractPurlsFromProjectList()

		res, err := lifecycle.AuditPackages(purls)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

		fmt.Println()
		if res.IsError {
			logLady.WithField("res", res).Error("An error occurred with the request to IQ Server")
			panic(fmt.Errorf("Uh oh! There was an error with your request to Nexus IQ Server"))
		}

		showPolicyActionMessage(res, os.Stdout)
		if res.PolicyAction == iq.PolicyActionFailure {
			os.Exit(1)
			return
		}
		return
	},
}

func showPolicyActionMessage(res iq.StatusURLResult, writer io.Writer) {
	switch res.PolicyAction {
	case iq.PolicyActionFailure:
		logLady.WithField("res", res).Debug("Successful in communicating with IQ Server")
		_, _ = fmt.Fprintln(writer, "Ahoy, Ahab here matey, avast ye work, ye have some policy violations to clean up!")
		_, _ = fmt.Fprintln(writer, "Report URL: ", res.AbsoluteReportHTMLURL)
	case iq.PolicyActionWarning:
		logLady.WithField("res", res).Debug("Successful in communicating with IQ Server")
		_, _ = fmt.Fprintln(writer, "A shot across the bow, there be policy warnings!")
		_, _ = fmt.Fprintln(writer, "Report URL: ", res.AbsoluteReportHTMLURL)
	default:
		logLady.WithField("res", res).Debug("Successful in communicating with IQ Server")
		_, _ = fmt.Fprintln(writer, "Wonderbar! No policy violations reported for this audit!")
		_, _ = fmt.Fprintln(writer, "Report URL: ", res.AbsoluteReportHTMLURL)
	}
}

func bindViperIq(cmd *cobra.Command) {
	// need to defer bind call until command is run. see: https://github.com/spf13/viper/issues/233

	// need to ensure ossi CLI flags will override ossi config file values when running IQ command
	bindViperRootCmd()

	// Bind viper to the flags passed in via the command line, so it will override config from file
	if err := viper.BindPFlag(configuration.ViperKeyIQUsername, lookupFlagNotNil(flagNameIqUsername, cmd)); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag(configuration.ViperKeyIQToken, lookupFlagNotNil(flagNameIqToken, cmd)); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag(configuration.ViperKeyIQServer, lookupFlagNotNil(flagNameIqServerUrl, cmd)); err != nil {
		panic(err)
	}
}

func lookupFlagNotNil(flagName string, cmd *cobra.Command) *pflag.Flag {
	// see: https://github.com/spf13/viper/pull/949
	foundFlag := cmd.Flags().Lookup(flagName)
	if foundFlag == nil {
		panic(fmt.Errorf("flag lookup for name: '%s' returned nil", flagName))
	}
	return foundFlag
}

func initIQConfig() {
	var cfgFileToCheck string
	if cfgFileIQ != "" {
		viper.SetConfigFile(cfgFileIQ)
		viper.SetConfigType(configuration.ConfigTypeYaml)
		cfgFileToCheck = cfgFileIQ
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configPath := path.Join(home, types.IQServerDirName)

		viper.AddConfigPath(configPath)
		viper.SetConfigType(configuration.ConfigTypeYaml)
		viper.SetConfigName(types.IQServerConfigFileName)

		cfgFileToCheck = path.Join(configPath, types.IQServerConfigFileName)
	}

	if fileExists(cfgFileToCheck) {
		// 'merge' IQ config here, since we also need OSSI config, and load order is not guaranteed
		if err := viper.MergeInConfig(); err != nil {
			panic(err)
		}
	}
}
