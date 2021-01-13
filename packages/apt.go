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

package packages

import (
	"fmt"

	"github.com/DarthHater/packageurl-go"
	"github.com/sonatype-nexus-community/ahab/parse"
)

type Apt struct {
	ProjectList parse.ProjectList
}

func (a Apt) ExtractPurlsFromProjectList() (purls []string) {
	for _, s := range a.ProjectList.Projects {
		var purl = fmt.Sprintf("pkg:deb/debian/%s@%s", s.Name, s.Version)
		purls = append(purls, purl)
	}
	return
}

func (a Apt) ExtractPurlObjectsFromProjectList() (purls []packageurl.PackageURL) {
	for _, s := range a.ProjectList.Projects {
		purl, err := packageurl.FromString(fmt.Sprintf("pkg:deb/debian/%s@%s", s.Name, s.Version))
		if err != nil {
			continue
		}
		purls = append(purls, purl)
	}
	return
}
