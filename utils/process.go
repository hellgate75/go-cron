package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Execute a set of command in a single string
func ExecuteCommandString(commandString string) (string, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	var tmp = strings.Split(commandString, " ")
	var command = make([]string, 0)
	for _, v := range tmp {
		if len(strings.TrimSpace(v)) > 0 {
			command = append(command, strings.TrimSpace(v))
		}
	}
	if len(command) == 0 {
		return "", errors.New("Command must have a least one not empty value")
	}
	cmdSubject := command[0]
	cmdArgs := command[1:]
	cmd := exec.Command(cmdSubject, cmdArgs...)
	if cmd == nil {
		return "", errors.New("Nil command cannot be executed")
	}
	var stdoutStderr = make([]byte, 0)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", stdoutStderr), err
}

// Execute a Command by tokens
func ExecuteCommandArgs(command ...string) (string, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	if len(command) == 0 {
		return "", errors.New("Command must have a least one value")
	}
	cmdSubject := command[0]
	if len(cmdSubject) == 0 {
		return "", errors.New("Command subject must not be empty")
	}
	cmdArgs := command[1:]
	cmd := exec.Command(cmdSubject, cmdArgs...)
	if cmd == nil {
		return "", errors.New("Nil command cannot be executed")
	}
	var stdoutStderr = make([]byte, 0)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", stdoutStderr), err
}

// Execute a Command with given arguments
func ExecuteCommand(command string, cmdArgs ...string) (string, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	if len(command) == 0 {
		return "", errors.New("Command subject must not be empty")
	}
	cmd := exec.Command(command, cmdArgs...)
	if cmd == nil {
		return "", errors.New("Nil command cannot be executed")
	}
	var stdoutStderr = make([]byte, 0)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", stdoutStderr), err
}
