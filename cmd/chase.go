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
	"regexp"
	"strings"
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

type CveListFlag struct {
	Cves []string
}

func (cve *CveListFlag) String() string {
	return fmt.Sprint(cve.Cves)
}

func (cve *CveListFlag) Set(value string) error {
	if len(cve.Cves) > 0 {
		return fmt.Errorf("The CVE Exclude Flag is already set")
	}
	cve.Cves = strings.Split(strings.ReplaceAll(value, " ", ""), ",")

	return nil
}

func (cve *CveListFlag) Type() string { return "CveListFlag" }

var (
	packageManager               string
	cleanCache                   bool
	ossIndexUser                 string
	ossIndexToken                string
	output                       string
	loud                         bool
	quiet                        bool
	noColor                      bool
	excludeVulnerabilityFilePath string
	cveList                      CveListFlag
	unixComments                 = regexp.MustCompile(`#.*$`)
	untilComment                 = regexp.MustCompile(`(until=)(.*)`)
	ossi                         *ossindex.Server
)

func init() {
	rootCmd.AddCommand(chaseCmd)

	pf := chaseCmd.PersistentFlags()
	pf.StringVar(&packageManager, "os", "", "Specify a value for the operating system type you want to scan (alpine, debian, fedora). Useful if autodetection fails and/or you want to explicitly set it.")
	pf.StringVar(&packageManager, "package-manager", "", "Specify package manager type you want to scan (apk, dnf, dpkg or yum). Useful if autodetection fails and/or you want to explicitly set it.")
	pf.BoolVar(&cleanCache, "clean-cache", false, "Flag to clean the database cache for OSS Index")
	pf.StringVar(&ossIndexUser, "user", "", "Specify your OSS Index Username")
	pf.StringVar(&ossIndexToken, "token", "", "Specify your OSS Index API Token")
	pf.StringVar(&output, "output", "text", "Specify the output type you want (json, text, csv)")
	pf.BoolVar(&loud, "loud", false, "Specify if you want non vulnerable packages included in your output")
	pf.BoolVar(&quiet, "quiet", false, "Quiet removes the header from being printed")
	pf.BoolVar(&noColor, "no-color", false, "Specify if you want no color in your results")
	pf.CountVarP(&verbose, "", "v", "Set log level, higher is more verbose")

	chaseCmd.Flags().VarP(&cveList, "exclude-vulnerability", "e", "Comma separated list of CVEs to exclude")
	chaseCmd.Flags().StringVarP(&excludeVulnerabilityFilePath, "exclude-vulnerability-file", "x", "./.ahab-ignore", "Path to a file containing newline separated CVEs to be excluded")

	chaseCmd.Flag("os").Deprecated = "use package-manager"
}

var chaseCmd = &cobra.Command{
	Use:   "chase",
	Short: "chase is used for auditing projects with OSS Index",
	Example: `
	dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase
	yum list installed | ./ahab chase
	dnf list installed | ./ahab chase
	apk info -vv | sort | ./ahab chase
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

		err = getCVEExcludesFromFile(excludeVulnerabilityFilePath)

		if packageManager == "" {
			logLady.Trace("Attempting to detect package manager for you")
			manager, err := packages.DetectPackageManager(logLady)
			if err != nil {
				logLady.Error(err)
				panic(err)
			}
			packageManager = manager
		}

		logLady.Trace("Attempting to audit list of strings from standard in")
		pkgs, err := parseStdIn(&packageManager)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

		logLady.WithField("package-manager", packageManager).Trace("Attempting to extract purls from Project List")
		purls := pkgs.ExtractPurlsFromProjectList()

		logLady.Trace("Attempting to Audit Packages with OSS Index")
		coordinates, err := ossi.AuditPackages(purls)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

		logLady.Trace("Attempting to output audited packages results")
		count, results, err := audit.LogResults(noColor, loud, output, coordinates, cveList.Cves)
		if err != nil {
			logLady.Error(err)
			panic(err)
		}

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

func parseStdInList(list []string, packageManager *string) (packages.IPackage, error) {
	thing := *packageManager
	logLady.WithFields(logrus.Fields{
		"list": list,
	}).Trace("Chasing ", thing)
	switch thing {
	case "dpkg", "debian":
		var aptResult packages.Apt
		aptResult.ProjectList = parse.ParseDpkgList(list)

		logLady.WithFields(logrus.Fields{
			"project_list": aptResult.ProjectList,
		}).Trace("Obtained dpkg project list")
		return aptResult, nil
	case "apk", "alpine":
		var apkResult packages.Apk
		apkResult.ProjectList = parse.ParseApkShow(list)

		logLady.WithFields(logrus.Fields{
			"project_list": apkResult.ProjectList,
		}).Trace("Obtained apk project list")
		return apkResult, nil
	case "yum", "dnf", "fedora":
		var dnfResult packages.Yum
		dnfResult.ProjectList = parse.ParseYumListFromStdIn(list)

		logLady.WithFields(logrus.Fields{
			"project_list": dnfResult.ProjectList,
		}).Trace("Obtained dnf project list")
		return dnfResult, nil
	default:
		var yumResult packages.Yum
		yumResult.ProjectList = parse.ParseYumListFromStdIn(list)

		logLady.WithFields(logrus.Fields{
			"project_list": yumResult.ProjectList,
		}).Trace("Obtained yum project list (default case)")
		return yumResult, nil
	}
}

func parseStdIn(packageManager *string) (packages.IPackage, error) {
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
		return nil, err
	}

	return parseStdInList(list, packageManager)
}

func printHeader() {
	figure.NewFigure("Ahab", "larry3d", true).Print()
	figure.NewFigure("By Sonatype & Friends", "pepper", true).Print()
	fmt.Println("Ahab version: " + buildversion.BuildVersion)
}

func getCVEExcludesFromFile(excludeVulnerabilityFilePath string) error {
	fi, err := os.Stat(excludeVulnerabilityFilePath)
	if (fi != nil && fi.IsDir()) || (err != nil && os.IsNotExist(err)) {
		return nil
	}
	file, err := os.Open(excludeVulnerabilityFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ogLine := scanner.Text()
		err := determineIfLineIsExclusion(ogLine)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func determineIfLineIsExclusion(ogLine string) error {
	line := unixComments.ReplaceAllString(ogLine, "")
	until := untilComment.FindStringSubmatch(line)
	line = untilComment.ReplaceAllString(line, "")
	cveOnly := strings.TrimSpace(line)

	if len(cveOnly) > 0 {
		if until != nil {
			parseDate, err := time.Parse("2006-01-02", strings.TrimSpace(until[2]))
			if err != nil {
				return fmt.Errorf("failed to parse until at line %q. Expected format is 'until=yyyy-MM-dd'", ogLine)
			}
			if parseDate.After(time.Now()) {
				cveList.Cves = append(cveList.Cves, cveOnly)
			}
		} else {
			cveList.Cves = append(cveList.Cves, cveOnly)
		}
	}

	return nil
}
