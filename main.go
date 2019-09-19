// Copyright 2019 Sonatype Inc.
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
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	packages "github.com/sonatype-nexus-community/ahab/packages"
	parse "github.com/sonatype-nexus-community/ahab/parse"
	"github.com/sonatype-nexus-community/nancy/audit"
	"github.com/sonatype-nexus-community/nancy/ossindex"
)

func main() {
	chaseCommand := flag.NewFlagSet("chase", flag.ExitOnError)
	whalesPtr := chaseCommand.String("whales", "", "A comma separated list of packages to parse")
	packagePtr := chaseCommand.String("os", "debian", "Your target operating system")

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("chase subcommand is required")
		os.Exit(1)
	}

	switch args[0] {
	case "chase":
		parseChaseCommandArgs(chaseCommand, whalesPtr, packagePtr)
	}
}

func parseChaseCommandArgs(command *flag.FlagSet, flag *string, operating *string) {
	command.Parse(os.Args[2:])

	if command.Parsed() {
		if *flag == "" {
			tryParseStdIn(operating)
		} else {
			tryParseFlag(flag, operating)
		}
	}
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
	if thing == "debian" {
		var aptResult packages.Apt
		//aptResult.ProjectList = parse.ParseAptListFromStdIn(list)
		aptResult.ProjectList = parse.ParseDpkgList(list)
		tryExtractAndAudit(aptResult, thing)
	} else {
		var yumResult packages.Yum
		yumResult.ProjectList = parse.ParseYumListFromStdIn(list)
		tryExtractAndAudit(yumResult, thing)
	}
}

func tryAuditPackages(purls []string, count int) {
	fmt.Print(purls)
	coordinates, err := ossindex.AuditPackages(purls)
	if err != nil {
		fmt.Print(err)
	}
	if count := audit.LogResults(true, count, coordinates); count > 0 {
		os.Exit(1)
	}
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
