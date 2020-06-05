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
	"bufio"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sonatype-nexus-community/ahab/buildversion"
	"github.com/sonatype-nexus-community/ahab/packages"
	"github.com/sonatype-nexus-community/ahab/parse"
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

var iqCmd = &cobra.Command{
	Use:   "iq",
	Short: "iq is used for auditing your projects with Nexus IQ Server",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fflags := cmd.Flags()

		err = checkRequiredFlags(fflags)
		if err != nil {
			return
		}

		printHeader()
		logger = logrus.New()
		switch verbose {
		case 1:
			logger.Level = logrus.InfoLevel
		case 2:
			logger.Level = logrus.DebugLevel
		case 3:
			logger.Level = logrus.TraceLevel
		default:
			logger.Level = logrus.ErrorLevel
		}

		lifecycle = iq.New(logger,
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

		var list []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			list = append(list, scanner.Text())
		}
		if err = scanner.Err(); err != nil {
			if err != nil {
				return
			}
		}

		res, err := lifecycle.AuditPackages(tryParseStdInListIQ(list, &operating), application)
		if err != nil {
			return
		}

		fmt.Println()
		if res.IsError {
			logger.WithField("res", res).Error("An error occurred with the request to IQ Server")
			return fmt.Errorf("Uh oh! There was an error with your request to Nexus IQ Server")
		}

		if res.PolicyAction != "Failure" {
			logger.WithField("res", res).Debug("Successful in communicating with IQ Server")
			fmt.Println("Wonderbar! No policy violations reported for this audit!")
			fmt.Println("Report URL: ", res.ReportHTMLURL)
			os.Exit(0)
		} else {
			logger.WithField("res", res).Debug("Successful in communicating with IQ Server")
			fmt.Println("Hi, Nancy here, you have some policy violations to clean up!")
			fmt.Println("Report URL: ", res.ReportHTMLURL)
			os.Exit(1)
		}
		return
	},
}

func tryParseStdInListIQ(list []string, operating *string) (purls []string) {
	var thing string
	thing = *operating
	switch thing {
	case "debian":
		logger.Trace("Chasing Debian")
		var aptResult packages.Apt
		aptResult.ProjectList = parse.ParseDpkgList(list)
		purls = aptResult.ExtractPurlsFromProjectList(*operating)
	case "alpine":
		logger.Trace("Chasing Alpine")
		var apkResult packages.Apk
		apkResult.ProjectList = parse.ParseApkShow(list)
		purls = apkResult.ExtractPurlsFromProjectList(*operating)
	default:
		logger.Trace("Chasing Yum")
		var yumResult packages.Yum
		yumResult.ProjectList = parse.ParseYumListFromStdIn(list)
		purls = yumResult.ExtractPurlsFromProjectList(*operating)
	}
	return
}

func checkRequiredFlags(flags *pflag.FlagSet) error {
	if !flags.Changed("application") {
		return fmt.Errorf("Application not set, see usage for more information")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(iqCmd)

	pf := iqCmd.PersistentFlags()
	pf.StringVar(&operating, "os", "debian", "")
	pf.StringVar(&iqUsername, "user", "admin", "Specify Nexus IQ username for request")
	pf.StringVar(&iqToken, "token", "admin123", "Specify Nexus IQ token/password for request")
	pf.StringVar(&ossIndexUser, "oss-index-user", "", "")
	pf.StringVar(&ossIndexToken, "oss-index-token", "", "")
	pf.StringVar(&iqHost, "host", "http://localhost:8070", "Specify Nexus IQ Server URL")
	pf.StringVar(&application, "application", "", "Specify public application ID for request (required)")
	pf.StringVar(&stage, "stage", "develop", "Specify stage for application")
	pf.BoolVar(&noColor, "no-color", false, "")
	pf.IntVar(&maxRetries, "max-retries", 300, "Specify maximum number of tries to poll Nexus IQ Server")
	pf.CountVarP(&verbose, "", "v", "Set log level, higher is more verbose")
}
