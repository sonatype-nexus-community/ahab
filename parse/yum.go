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
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	types "github.com/sonatype-nexus-community/nancy/types"
)

func ParseYumListFromStdIn(stdin []string) (projectList types.ProjectList) {
	for _, pkg := range stdin {

		if strings.Contains(pkg, "Loaded plugins") {
			log.Println("Found beginning line of Yum Install List")
		} else if strings.Contains(pkg, "Installed Packages") {
			log.Println("Found install header line of Yum Install List")
		} else {
			parsedProject, err := doYumParseStdIn(pkg)
			if err == nil {
				projectList.Projects = append(projectList.Projects, parsedProject)
			}
		}
	}
	return
}

func ParseYumList(packages []string) (projectList types.ProjectList) {
	for _, pkg := range packages {
		parsedProject, err := doYumParse(pkg)
		if err == nil {
			projectList.Projects = append(projectList.Projects, parsedProject)
		}
	}
	return
}

func doYumParseStdIn(pkg string) (parsedProject types.Projects, err error) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Fields(pkg)
	newVersion, err := doParseYumVersionIntoPurl(splitPackage[1])
	if err != nil {
		return parsedProject, err
	}
	parsedProject.Name = doParseYumName(splitPackage[0])
	parsedProject.Version = newVersion
	return
}

func doYumParse(pkg string) (parsedProject types.Projects, err error) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Split(pkg, " ")
	newVersion, err := doParseYumVersionIntoPurl(splitPackage[1])
	if err != nil {
		return parsedProject, err
	}
	parsedProject.Name = doParseYumName(splitPackage[0])
	parsedProject.Version = newVersion
	return
}

func doParseYumName(name string) (newName string) {
	newSlice := strings.Split(name, ".")
	newName = newSlice[0]
	return
}

func doParseYumVersionIntoPurl(version string) (newVersion string, err error) {
	re, err := regexp.Compile(`^([0-9]+[:-])?(([0-9]+)(\.([0-9]+))?(\.([0-9]+))?)`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(">>>>" + version)
	newSlice := re.FindStringSubmatch(version)
	fmt.Println(newSlice)
	if len(newSlice) >= 3 {
		newVersion = newSlice[2]
		return newVersion, nil
	}else{
		return newVersion, errors.New(fmt.Sprintf("Version of %s is not semVer", newSlice))
	}
}
