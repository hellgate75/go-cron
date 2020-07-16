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
	ContextMap *map[string]interface{}
	// Node global cache map, to store and recover common data
	StaticMap *map[string]interface{}
	// Infra-Node Cluster global cache map, to store and recover common data
	GlobalMap *map[string]interface{}
	// Allows developers to send warnings in the log
	WarningsPipe	chan error
	// Allows developers to send errors in the log
	ErrorsPipe		chan error
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
	UUID    string        		`yaml:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
	Command CommandConfig 		`yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
	Next    time.Time     		`yaml:"nextExecution,omitempty" json:"nextExecution,omitempty" xml:"next-execution,omitempty"`
	Last    time.Time			`yaml:"lastExecution,omitempty" json:"lastExecution,omitempty" xml:"last-execution,omitempty"`
	Times   int     			`yaml:"numberOfExecutions,omitempty" json:"numberOfExecutions,omitempty" xml:"number-of-executions,omitempty"`
	Scheduled bool	   			`yaml:"scheduled,omitempty" json:"scheduled,omitempty" xml:"scheduled,omitempty"`
	Map map[string]interface{}	`yaml:"scheduled,omitempty" json:"scheduled,omitempty" xml:"scheduled,omitempty"`
}

// Reset Command time table
func (e Execution) Reset() {
	e.Next = e.Last
}

// Reset Command time table
func (e Execution) Expired() bool {
	e.UpdateNext()
	if time.Since(e.Next).Nanoseconds() >= 0 || e.Scheduled {
		return false
	}
	return true
}

// Update Command time table and calculate Next Execution
func (e Execution) UpdateNext() {
	if ! e.Scheduled {
		c := e.Command
		if c.OnDemand {
			e.Last = time.Now()
			e.Next = time.Now().Add(20 * time.Second)
		} else {
			if c.Period != "" {
				d, err := time.ParseDuration(c.Period)
				if err == nil {
					if time.Since(c.Since).Nanoseconds() > 0 {
						if e.Last.Sub(c.Since).Nanoseconds() > 0 {
							e.Next = c.Since
						} else {
							e.Next = e.Last.Add(d)
						}
					} else {
						e.Next = e.Last.Add(d)
					}
				} else {
					e.Next = e.Last
				}
			} else if c.Repeat > 0 {
				if e.Times <= c.Repeat {
					if time.Since(c.Since).Nanoseconds() > 0 {
						if e.Last.Sub(c.Since).Nanoseconds() > 0 {
							e.Next = c.Since
						} else {
							e.Next = time.Now().Add(20 * time.Second)
						}
					} else {
						e.Next = time.Now().Add(20 * time.Second)
					}
				} else {
					e.Next = e.Last
				}
			}
		}
	}
	//e.Next =
}

func (e Execution) NeedScheduling() bool {
	e.UpdateNext()
	return ! e.Scheduled && time.Since(e.Next).Nanoseconds() >= 0
}

// Describes Scheduler behaviours and capabilities
type Scheduler interface {
	// Checks if scheduler is still/already running
	IsRunning() bool
	// Start scheduler process
	Start() error
	// Stop scheduler
	Stop() error
	// Run scheduler once and exit
	RunOnce() error
	// Collects all running tasks
	Running() []Execution
	// Load scheduler data from the device
	Load() error
	// Collects all planned tasks reference information
	References() []CommandConfigRef
	// Collects all planned tasks
	Planned() []CommandConfig
	// Collects all next running tasks
	NextRunningTasks() []Execution
	// Add a task and persist data
	AddAndPersist(cmd CommandConfig) error
	// Update a task and persist data
	UpdateAndPersist(cmd CommandConfig, index int) error
	// Delete a task and persist data
	DeleteAndPersist(index int) error
	// Add a task to cache without persist (function executable tasks)
	AddToCache(cmd CommandConfig) error
	// Update a task to cache without persist (function executable tasks)
	UpdateToCache(cmd CommandConfig, index int) error
	// Delete a task to cache without persist (function executable tasks)
	DeleteFromCache(index int) error
	// Waits until scheduler finish
	Wait()
	// Retrieves the scheduler errors channel, used to report live errors from scheduler or scheduler tasks
	Errors() chan error
	// Retrieves the scheduler warnings channel, used to report live warnings from scheduler or scheduler tasks
	Warnings() chan error
	// Stop Scheduler if required, then Destroy content, data and shared memory items
	// and eventually save last state fto the device if required
	Destroy(saveState bool)
}
