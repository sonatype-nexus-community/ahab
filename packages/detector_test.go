// Copyright 2019 Nathan Zender
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
package packages

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sonatype-nexus-community/ahab/logger"
)

var testType string

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=" + testType}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	fmt.Println(">>> TEST PROCESS <<<<")
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "notinstalled" ||
		os.Getenv("GO_WANT_HELPER_PROCESS") == "yuminstalled" ||
		os.Getenv("GO_WANT_HELPER_PROCESS") == "dpkgqueryinstalled" ||
		os.Getenv("GO_WANT_HELPER_PROCESS") == "apkinstalled" ||
		os.Getenv("GO_WANT_HELPER_PROCESS") == "dnfinstalled" {
		expectedProgram := os.Args[4]
		if os.Getenv("GO_WANT_HELPER_PROCESS") == "notinstalled" {
			fmt.Println(">>> NOT INSTALLED CODE ... ALWAYS RETURN 1 <<<<")
			os.Exit(1)
		} else if os.Getenv("GO_WANT_HELPER_PROCESS") == "yuminstalled" {
			fmt.Fprintln(os.Stdout, ">>> YUM CODE <<<<")
			if expectedProgram == "yum" {
				fmt.Fprintln(os.Stdout, ">>> YUM CODE : ASKING FOR YUM ... ITS INSTALLED <<<<")
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stdout, ">>> YUM CODE : ASKING FOR SOMETHING ELSE ... ITS NOT INSTALLED <<<<")
				os.Exit(1)
			}
		} else if os.Getenv("GO_WANT_HELPER_PROCESS") == "dpkgqueryinstalled" {
			fmt.Println(">>> DKPG-QUERY CODE <<<<")
			if expectedProgram == "dpkg-query" {
				fmt.Fprintln(os.Stdout, ">>> DKPG-QUERY CODE : ASKING FOR APT ... ITS INSTALLED <<<<")
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stdout, ">>> DKPG-QUERY CODE : ASKING FOR SOMETHING ELSE ... ITS NOT INSTALLED <<<<")
				os.Exit(1)
			}
		} else if os.Getenv("GO_WANT_HELPER_PROCESS") == "apkinstalled" {
			fmt.Println(">>> APK CODE <<<<")
			if expectedProgram == "apk" {
				fmt.Fprintln(os.Stdout, ">>> APK CODE : ASKING FOR APK ... ITS INSTALLED <<<<")
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stdout, ">>> APK CODE : ASKING FOR SOMETHING ELSE ... ITS NOT INSTALLED <<<<")
				os.Exit(1)
			}
		} else if os.Getenv("GO_WANT_HELPER_PROCESS") == "dnfinstalled" {
			fmt.Println(">>> DNF CODE <<<<")
			if expectedProgram == "dnf" {
				fmt.Fprintln(os.Stdout, ">>> DNF CODE : ASKING FOR DNF ... ITS INSTALLED <<<<")
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stdout, ">>> DNF CODE : ASKING FOR SOMETHING ELSE ... ITS NOT INSTALLED <<<<")
				os.Exit(1)
			}
		}
	} else {
		fmt.Println(">>> RETURNING <<<<")
		return
	}
}

func TestDetectPackageManager(t *testing.T) {
	tests := map[string]struct {
		expectedInstalledPackageManager string
		expectedResult                  string
		expectedErr                     error
	}{
		"apk":            {expectedInstalledPackageManager: "apkinstalled", expectedResult: "apk", expectedErr: nil},
		"yum":            {expectedInstalledPackageManager: "yuminstalled", expectedResult: "yum", expectedErr: nil},
		"dnf":            {expectedInstalledPackageManager: "dnfinstalled", expectedResult: "dnf", expectedErr: nil},
		"dpkg-query":     {expectedInstalledPackageManager: "dpkgqueryinstalled", expectedResult: "dpkg", expectedErr: nil},
		"none installed": {expectedInstalledPackageManager: "notinstalled", expectedResult: "", expectedErr: errors.New(SupportedPackageManagers)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			testType = test.expectedInstalledPackageManager

			execCommand = fakeExecCommand
			defer func() { execCommand = exec.Command }()

			log, _ := logger.GetLogger(logrus.TraceLevel)
			actualResult, actualErr := DetectPackageManager(log)

			if actualResult != test.expectedResult {
				t.Errorf("Expected %q, got %q", test.expectedResult, actualResult)
			}
			if actualErr != nil && test.expectedErr != nil && actualErr.Error() != test.expectedErr.Error() {
				t.Errorf("Expected %q, got %q", test.expectedErr, actualErr)
			}
		})
	}
}
