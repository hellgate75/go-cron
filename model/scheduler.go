package model

import (
	"fmt"
	"time"
)

// Generic Command Type
type CommandValue interface{}

func CommandValueToString(c CommandValue) string {
	return fmt.Sprintf("%v", c)
}

func TypeOfCommandValue(c CommandValue) string {
	return fmt.Sprintf("%T", c)
}

type Execution struct{
	Command				CommandValue								`yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
	Next				time.Time									`yaml:"nextExecution,omitempty" json:"nextExecution,omitempty" xml:"next-execution,omitempty"`
	Last				time.Time									`yaml:"lastExecution,omitempty" json:"lastExecution,omitempty" xml:"last-execution,omitempty"`
}

type Scheduler interface {
	IsRunning() bool
	Start() error
	Stop() error
	RunOnce() error
	Running() []CommandValue
	LoadFrom(path string)
	References() []CommandConfigRef
	Planned() []CommandConfig
	NextRunningTasks() []Execution
	AddAndSave(cmd CommandConfig) error
	UpdateAndSave(cmd CommandConfig, index int) error
	DeleteAndSave(index int) error
}