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

func TestParseYumList(t *testing.T) {
	var list []string
	list = append(list, "bzip2-libs.x86_64 1.0.6-13.el7")
	list = append(list, "cpio.x86_64 2.11-27.el7")
	list = append(list, "elfutils-default-yama-scope.noarch 0.172-2.el7")
	result := ParseYumList(list)

	if len(result.Projects) != 3 {
		t.Errorf("Didn't work")
	}

	if result.Projects[0].Name != "bzip2-libs" || result.Projects[0].Version != "1.0.6" {
		t.Errorf("bzip2-libs dep did not match result. Actual %s", result.Projects[0])
	}
	if result.Projects[1].Name != "cpio" || result.Projects[1].Version != "2.11" {
		t.Errorf("cpio dep did not match result. Actual %s", result.Projects[1])
	}
	if result.Projects[2].Name != "elfutils-default-yama-scope" || result.Projects[2].Version != "0.172" {
		t.Errorf("elfutils-default-yama-scope dep did not match result. Actual %s", result.Projects[2])
	}
}

func TestParseYumListFromStdIn(t *testing.T) {
	var list []string
	list = append(list, "Loaded plugins: fastestmirror, ovl")
	list = append(list, "Installed Packages")
	list = append(list, "ncurses.x86_64 5.9-14.20130511.el7_4 @CentOS")
	list = append(list, "coreutils.x86_64 8.22-23.el7 @CentOS")
	list = append(list, "expat.x86_64 2.1.0-10.el7_3 @CentOS")
	result := ParseYumListFromStdIn(list)

	if len(result.Projects) != 3 {
		t.Errorf("Didn't work")
	}

	if result.Projects[0].Name != "ncurses" || result.Projects[0].Version != "5.9" {
		t.Errorf("ncurses dep did not match result. Actual %s", result.Projects[0])
	}
	if result.Projects[1].Name != "coreutils" || result.Projects[1].Version != "8.22" {
		t.Errorf("coreutils dep did not match result. Actual %s", result.Projects[1])
	}
	if result.Projects[2].Name != "expat" || result.Projects[2].Version != "2.1.0" {
		t.Errorf("expat dep did not match result. Actual %s", result.Projects[2])
	}
}
