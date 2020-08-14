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
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

var execCommand = exec.Command

const SupportedPackageManagers = "supported package managers are dpkg-query, apk or yum, could not find any. Possible issues: 1.) dpkg-query, apk or yum is not installed. 2.) 'which' program is not installed to do auto detection"

func DetectPackageManager(logger *logrus.Logger) (string, error) {
	var os string

	installed := determineIfPackageManagerInstalled("apk", logger)
	if installed {
		//Having this be OS is a little weird. It probably should have been just package manager based flag.
		os = "alpine"
		return os, nil
	}
	installed = determineIfPackageManagerInstalled("yum", logger)
	if installed {
		//Having this be OS is a little weird. It probably should have been just package manager based flag.
		os = "fedora"
		return os, nil
	}
	installed = determineIfPackageManagerInstalled("dpkg-query", logger)
	if installed {
		os = "debian"
		return os, nil
	} else {
		return os, errors.New(SupportedPackageManagers)
	}
}

func determineIfPackageManagerInstalled(packageManager string, logger *logrus.Logger) bool {
	cmd := execCommand("which", packageManager)
	var waitStatus syscall.WaitStatus
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Errorf("Error: %s\n", err.Error())
		logger.Infof(string(output))
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			logger.Infof("Output 1: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
			if waitStatus == 0 {
				return true
			}else{
				return false
			}
		}else{
			return false
		}
	} else {
		// Success
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		logger.Info(string(output))
		logger.Infof("Output 2: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		if waitStatus == 0 {
			return true
		}else{
			return false
		}
	}
}
