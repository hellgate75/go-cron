package cron

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-cron/io"
)

var Commands = []string{"help", "daemon", "add", "remove", "update", "list", "active", "next"}

func Exec(command string, path string, encoding io.Encoding) error {
	var err error
	switch command {
	case "daemon":
		return executeDaemonCommand()
	case "add":
		return executeAddCommand()
	case "remove":
		return executeRemoveCommand()
	case "update":
		return executeUpdateCommand()
	case "list":
		return executeListCommand()
	case "active":
		return executeActiveCommand()
	case "next":
		return executeNextCommand()
	default:
		fmt.Printf("Cannot describe unknown command: <%s>\n", command)
		fmt.Printf("Available commands: %v\n", Commands)
	}
	return err
}

func executeDaemonCommand() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	//TODO: Implement method
	err = parse(getDaemonCommandArgsParser())
	if err != nil {
		return err
	}
	return err
}

func executeAddCommand() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	//TODO: Implement method
	err = parse(getAddCommandArgsParser())
	if err != nil {
		return err
	}
	return err
}

func executeRemoveCommand() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	//TODO: Implement method
	err = parse(getRemoveCommandArgsParser())
	if err != nil {
		return err
	}
	return err
}


func executeUpdateCommand() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	//TODO: Implement method
	err = parse(getUpdateCommandArgsParser())
	if err != nil {
		return err
	}
	return err
}


func executeListCommand() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	//TODO: Implement method
	err = parse(getListCommandArgsParser())
	if err != nil {
		return err
	}
	return err
}

func executeActiveCommand() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	//TODO: Implement method
	err = parse(getActiveCommandArgsParser())
	if err != nil {
		return err
	}
	return err
}

func executeNextCommand() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	//TODO: Implement method
	err = parse(getNextCommandArgsParser())
	if err != nil {
		return err
	}
	return err
}
