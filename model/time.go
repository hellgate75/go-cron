package model

import (
	"strings"
	"time"
)

// Time unit type for scheduler
type TimeUnit string

// Describes unit string value
func (t TimeUnit) String() string {
	return strings.ToLower(string(t))
}

// Compares a given time unit with current one
func (t TimeUnit) Equals(tu TimeUnit) bool {
	return t.String() == tu.String()
}

// Compares a given time unit with current one
func (t TimeUnit) Bigger(tu TimeUnit) bool {
	return t.Duration() > tu.Duration()
}

// Compares a given time unit with current one
func (t TimeUnit) Smaller(tu TimeUnit) bool {
	return t.Duration() < tu.Duration()
}

func (t TimeUnit) Duration() time.Duration {
	switch t.String() {
	case "nanoseconds":
		return time.Nanosecond
	case "microseconds":
		return time.Microsecond
	case "millis":
		return time.Millisecond
	case "seconds":
		return time.Second
	case "minutes":
		return time.Minute
	case "hours":
		return time.Hour
	case "days":
		return 24 * time.Hour
	case "weeks":
		return 168 * time.Hour
	case "months":
		return 720 * time.Hour
	case "years":
		return 8760 * time.Hour
	default:
		return time.Nanosecond
	}
}

// Describes the state of the scheduler
type SchedulerState string

const (
	TimeUnitNanoseconds		= TimeUnit("ns")
	TimeUnitMicroseconds	= TimeUnit("mis")
	TimeUnitMilliseconds	= TimeUnit("ms")
	TimeUnitSeconds			= TimeUnit("sec")
	TimeUnitMinutes			= TimeUnit("min")
	TimeUnitHours			= TimeUnit("hr")
	TimeUnitDays			= TimeUnit("dd")
	TimeUnitWeeks			= TimeUnit("wk")
	TimeUnitMonths			= TimeUnit("mt")
	TimeUnitYears			= TimeUnit("yr")
	SchedulerStateRunning	= SchedulerState("running")
	SchedulerStateStopped	= SchedulerState("stopped")
	SchedulerStatePaused	= SchedulerState("paused")
)

// Defines the scheduler configuration
type CommandConfig struct {
	OnDemand			bool										`yaml:"onDemand,omitempty" json:"onDemand,omitempty" xml:"onDemand,omitempty"`
	Period				string										`yaml:"period,omitempty" json:"period,omitempty" xml:"period,omitempty"`
	Repeat				int											`yaml:"repeat,omitempty" json:"repeat,omitempty" xml:"repeat,omitempty"`
	Since				int											`yaml:"since,omitempty" json:"since,omitempty" xml:"since,omitempty"`
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
