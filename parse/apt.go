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
package parse

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	types "github.com/sonatype-nexus-community/nancy/types"
)

func ParseAptListFromStdIn(stdin []string) (projectList types.ProjectList) {
	for _, pkg := range stdin {

		if strings.TrimSpace(pkg) == "Listing... Done" {
			log.Println("Found beginning line of Apt Install List")
		} else {
			projectList.Projects = append(projectList.Projects, doAptParseStdIn(pkg))
		}
	}
	return
}

func ParseAptList(packages []string) (projectList types.ProjectList) {
	for _, pkg := range packages {
		projectList.Projects = append(projectList.Projects, doAptParse(pkg))
	}
	return
}

func doAptParseStdIn(pkg string) (parsedProject types.Projects) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Split(pkg, " ")
	parsedProject.Name = strings.Split(splitPackage[0], "/")[0]
	parsedProject.Version = doParseAptVersionIntoPurl(splitPackage[1])
	return
}

func doAptParse(pkg string) (parsedProject types.Projects) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Split(pkg, " ")
	parsedProject.Name = splitPackage[0]
	parsedProject.Version = doParseAptVersionIntoPurl(splitPackage[1])
	return
}

func doParseAptVersionIntoPurl(version string) (newVersion string) {
	re, err := regexp.Compile(`^([0-9]+:)?(([0-9]+)\.([0-9]+)(\.([0-9]+))?)`)
	if err != nil {
		fmt.Println(err)
	}
	newSlice := re.FindStringSubmatch(version)
	fmt.Println(newSlice)
	newVersion = newSlice[2]
	return
}
