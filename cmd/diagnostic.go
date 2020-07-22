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
	"time"

	"github.com/sonatype-nexus-community/ahab/buildversion"
	"github.com/sonatype-nexus-community/ahab/logger"
	"github.com/sonatype-nexus-community/go-sona-types/cyclonedx"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
	"github.com/spf13/cobra"
)

var diagnosticCmd = &cobra.Command{
	Use:   "diagnostic",
	Short: "diagnostic is used for generating diagnostic things for support",
	Example: `
	./ahab diagnostic --os debian
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

		logLady, err = getLogger(verbose)
		if err != nil {
			panic(err)
		}

		ossi = ossindex.New(logLady,
			types.Options{
				DBCacheName: "ahab-cache",
				TTL:         time.Now().Local().Add(time.Hour * 12),
				Tool:        "ahab-client",
				Version:     buildversion.BuildVersion,
				Username:    ossIndexUser,
				Token:       ossIndexToken,
			})

		logLady.Trace("Attempting to audit list of strings from standard in")
		pkgs, err := parseStdIn(&operating)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

		logLady.Trace("Attempting to extract purls from Project List")
		purls := pkgs.ExtractPurlsFromProjectList(operating)

		logLady.Trace("Attempting to Audit Packages with OSS Index")
		coordinates, err := ossi.AuditPackages(purls)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

		dx := cyclonedx.Default(logLady)

		sbom := dx.FromCoordinates(coordinates)

		fmt.Print(sbom)

		return
	},
}

func init() {
	rootCmd.AddCommand(diagnosticCmd)
	diagnosticCmd.PersistentFlags().StringVar(&operating, "os", "debian", "Specify a value for the operating system type you want to scan (alpine, debian, fedora)")
	diagnosticCmd.PersistentFlags().CountVarP(&verbose, "", "v", "Set log level, higher is more verbose")
}
