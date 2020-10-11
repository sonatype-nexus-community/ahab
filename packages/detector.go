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

// Package packages implements package manager detection and package
// manager specific output formatting.
package packages

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

var execCommand = exec.Command

// SupportedPackageManagers represents the standard error string used
// when OS package maanger can not be identified.
const SupportedPackageManagers = "No supported package manager found; apk, dpkg, dnf or yum installed?"

// Release allows tests to inject ID and VERSION_ID typically
// parsed from os-release.
type Release struct {
	ID      []byte
	Version []byte
}

// DetectPackageManager parses os-release file to determine package
// manager based on distribution ID.
//
// Optional release short circuits os-release parsing to force
// appropriate package manager.
func DetectPackageManager(logger *logrus.Logger, r ...Release) (string, error) {
	var packageManager string
	var id []byte
	var version []byte

	if len(r) == 0 {
		data, err := readReleaseFile()
		if err != nil {
			return packageManager, err
		}
		raw := bytes.Split(data, []byte("\n"))

		id, err = parseField(raw, "ID")
		if err != nil {
			return packageManager, err
		}

		version, err = parseField(raw, "VERSION_ID")
		if err != nil {
			return packageManager, err
		}
	} else {
		id = r[0].ID
		version = r[0].Version
	}

	switch string(id) {
	case "alpine":
		packageManager = "apk"
	case "debian", "ubuntu", "mint":
		packageManager = "dpkg"
	case "fedora", "centos", "rhel", "redhat":
		if string(id) == "fedora" {
			if v, _ := strconv.Atoi(string(version)); v <= 22 {
				packageManager = "yum"
				break
			}
		}
		if v, _ := strconv.Atoi(string(version)); v <= 7 {
			packageManager = "yum"
			break
		}
		packageManager = "dnf"
	default:
		err := errors.New(SupportedPackageManagers)
		logger.Errorf("Error: %s\n", err.Error())
		return packageManager, err
	}

	logger.Infof("Detected package manager: %s\n", packageManager)
	return packageManager, nil
}

// Try to read os-release file.
// https://www.freedesktop.org/software/systemd/man/os-release.html
func readReleaseFile() ([]byte, error) {
	var data []byte

	files := []string{
		"/etc/os-release",
		"/usr/lib/os-release",
	}

	var file string
	for _, f := range files {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			file = f
			break
		}
	}

	if file == "" {
		return data, errors.New("Unable to read os-release")
	}

	f, err := os.Open(file)
	if err != nil {
		return data, err
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return data, err
	}

	data = make([]byte, s.Size())
	_, err = f.Read(data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Given os-release contents and field name, return value.
// FIELD="value" -> value
func parseField(raw [][]byte, field string) ([]byte, error) {
	var parsed []byte
	for _, v := range raw {
		if matched, _ := regexp.Match(fmt.Sprintf("^%s=.*$", field), v); matched {
			parsed = bytes.Split(v, []byte("="))[1]
			parsed = bytes.Trim(parsed, "\" ")
			return bytes.ToLower(parsed), nil
		}
	}
	return parsed, fmt.Errorf("Failed to parse os-release field: %s", field)
}
