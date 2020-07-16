package cron

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-cron/model"
	"github.com/hellgate75/go-cron/utils"
	"strings"
	"sync"
	"time"
)

var itemsLock = make(map[string]map[string]*sync.Mutex)

var NodeMap = make(map[string]interface{})
var ClusterMap = make(map[string]interface{})

func filterFirstExecution(list []model.Execution, match func(model.Execution) bool) *model.Execution {
	for _, c := range list {
		if match(c) {
			return &c
		}
	}
	return nil
}

func sendCommands(c chan model.CommandConfigRef, scheduler0 *scheduler) {
	go func(scheduler *scheduler) {
		for _, com := range scheduler.cacheCommands {
			c <- com
		}
		for _, com := range scheduler.commands {
			c <- com
		}
	}(scheduler0)
}

func checkNextSchedulerTasks(scheduler0 *scheduler) bool {
	var execAtLEastOnce bool
	for _, ref := range scheduler0.cacheCommands {
		if cmd, ok := scheduler0.cache[ref.UUID]; ok {
			exec := scheduler0.ToExecutionWith(ref, cmd)
			if exec.NeedScheduling() {
				execAtLEastOnce = true
			}
		}
	}
	for _, ref := range scheduler0.commands {
		if cmd, err := scheduler0.loadItem(ref.UUID); err == nil {
			exec := scheduler0.ToExecutionWith(ref, *cmd)
			if exec.NeedScheduling() {
				execAtLEastOnce = true
			}
		}
	}
	return execAtLEastOnce
}

func createExecutionContextFrom(execution *model.Execution, scheduler *scheduler) model.ExecutionContext {
	var refs = make([]model.CommandConfigRef, 0)
	refs = append(refs, scheduler.cacheCommands...)
	refs = append(refs, scheduler.commands...)
	return model.ExecutionContext{
		Configuration: &model.SchedulerConfig{
			Commands:    refs,
			Sync: scheduler.syncRun,
		},
		CommandInfo: &execution.Command,
		ContextMap: &execution.Map,
		StaticMap: &NodeMap,
		GlobalMap: &ClusterMap,
		ErrorsPipe: scheduler.errors,
		WarningsPipe: scheduler.warnings,
	}
}

func runTextArrayCommand(scheduler *scheduler, id string, cmdArr []string) {
	out, err := utils.ExecuteCommandArgs(cmdArr...)
	if err != nil {
		scheduler.errors <- err
	} else {
		scheduler.warnings <- errors.New(fmt.Sprintf("Execution of command id : %s, completed, output: %s", id, out))
	}
}

func runTextCommand(scheduler *scheduler, id string, cmd string) {
	out, err := utils.ExecuteCommand(cmd)
	if err != nil {
		scheduler.errors <- err
	} else {
		scheduler.warnings <- errors.New(fmt.Sprintf("Execution of command id : %s, completed, output: %s", id, out))
	}
}

func runFunctionCommand(scheduler *scheduler, id string, execution *model.Execution, function func(model.ExecutionContext) error) {
	var context = createExecutionContextFrom(execution, scheduler)
	err := function(context)
	if err != nil {
		scheduler.errors <- err
	} else {
		scheduler.warnings <- errors.New(fmt.Sprintf("Execution of command id : %s, completed, type: func(model.ExecutionContext) error", id))
	}
}

func runComputableCommand(scheduler *scheduler, id string, execution *model.Execution, computable model.ComputableValue) {
	var context = createExecutionContextFrom(execution, scheduler)
	err := computable.Compute(context)
	if err != nil {
		scheduler.errors <- err
	} else {
		scheduler.warnings <- errors.New(fmt.Sprintf("Execution of command id : %s, completed, type: model.ComputableValue", id))
	}
}

