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

func ParseApkShow(packages []string) (projectList types.ProjectList) {
	for _, pkg := range packages {
		if !strings.Contains(pkg, "WARNING") {
			projectList.Projects = append(projectList.Projects, doApkShowParse(pkg))
		}
	}
	return
}

func doApkShowParse(pkg string) (parsedProject types.Projects) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Split(pkg, " ")
	re, err := regexp.Compile(`^((.*)-([^a-zA-Z].*)-.*)`)
	if err != nil {
		panic(err)
	}
	newSlice := re.FindStringSubmatch(splitPackage[0])
	if newSlice != nil {
		parsedProject.Name = newSlice[2]
		parsedProject.Version = newSlice[3]
	} else {
		fmt.Printf("Failure parsing name, version for package")
	}
	return
}
