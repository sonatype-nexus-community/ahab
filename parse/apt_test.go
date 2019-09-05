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
	"testing"
)

func TestParseAptList(t *testing.T) {
	var list []string
	list = append(list, "libedit2 3.1-20170329-1")
	list = append(list, "libmount1 2.31.1-0.4ubuntu3.3")
	list = append(list, "zlib1g 1:1.2.11.dfsg-0ubuntu2")
	result := ParseAptList(list)

	if len(result.Projects) != 3 {
		t.Errorf("Didn't work")
	}
}
