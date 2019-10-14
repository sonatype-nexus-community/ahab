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
		if strings.Contains(pkg, "WARNING") {
		} else {
			projectList.Projects = append(projectList.Projects, doApkShowParse(pkg))
		}
	}
	return
}

func doApkShowParse(pkg string) (parsedProject types.Projects) {
	pkg = strings.TrimSpace(pkg)
	splitPackage := strings.Split(pkg, " ")
	parsedProject.Name = doApkShowParseName(splitPackage[0])
	parsedProject.Version = doApkShowParseVersion(splitPackage[0])
	return
}

func doApkShowParseName(pkg string) (name string) {
	re, err := regexp.Compile("([a-zA-Z][a-zA-Z0-9]+-)+")
	if err != nil {
		fmt.Println(err)
	}
	results := re.FindStringSubmatch(pkg)
	if results != nil {
		name = strings.TrimSuffix(results[0], "-")
	}
	return
}

func doApkShowParseVersion(pkg string) (version string) {
	re, err := regexp.Compile("([0-9]+)(\\.[0-9]+)(\\.[0-9]+)?(-[r][0-9]+)")
	if err != nil {
		fmt.Println(err)
	}
	results := re.FindStringSubmatch(pkg)
	if results != nil {
		version = results[0]
	}
	return
}
