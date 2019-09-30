package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

func TestSuccessfullyBuildAptDockerWithAhab(t *testing.T) {
	_, goStatus := runCommand("go", "build", "-o", "ahab", "../main.go")
	defer os.Remove("ahab")
	if goStatus == false {
		t.Error("Could not build ahab")
		return
	}

	output, status := runCommand("docker", "build", "-f", "apt/Dockerfile", ".")
	if status == false {
		if !strings.Contains(output, "Audited dependencies:"){
			t.Error("Docker build for apt failed and was not due to vulnerable packages. See test output for more details.")
			return
		}
	}
}

func TestSuccessfullyBuildYumDockerWithAhab(t *testing.T) {
	_, goStatus := runCommand("go", "build", "-o", "ahab", "../main.go")
	defer os.Remove("ahab")
	if goStatus == false {
		t.Error("Could not build ahab")
		return
	}

	output, status := runCommand("docker", "build", "-f", "yum/Dockerfile", ".")
	if status == false {
		if !strings.Contains(output, "Audited dependencies:"){
			t.Error("Docker build for yum failed and was not due to vulnerable packages. See test output for more details.")
			return
		}
	}
}

func runCommand(command string, args ...string) (output string, status bool){
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=amd64")
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")

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
			}else{
				return combinedOutputStr,false
			}
		}else{
			return combinedOutputStr, false
		}
	} else {
		// Success
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("Output 2: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		if waitStatus == 0 {
			return combinedOutputStr, true
		}else{
			return combinedOutputStr, false
		}
	}
}
