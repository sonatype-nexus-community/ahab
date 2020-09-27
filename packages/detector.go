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

const SupportedPackageManagers = "Supported package managers are apk, dnf, dpkg or yum; could not find any. Possible issues: 1.) dpkg, apk, yum or dnf is not installed. 2.) 'which' program is not installed to do auto detection"

func DetectPackageManager(logger *logrus.Logger) (string, error) {
	var packageManager string

	installed := determineIfPackageManagerInstalled("apk", logger)
	if installed {
		packageManager = "apk"
		return packageManager, nil
	}
	installed = determineIfPackageManagerInstalled("dnf", logger)
	if installed {
		packageManager = "dnf"
		return packageManager, nil
	}
	installed = determineIfPackageManagerInstalled("yum", logger)
	if installed {
		packageManager = "yum"
		return packageManager, nil
	}
	installed = determineIfPackageManagerInstalled("dpkg-query", logger)
	if installed {
		packageManager = "dpkg"
		return packageManager, nil
	} else {
		return packageManager, errors.New(SupportedPackageManagers)
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
			return waitStatus.ExitStatus() == 0
		}
		return false
	}
	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	logger.Info(string(output))
	logger.Infof("Output 2: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	return waitStatus.ExitStatus() == 0
}
