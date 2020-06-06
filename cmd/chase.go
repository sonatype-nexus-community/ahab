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
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/sonatype-nexus-community/ahab/audit"
	"github.com/sonatype-nexus-community/ahab/buildversion"
	"github.com/sonatype-nexus-community/ahab/packages"
	"github.com/sonatype-nexus-community/ahab/parse"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
	"github.com/spf13/cobra"
)

var (
	whales        string
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

var chaseCmd = &cobra.Command{
	Use:   "chase",
	Short: "chase is used for auditing projects with OSS Index",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
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

		ossi = ossindex.New(logger,
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
				logger.Error(err)
				return
			}
			logger.Trace("Cleaned Ahab Cache")
			return
		}

		if whales != "" {
			logger.WithField("packages", whales).Trace("Attempting to audit list of strings from command line")
			tryParseFlag(&whales, &operating)
		} else {
			logger.Trace("Attempting to audit list of strings from standard in")
			tryParseStdIn(&operating)
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(chaseCmd)
	chaseCmd.PersistentFlags().StringVar(&whales, "whales", "", "")
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

func tryParseFlag(flag *string, operating *string) {
	whales := strings.Split(*flag, ",")
	var aptResult packages.Apt
	aptResult.ProjectList = parse.ParseAptList(whales)
	var thing string
	thing = *operating
	tryExtractAndAudit(aptResult, thing)
}

func tryExtractAndAudit(pkgs packages.IPackage, operating string) {
	purls := pkgs.ExtractPurlsFromProjectList(operating)
	tryAuditPackages(purls, len(purls))
}

func tryParseStdInList(list []string, operating *string) {
	var thing string
	thing = *operating
	switch thing {
	case "debian":
		logger.Trace("Chasing Debian")
		var aptResult packages.Apt
		aptResult.ProjectList = parse.ParseDpkgList(list)
		tryExtractAndAudit(aptResult, thing)
	case "alpine":
		logger.Trace("Chasing Alpine")
		var apkResult packages.Apk
		apkResult.ProjectList = parse.ParseApkShow(list)
		tryExtractAndAudit(apkResult, thing)
	default:
		logger.Trace("Chasing Yum")
		var yumResult packages.Yum
		yumResult.ProjectList = parse.ParseYumListFromStdIn(list)
		tryExtractAndAudit(yumResult, thing)
	}
}

func tryAuditPackages(purls []string, count int) {
	coordinates, err := ossi.AuditPackages(purls)
	if err != nil {
		logger.Error(err)
	}
	logger.Trace(coordinates)
	count, results := audit.LogResults(quiet, noColor, loud, output, coordinates)
	logger.Trace(count)
	logger.Trace(results)
	fmt.Print(results)
}

func tryParseStdIn(operating *string) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if (fi.Mode() & os.ModeNamedPipe) == 0 {
		os.Exit(1)
	} else {
		doRead(operating)
	}
}

func doRead(operating *string) {
	var list []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != nil {
			panic(err)
		}
	}
	tryParseStdInList(list, operating)
}

func printHeader() {
	figure.NewFigure("Ahab", "larry3d", true).Print()
	figure.NewFigure("By Sonatype & Friends", "pepper", true).Print()
	fmt.Println("Ahab version: " + buildversion.BuildVersion)
}
