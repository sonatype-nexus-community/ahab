//
// Copyright (c) 2019-present Sonatype, Inc.
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

package parse

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func ParseYumListFromStdIn(stdin []string) (projectList ProjectList) {
	for _, pkg := range stdin {
		if strings.Contains(pkg, "Loaded plugins:") {
			log.Println("Skipping loaded plugins")
		} else if strings.TrimSpace(pkg) == "Installed Packages" {
			log.Println("Found beginning line of Yum Install List")
		} else if len(strings.Fields(pkg)) < 3 {
			log.Println("Skipping malformed package:", pkg)
		} else {
			projectList.Projects = append(projectList.Projects, doYumParseStdIn(pkg))
		}
	}
	return
}

func ParseYumList(packages []string) (projectList ProjectList) {
	for _, pkg := range packages {
		projectList.Projects = append(projectList.Projects, doYumParse(pkg))
	}
	return
}

func doYumParseStdIn(pkg string) (parsedProject Projects) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Fields(pkg)
	parsedProject.Name = doParseYumName(splitPackage[0])
	parsedProject.Version = doParseYumVersionIntoPurl(splitPackage[1])
	return
}

func doYumParse(pkg string) (parsedProject Projects) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Split(pkg, " ")
	parsedProject.Name = doParseYumName(splitPackage[0])
	parsedProject.Version = doParseYumVersionIntoPurl(splitPackage[1])
	return
}

func doParseYumName(name string) (newName string) {
	newSlice := strings.Split(name, ".")
	newName = newSlice[0]
	return
}

func doParseYumVersionIntoPurl(version string) (newVersion string) {
	re, err := regexp.Compile(`^([0-9]+[:-])?(([0-9]+)(\.([0-9]+))?(\.([0-9]+))?)`)
	if err != nil {
		fmt.Println(err)
	}
	newSlice := re.FindStringSubmatch(version)
	newVersion = newSlice[2]
	return
}
