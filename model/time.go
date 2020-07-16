package model

import (
	"time"
)


// Describes the state of the scheduler
type SchedulerState string

const (
	SchedulerStateRunning	= SchedulerState("running")
	SchedulerStateStopped	= SchedulerState("stopped")
	SchedulerStatePaused	= SchedulerState("paused")
)

// Defines the scheduler configuration
type CommandConfig struct {
	OnDemand			bool										`yaml:"onDemand,omitempty" json:"onDemand,omitempty" xml:"onDemand,omitempty"`
	Period				string										`yaml:"period,omitempty" json:"period,omitempty" xml:"period,omitempty"`
	Repeat				int											`yaml:"repeat,omitempty" json:"repeat,omitempty" xml:"repeat,omitempty"`
	Since				time.Time									`yaml:"since,omitempty" json:"since,omitempty" xml:"since,omitempty"`
	Command				CommandValue								`yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
}

// Defines reference the scheduler configuration
type CommandConfigRef struct {
	UUID				string										`yaml:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
	Command				CommandValue								`yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
	Created				time.Time									`yaml:"created,omitempty" json:"created,omitempty" xml:"created,omitempty"`
	Updated				time.Time									`yaml:"updated,omitempty" json:"updated,omitempty" xml:"updated,omitempty"`
	FirstRun			time.Time									`yaml:"fistExecution,omitempty" json:"fistExecution,omitempty" xml:"first-execution,omitempty"`
	LastRun				time.Time									`yaml:"lastExecution,omitempty" json:"lastExecution,omitempty" xml:"last-execution,omitempty"`
}


// Defines the scheduler configuration if the scheduler is configured in sync mode it will run all tasks immediately all together.
type SchedulerConfig struct {
	Sync				bool										`yaml:"sync,omitempty" json:"sync,omitempty" xml:"sync,omitempty"`
	Commands			[]CommandConfigRef							`yaml:"commands,omitempty" json:"commands,omitempty" xml:"command,omitempty"`
}
