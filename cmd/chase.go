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
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/sonatype-nexus-community/ahab/audit"
	"github.com/sonatype-nexus-community/ahab/buildversion"
	"github.com/sonatype-nexus-community/ahab/logger"
	"github.com/sonatype-nexus-community/ahab/packages"
	"github.com/sonatype-nexus-community/ahab/parse"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
	"github.com/spf13/cobra"
)

var (
	operating     string
	cleanCache    bool
	ossIndexUser  string
	ossIndexToken string
	output        string
	loud          bool
	quiet         bool
	noColor       bool
	ossi          *ossindex.Server
)

func init() {
	rootCmd.AddCommand(chaseCmd)
	chaseCmd.PersistentFlags().StringVar(&operating, "os", "debian", "")
	chaseCmd.PersistentFlags().BoolVar(&cleanCache, "clean-cache", false, "")
	chaseCmd.PersistentFlags().StringVar(&ossIndexUser, "user", "", "")
	chaseCmd.PersistentFlags().StringVar(&ossIndexToken, "token", "", "")
	chaseCmd.PersistentFlags().StringVar(&output, "output", "text", "")
	chaseCmd.PersistentFlags().BoolVar(&loud, "loud", false, "")
	chaseCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "")
	chaseCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "")
	chaseCmd.PersistentFlags().CountVarP(&verbose, "", "v", "Set log level, higher is more verbose")
}

var chaseCmd = &cobra.Command{
	Use:   "chase",
	Short: "chase is used for auditing projects with OSS Index",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		logLady, err = getLogger(verbose)
		if err != nil {
			return
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

		if output == "text" && !quiet {
			printHeader()
		}

		if cleanCache {
			err = ossi.NoCacheNoProblems()
			if err != nil {
				logLady.Error(err)
				return
			}
			logLady.Trace("Cleaned Ahab Cache")
			return
		}

		logLady.Trace("Attempting to audit list of strings from standard in")
		pkgs, err := parseStdIn(&operating)
		if err != nil {
			logLady.Error(err)
			return
		}

		logLady.Trace("Attempting to extract purls from Project List")
		purls := pkgs.ExtractPurlsFromProjectList(operating)

		logLady.Trace("Attempting to Audit Packages with OSS Index")
		coordinates, err := ossi.AuditPackages(purls)
		if err != nil {
			logLady.Error(err)
			return
		}

		logLady.Trace("Attempting to output audited packages results")
		count, results := audit.LogResults(quiet, noColor, loud, output, coordinates)
		fmt.Print(results)
		if count > 0 {
			os.Exit(1)
		}

		return
	},
}

func getLogger(level int) (*logrus.Logger, error) {
	switch level {
	case 1:
		return logger.GetLogger(logrus.InfoLevel)
	case 2:
		return logger.GetLogger(logrus.DebugLevel)
	case 3:
		return logger.GetLogger(logrus.TraceLevel)
	default:
		return logger.GetLogger(logrus.ErrorLevel)
	}
}

func parseStdInList(list []string, operating *string) (packages.IPackage, error) {
	var thing string
	thing = *operating
	switch thing {
	case "debian":
		logLady.Trace("Chasing Debian")
		var aptResult packages.Apt
		aptResult.ProjectList = parse.ParseDpkgList(list)
		return aptResult, nil
	case "alpine":
		logLady.Trace("Chasing Alpine")
		var apkResult packages.Apk
		apkResult.ProjectList = parse.ParseApkShow(list)
		return apkResult, nil
	default:
		logLady.Trace("Chasing Yum")
		var yumResult packages.Yum
		yumResult.ProjectList = parse.ParseYumListFromStdIn(list)
		return yumResult, nil
	}
}

func parseStdIn(operating *string) (packages.IPackage, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if (fi.Mode() & os.ModeNamedPipe) == 0 {
		return nil, fmt.Errorf("Nothing passed in to Standard In")
	}

	var list []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != nil {
			return nil, err
		}
	}
	return parseStdInList(list, operating)
}

func printHeader() {
	figure.NewFigure("Ahab", "larry3d", true).Print()
	figure.NewFigure("By Sonatype & Friends", "pepper", true).Print()
	fmt.Println("Ahab version: " + buildversion.BuildVersion)
}
