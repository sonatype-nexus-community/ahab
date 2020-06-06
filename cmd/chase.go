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
	chaseCmd.PersistentFlags().StringVar(&operating, "os", "debian", "Specify a value for the operating system type you want to scan (alpine, debian, fedora)")
	chaseCmd.PersistentFlags().BoolVar(&cleanCache, "clean-cache", false, "Flag to clean the database cache for OSS Index")
	chaseCmd.PersistentFlags().StringVar(&ossIndexUser, "user", "", "Specify your OSS Index Username")
	chaseCmd.PersistentFlags().StringVar(&ossIndexToken, "token", "", "Specify your OSS Index API Token")
	chaseCmd.PersistentFlags().StringVar(&output, "output", "text", "Specify the output type you want (json, text, csv)")
	chaseCmd.PersistentFlags().BoolVar(&loud, "loud", false, "Specify if you want non vulnerable packages included in your output")
	chaseCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "Quiet removes the header from being printed")
	chaseCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Specify if you want no color in your results")
	chaseCmd.PersistentFlags().CountVarP(&verbose, "", "v", "Set log level, higher is more verbose")
}

var chaseCmd = &cobra.Command{
	Use:   "chase",
	Short: "chase is used for auditing projects with OSS Index",
	Example: `
	dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase --os debian
	yum list installed | ./ahab chase --os fedora
	apk info -vv | sort | ./ahab chase --os alpine
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
				cmd.Usage()
				logger.PrintErrorAndLogLocation(err)
			}
		}()

		if output == "text" && !quiet {
			printHeader()
		}

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

		if cleanCache {
			err = ossi.NoCacheNoProblems()
			if err != nil {
				logLady.Error(err)
				panic(err)
			}
			logLady.Trace("Cleaned Ahab Cache")
			return
		}

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
		return nil, fmt.Errorf("Nothing passed in to standard in")
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
