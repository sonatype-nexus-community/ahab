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
	"os"

	"github.com/sonatype-nexus-community/ahab/buildversion"
	"github.com/sonatype-nexus-community/ahab/logger"
	"github.com/sonatype-nexus-community/ahab/packages"
	"github.com/sonatype-nexus-community/go-sona-types/iq"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	iqUsername  string
	iqToken     string
	iqHost      string
	stage       string
	application string
	maxRetries  int
	lifecycle   *iq.Server
)

func init() {
	rootCmd.AddCommand(iqCmd)

	pf := iqCmd.PersistentFlags()
	pf.StringVar(&packageManager, "os", "", "Specify a value for the operating system type you want to scan (alpine, debian, fedora). Useful if autodetection fails and/or you want to explicitly set it.")
	pf.StringVar(&packageManager, "package-manager", "", "Specify package manager type you want to scan (apk, dnf, dpkg or yum). Useful if autodetection fails and/or you want to explicitly set it.")
	pf.StringVar(&iqUsername, "user", "admin", "Specify Nexus IQ Username for request")
	pf.StringVar(&iqToken, "token", "admin123", "Specify Nexus IQ Token/Password for request")
	pf.StringVar(&ossIndexUser, "oss-index-user", "", "Specify your OSS Index Username")
	pf.StringVar(&ossIndexToken, "oss-index-token", "", "Specify your OSS Index API Token")
	pf.StringVar(&iqHost, "host", "http://localhost:8070", "Specify Nexus IQ Server URL")
	pf.BoolVar(&quiet, "quiet", false, "Quiet removes the header from being printed")
	pf.StringVar(&application, "application", "", "Specify public application ID for request (required)")
	pf.StringVar(&stage, "stage", "develop", "Specify stage for application")
	pf.IntVar(&maxRetries, "max-retries", 300, "Specify maximum number of tries to poll Nexus IQ Server")
	pf.CountVarP(&verbose, "", "v", "Set log level, higher is more verbose")

	iqCmd.Flag("os").Deprecated = "use package-manager"
}

var iqCmd = &cobra.Command{
	Use:   "iq",
	Short: "iq is used for auditing your projects with Nexus IQ Server",
	Example: `
	dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab iq --application testapp
	yum list installed | ./ahab iq --application testapp
	dnf list installed | ./ahab iq --application testapp
	apk info -vv | sort | ./ahab iq	--application testapp
	`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if r := recover(); r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("pkg: %v", r)
				}
				_ = cmd.Usage()
				logger.PrintErrorAndLogLocation(err)
			}
		}()

		if !quiet {
			printHeader()
		}

		logLady, err = getLogger(verbose)
		if err != nil {
			panic(err)
		}

		fflags := cmd.Flags()

		err = checkRequiredFlags(fflags)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

		lifecycle = iq.New(logLady,
			iq.Options{
				User:          iqUsername,
				Token:         iqToken,
				Server:        iqHost,
				Application:   application,
				Stage:         stage,
				Tool:          "ahab-client",
				Version:       buildversion.BuildVersion,
				OSSIndexToken: ossIndexToken,
				OSSIndexUser:  ossIndexUser,
				DBCacheName:   "ahab-cache",
				MaxRetries:    maxRetries,
			})

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

		res, err := lifecycle.AuditPackages(purls, application)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

		fmt.Println()
		if res.IsError {
			logLady.WithField("res", res).Error("An error occurred with the request to IQ Server")
			panic(fmt.Errorf("Uh oh! There was an error with your request to Nexus IQ Server"))
		}

		if res.PolicyAction == "Failure" {
			logLady.WithField("res", res).Debug("Successful in communicating with IQ Server")
			fmt.Println("Ahoy, Ahab here matey, avast ye work, ye have some policy violations to clean up!")
			fmt.Println("Report URL: ", res.ReportHTMLURL)
			os.Exit(1)
			return
		}

		logLady.WithField("res", res).Debug("Successful in communicating with IQ Server")
		fmt.Println("Wonderbar! No policy violations reported for this audit!")
		fmt.Println("Report URL: ", res.ReportHTMLURL)
		return
	},
}

func checkRequiredFlags(flags *pflag.FlagSet) error {
	if !flags.Changed("application") {
		return fmt.Errorf("Application not set, see usage for more information")
	}
	return nil
}
