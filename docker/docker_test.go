package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

func TestMain(m *testing.M) {
	if setupSubtest() {
		os.Exit(222)
	}
	exitVal := m.Run()

	err := teardownSubtest()
	if err != nil {
		fmt.Printf("welp teardown went terribly...not sure if it matters or not but here is your error: %v", err)
		os.Exit(333)
	}

	os.Exit(exitVal)
}

func TestDockerIntegration(t *testing.T) {
	tests := map[string]struct {
		expectedDockerfile string
	}{
		"dpkg-query using autodetect": {expectedDockerfile: "dpkg-query-autodetect/Dockerfile"},
		"dpkg-query":                  {expectedDockerfile: "dpkg-query/Dockerfile"},
		"yum":                         {expectedDockerfile: "yum/Dockerfile"},
		"yum using autodetect":        {expectedDockerfile: "yum-autodetect/Dockerfile"},
		"dnf":                         {expectedDockerfile: "dnf/Dockerfile"},
		"dnf using autodetect":        {expectedDockerfile: "dnf-autodetect/Dockerfile"},
		"apk using autodetect":        {expectedDockerfile: "apk-autodetect/Dockerfile"},
		"apk":                         {expectedDockerfile: "apk/Dockerfile"},
	}
	t.Run("docker", func(t *testing.T) {
		for name, test := range tests {
			name := name
			test := test
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				t.Logf(">>>>>>>>>>>>>> %v: Started", name)
				t.Logf(">>>>>>>>>>>>>> %v: %v", name, test.expectedDockerfile)
				output, status := runCommand("docker", "build", "--no-cache", "-f", test.expectedDockerfile, ".")
				t.Logf(">>>>>>>>>>>>>> %v: %v", name, output)
				if status == false {
					if !strings.Contains(output, "Audited Dependencies") {
						t.Errorf(">>>>>>>>>>>>>> %v: Docker build failed and was not due to vulnerable packages. See test output for more details.", name)
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
	output, goStatus := runCommand("go", "build", "-o", "ahab", "../main.go")
	if goStatus == false {
		fmt.Println("Could not build ahab")
		fmt.Println(output)
		return true
	}
	return false
}

func runCommand(command string, args ...string) (output string, status bool) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=amd64")
	// 'DOCKER_BUILDKIT=0' fixes errors like this: `#10 2.350 [output clipped, log limit 100KiB/s reached]`
	cmd.Env = append(cmd.Env, "DOCKER_BUILDKIT=0")

	var waitStatus syscall.WaitStatus
	combinedOutput, err := cmd.CombinedOutput()
	combinedOutputStr := string(combinedOutput)
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