func executeSingleTask(scheduler *scheduler, execution *model.Execution, id string) {
	defer func() {
		var err error
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		if err != nil {
			scheduler.errors <- err
		}
	}()
	// Increase number of executions
	execution.Times++
	var typeOfCommand = fmt.Sprintf("%T", execution.Command.Command)
	switch typeOfCommand {
	case "string":
		cmd := fmt.Sprintf("%v", execution.Command.Command)
		runTextCommand(scheduler, id, cmd)
	case "[]string":
		var cmdArr = execution.Command.Command.([]string)
		runTextArrayCommand(scheduler, id, cmdArr)
	default:
		if strings.Contains(typeOfCommand, "model.ComputableValue") {
			var computable = execution.Command.Command.(model.ComputableValue)
			runComputableCommand(scheduler, id, execution, computable)
		} else if strings.Contains(typeOfCommand, "func") {
			// Try available command
			var function = execution.Command.Command.(func(model.ExecutionContext) error)
			runFunctionCommand(scheduler, id, execution, function)
		} else if strings.Contains(typeOfCommand, "[]") {
			// Slice of something ...
			// We hope model/ComputableValue or func
			//TODO: Implement reflection to seek into the array and collect data, if check of types matches with
			// on one of following: string, []string, model.ComputationValue or func(model.ExecutionContext) error
			// in this case I run the matching elements in the array, as in the other cases
			scheduler.errors <- errors.New(fmt.Sprintf("Unable to execute command of type %s for command id %s, slice of objects execution not implemented yet", typeOfCommand, id))
		} else {
			// Unknown type
			scheduler.errors <- errors.New(fmt.Sprintf("Unable to execute command of type %s", typeOfCommand))
		}
	}
}

func scheduleSingleTask(schedule *scheduler, execution *model.Execution, ref *model.CommandConfigRef) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	if execution.NeedScheduling() {
		go func(schedule *scheduler, execution *model.Execution, id string) {
			defer func() {
				if r := recover(); r != nil {
					schedule.errors <- errors.New(fmt.Sprintf("%v", r))
				} else {
					schedule.warnings <- errors.New(fmt.Sprintf("Scheduler tasks %s completed!!", id))
				}
				execution.Scheduled = false
				//Save with scheduler execution state for non cached tasks
				_ = schedule.saveExecutions()
			}()
			time.Sleep(time.Since(execution.Next))
			executeSingleTask(schedule, execution, id)
		}(schedule, execution, ref.UUID)
		execution.Scheduled = true
		//Save with scheduler execution state for non cached tasks
		_ = schedule.saveExecutions()
	}
	fmt.Printf("%v", execution)
	return err
}

func executeSchedulerTasks(scheduler0 *scheduler) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			scheduler0.errors <- errors.New(fmt.Sprintf("Stopping scheduler due to error: %v", err))
			scheduler0.running = false
		} else {
			scheduler0.warnings <- errors.New(fmt.Sprintf(fmt.Sprint("Scheduler tasks completed ..."), err))
		}
	}()
	var channel = make(chan model.CommandConfigRef)
	sendCommands(channel, scheduler0)
	pool := sync.WaitGroup{}
tasksCycle:
	for {
		select {
		case <- time.After(10 * time.Second):
			break tasksCycle
		case v := <- channel:
			go func(scheduler1 *scheduler, ref model.CommandConfigRef, pool *sync.WaitGroup) {
				var err error
				cmd := scheduler1.cacheValue(ref.UUID)
				if cmd == nil {
					cmd, err = scheduler1.loadItem(ref.UUID)
					if err != nil {
						scheduler1.errors <- err
						return
					}
				}
				if cmd != nil {
					var exec = scheduler1.ToExecutionWith(ref, *cmd)
					if exec != nil {
						if ! scheduler1.IsExecutionStored(exec.UUID) {
							scheduler1.runningTasks = append(scheduler1.runningTasks, *exec)
						}
						if exec.NeedScheduling() {
							pool.Add(1)
							defer func() {
								if r := recover(); r != nil {
									scheduler1.errors <- errors.New(fmt.Sprintf("%v", r))
								}
								pool.Done()
							}()
							err = scheduleSingleTask(scheduler1, exec, &ref)
							if err != nil {
								scheduler1.errors <- err
							}
						}
						// I do not track unscheduled tasks at the moment
						//} else {
						//	scheduler1.warnings <- errors.New(fmt.Sprintf("Task for id: %s already scheduled", ref.UUID))
						//}
					} else {
						scheduler1.errors <- errors.New(fmt.Sprintf("Unable to retrive or create task execution record for task id: %s", ref.UUID))
					}
				} else {
					scheduler1.errors <- errors.New(fmt.Sprintf("Unable to retrive command task with id: %s", ref.UUID))
				}
			}(scheduler0, v, &pool)
		}
	}
	pool.Wait()
	return err
}

