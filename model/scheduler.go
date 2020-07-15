package model

import (
	"fmt"
	"time"
)

// Generic Command Type
type CommandValue interface{}

// Execution context, that will be passed to the ComputableValue
type ExecutionContext struct {
	// Node configuration
	Configuration *SchedulerConfig
	// Information about the process
	CommandInfo *CommandConfig
	// Cache map used to repeat the execution of the same function, when repeated
	contextMap *map[string]interface{}
	// Node global cache map, to store and recover common data
	staticMap *map[string]interface{}
	// Infra-Node Cluster global cache map, to store and recover common data
	globalMap *map[string]interface{}
}

// Describes interface that can be executed in the Scheduler (passed as CommandValue) with self encapsulation of the running process
type ComputableValue interface {
	Compute(ExecutionContext) error
}

func CommandValueToString(c CommandValue) string {
	return fmt.Sprintf("%v", c)
}

func TypeOfCommandValue(c CommandValue) string {
	return fmt.Sprintf("%T", c)
}

type Execution struct {
	UUID    string        `yaml:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
	Command CommandConfig `yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
	Next    time.Time     `yaml:"nextExecution,omitempty" json:"nextExecution,omitempty" xml:"next-execution,omitempty"`
	Last    time.Time     `yaml:"lastExecution,omitempty" json:"lastExecution,omitempty" xml:"last-execution,omitempty"`
}

func (e Execution) Reset() {

}

func (e Execution) UpdateNext() {

}

type Scheduler interface {
	IsRunning() bool
	Start() error
	Stop() error
	RunOnce() error
	Running() []Execution
	Load() error
	References() []CommandConfigRef
	Planned() []CommandConfig
	NextRunningTasks() []Execution
	AddAndPersist(cmd CommandConfig) error
	UpdateAndPersist(cmd CommandConfig, index int) error
	DeleteAndPersist(index int) error
	AddToCache(cmd CommandConfig) error
	UpdateToCache(cmd CommandConfig, index int) error
	DeleteFromCache(index int) error
	Wait()
}
