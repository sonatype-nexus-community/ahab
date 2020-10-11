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
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sonatype-nexus-community/ahab/logger"
)

func TestDetectPackageManager(t *testing.T) {
	tests := map[string]struct {
		release        Release
		expectedResult string
		expectedErr    error
	}{
		"alpine": {
			release: Release{
				ID:      []byte("alpine"),
				Version: []byte(""),
			},
			expectedResult: "apk",
			expectedErr:    nil,
		},
		"centos_new": {
			release: Release{
				[]byte("centos"),
				[]byte("8"),
			},
			expectedResult: "dnf",
			expectedErr:    nil,
		},
		"centos_old": {
			release: Release{
				[]byte("centos"),
				[]byte("7"),
			},
			expectedResult: "yum",
			expectedErr:    nil,
		},
		"fedora_new": {
			release: Release{
				[]byte("fedora"),
				[]byte("32"),
			},
			expectedResult: "dnf",
			expectedErr:    nil,
		},
		"fedora_old": {
			release: Release{
				[]byte("fedora"),
				[]byte("22"),
			},
			expectedResult: "yum",
			expectedErr:    nil,
		},
		"debian": {
			release: Release{
				[]byte("debian"),
				[]byte(""),
			},
			expectedResult: "dpkg",
			expectedErr:    nil,
		},
		"none installed": {
			release: Release{
				[]byte("unsupported"),
				[]byte(""),
			},
			expectedResult: "",
			expectedErr:    errors.New(SupportedPackageManagers),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			log, _ := logger.GetLogger(logrus.TraceLevel)
			actualResult, actualErr := DetectPackageManager(log, test.release)

			if actualResult != test.expectedResult {
				t.Errorf("Expected %q, got %q", test.expectedResult, actualResult)
			}
			if actualErr != nil && test.expectedErr != nil && actualErr.Error() != test.expectedErr.Error() {
				t.Errorf("Expected %q, got %q", test.expectedErr, actualErr)
			}
		})
	}
}
