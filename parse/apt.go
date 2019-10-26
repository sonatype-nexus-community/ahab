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
	"regexp"
	"strings"

	"github.com/sonatype-nexus-community/nancy/types"
)

func ParseAptList(packages []string) (projectList types.ProjectList) {
	for _, pkg := range packages {
		parsedProject, err := doAptParse(pkg)
		if err == nil {
			projectList.Projects = append(projectList.Projects, parsedProject)
		}
	}
	return
}

func doAptParse(pkg string) (parsedProject types.Projects, err error) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Split(pkg, " ")
	newVersion := doParseAptVersionIntoPurl(splitPackage[0], splitPackage[1])
	parsedProject.Name = splitPackage[0]
	parsedProject.Version = newVersion
	return
}

func doParseAptVersionIntoPurl(name string, version string) (newVersion string) {
	// exclude prefix delimited by :, and drop suffixes after .
	re, err := regexp.Compile(`^([0-9]+:)?(([0-9]+)\.([0-9]+)(\.([0-9]+))?)`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(">>>>" + version)
	newSlice := re.FindStringSubmatch(version)
	if newSlice != nil {
		newVersion = newSlice[2]
	} else {
		// first approach failed, second attempt:
		// use prefix up to the first alphabetic character
		reNumericPrefix, err := regexp.Compile(`([^a-zA-Z]+)?`)
		if err != nil {
			fmt.Println(err)
		}
		numberPrefix := reNumericPrefix.FindStringSubmatch(version)
		if numberPrefix != nil {
			newVersion = numberPrefix[1]
		} else {
			// yikes, nothing we recognize. fallback to taking the string as is.
			fmt.Printf("package name: %s, using fallback value for version: %s\n", name, version)
			newVersion = version
		}
	}
	return
}

func ParseDpkgList(packages []string) (projectList types.ProjectList) {
	for _, pkg := range packages {
		projectList.Projects = append(projectList.Projects, doDpkgParse(pkg))
	}
	return
}

func doDpkgParse(pkg string) (parsedProject types.Projects) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Split(pkg, " ")
	parsedProject.Name = splitPackage[0]
	parsedProject.Version = doParseAptVersionIntoPurl(splitPackage[0], splitPackage[1])
	return
}
