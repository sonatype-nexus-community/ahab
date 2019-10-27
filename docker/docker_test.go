package docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"testing"
)

func TestMain(m *testing.M) {
	if setupSubtest() {
		os.Exit(222)
	}
	exitVal := m.Run()

	teardownSubtest()

	os.Exit(exitVal)
}

func TestDockerIntegration(t *testing.T) {
	tests := map[string]struct {
		expectedDockerfile string
	}{
		"dpkg-query":           {expectedDockerfile: "dpkg-query/Dockerfile"},
		"yum":                  {expectedDockerfile: "yum/Dockerfile"},
		"yum using autodetect": {expectedDockerfile: "yum-autodetect/Dockerfile"},
		"apk using autodetect": {expectedDockerfile: "apk-autodetect/Dockerfile"},
	}
	t.Run("docker", func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				output, status := runCommand("docker", "build", "-f", test.expectedDockerfile, ".")
				if status == false {
					if !strings.Contains(output, "Audited dependencies:") {
						t.Error("Docker build failed and was not due to vulnerable packages. See test output for more details.")
						return
					}
				}
			})
		}
	})
}

func teardownSubtest() error {
	fmt.Println("[TEARDOWN]")
	return os.Remove("ahab")
}

func setupSubtest() bool {
	fmt.Println("[SETUP]")
	_, goStatus := runCommand("go", "build", "-o", "ahab", "../main.go")
	if goStatus == false {
		fmt.Println("Could not build ahab")
		return true
	}
	return false
}

func runCommand(command string, args ...string) (output string, status bool) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=amd64")
	if runtime.GOOS == "darwin" {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	}

	var waitStatus syscall.WaitStatus
	combinedOutput, err := cmd.CombinedOutput()
	combinedOutputStr := string(combinedOutput)
	fmt.Println(combinedOutputStr)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			fmt.Printf("Output 1: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
			if waitStatus == 0 {
				return combinedOutputStr, true
			} else {
				return combinedOutputStr, false
			}
		} else {
			return combinedOutputStr, false
		}
	} else {
		// Success
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("Output 2: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		if waitStatus == 0 {
			return combinedOutputStr, true
		} else {
			return combinedOutputStr, false
		}
	}
}
