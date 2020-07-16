package cron

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-cron/io"
	"github.com/hellgate75/go-cron/model"
	"os"
	"strings"
	"time"
)

var Commands = []string{"help", "explain", "daemon", "add", "remove", "update", "list", "active", "next"}

func LogObj(d interface{}) {
	fmt.Printf("%v\n", d)
}

func LogText(s string) {
	fmt.Printf("%s\n", s)
}

func LogMany(format string, d ...interface{}) {
	fmt.Printf(format, d)
}


func LogResponse(err error, message string, out interface{}) {
	var response = struct {
		Error	error				`yaml:"error,omitempty" json:"error,omitempty" xml:"error,omitempty"`
		Message string				`yaml:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
		Content	interface{}			`yaml:"content,omitempty" json:"content,omitempty" xml:"content,omitempty"`
	}{
		err,
		message,
		out,
	}
	var data []byte
	if nativeGobOutFormat {
		data, _ = io.EncodeGobValue(&response)
	} else {
		if strings.ToLower(outputFormat) == "text" {
			data, _ = io.EncodeTextFormatSummary(response)
		} else {
			var outputEncoding = io.EncodingFromValue(outputFormat)
			data, _ = io.EncodeValue(&response, outputEncoding)
		}
	}
	LogText(string(data))
}

func LogListResponse(message string, out interface{}) {
	var response = struct {
		Title 		string				`yaml:"title,omitempty" json:"title,omitempty" xml:"title,omitempty"`
		Response	interface{}			`yaml:"response,omitempty" json:"response,omitempty" xml:"response,omitempty"`
	}{
		fmt.Sprintf("List of %s\n", message),
		out,
	}
	var data []byte
	if nativeGobOutFormat {
		data, _ = io.EncodeGobValue(&response)
	} else {
		if strings.ToLower(outputFormat) == "text" {
			data, _ = io.EncodeTextFormatSummary(out)
			var tmp = []byte(fmt.Sprintf("List of %s\n", message))
			tmp = append(tmp, data...)
			data = tmp
		} else {
			var outputEncoding = io.EncodingFromValue(outputFormat)
			data, _ = io.EncodeValue(&response, outputEncoding)
		}
	}
	LogText(string(data))
}

func Exec(command string) error {
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
		return executeListCommand(true)
	case "active":
		return executeActiveCommand(true)
	case "next":
		return executeNextCommand(true)
	default:
		LogMany("Cannot describe unknown command: <%s>\n", command)
		LogMany("Available commands: %v\n", Commands)
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
	//TODO: Implement logger for the daemon
	err = parse(getDaemonCommandArgsParser())
	if err != nil {
		return err
	}
	var scheduler model.Scheduler
	if configPath == "" {
		configPath = fmt.Sprintf("%s%c%s/%s.%s", io.HomeFolder(), os.PathSeparator, ",go-cron", "scheduler", encoding.String())
		scheduler, err = NewEmptyScheduler(configPath, encoding, true)
		if err != nil {
			return err
		}
	} else {
		scheduler, err = LoadSchedulerFrom(configPath, encoding, true)
		if err != nil {
			return err
		}
		err = scheduler.Load()
		if err != nil {
			return err
		}
	}
	err = scheduler.Start()
	if err != nil {
		return err
	}
	scheduler.Wait()
	return err
}

func executeAddCommand() error {
	var err error
	var scheduler model.Scheduler
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	err = parse(getAddCommandArgsParser())
	if err != nil {
		return err
	}
	if configPath == ""  || encoding.String() == "" {
		err = errors.New(fmt.Sprint("Invalid parameters"))
	} else {
		var inputCommand model.CommandConfig
		if nativeGobOutFormat {
			if inputText != "" {
				err = io.DecodeGobValue(&inputCommand, []byte(inputText))
			} else if inputFile != "" {
				err = io.ReadNative(inputFile, &inputCommand)
			} else {
				err = errors.New(fmt.Sprint("Invalid input source"))
			}
		} else {
			var inputEncoding = io.EncodingFromValue(inputFormat)
			if inputEncoding.String() == "" {
				err = errors.New(fmt.Sprint("Invalid input encoding format"))
			}
			if inputText != "" {
				err = io.DecodeValue(&inputCommand, []byte(inputText), inputEncoding)
			} else if inputFile != "" {
				err = io.ReadConfig(inputEncoding, inputFile, &inputCommand)
			} else {
				err = errors.New(fmt.Sprint("Invalid input source"))
			}
		}
		if err == nil {
			scheduler, err = LoadSchedulerFrom(configPath, encoding, true)
			if err != nil {
				return err
			}
			err = scheduler.Load()
			if err != nil {
				return err
			}
			err = scheduler.AddAndPersist(inputCommand)
			LogResponse(err, "Adding new command", inputCommand)
			err = nil
		}
	}
	return err
}

func executeRemoveCommand() error {
	var err error
	var scheduler model.Scheduler
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	err = parse(getRemoveCommandArgsParser())
	if err != nil {
		return err
	}
	if configPath == ""  || encoding.String() == "" {
		err = errors.New(fmt.Sprint("Invalid parameters"))
	} else {
		scheduler, err = LoadSchedulerFrom(configPath, encoding, true)
		if err != nil {
			return err
		}
		err = scheduler.Load()
		if err != nil {
			return err
		}

		if listFrom != listTo {
			if listFrom > listTo {
				var outRes = make([]interface{}, 0)
				//	Items range removal
				for i := listFrom; i <= listTo; i++ {
					//	Single item removal
					err = scheduler.DeleteAndPersist(i)
					outRes = append(outRes, struct{
						Index	int					`yaml:"index,omitempty" json:"index,omitempty" xml:"index,omitempty"`
						Error 	error				`yaml:"error,omitempty" json:"error,omitempty" xml:"error,omitempty"`
					}{
						i,
						err,
					})
				}
				LogResponse(err, "Removing command in range of indexes", outRes)
			} else {
				//	Items range removal
				for i := listTo; i <= listFrom; i++ {
					//	Single item removal
					err = scheduler.DeleteAndPersist(i)
					LogResponse(err, "Removing command at index", fmt.Sprintf("Removing item at index %v", i))
				}
			}
		} else {
			//	Single item removal
			err = scheduler.DeleteAndPersist(listIndex)
			LogResponse(err, "Removing command at index", fmt.Sprintf("Removing item at index %v", listIndex))
		}
		err = nil
	}
	return err
}


func executeUpdateCommand() error {
	var err error
	var scheduler model.Scheduler
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	err = parse(getUpdateCommandArgsParser())
	if err != nil {
		return err
	}
	if configPath == ""  || encoding.String() == "" {
		err = errors.New(fmt.Sprint("Invalid parameters"))
	} else {
		var inputCommand model.CommandConfig
		if nativeGobOutFormat {
			if inputText != "" {
				err = io.DecodeGobValue(&inputCommand, []byte(inputText))
			} else if inputFile != "" {
				err = io.ReadNative(inputFile, &inputCommand)
			} else {
				err = errors.New(fmt.Sprint("Invalid input source, only file is accepted for native decryption"))
			}
		} else {
			var inputEncoding = io.EncodingFromValue(inputFormat)
			if inputEncoding.String() == "" {
				err = errors.New(fmt.Sprint("Invalid input encoding format"))
			}
			if inputText != "" {
				err = io.DecodeValue(&inputCommand, []byte(inputText), inputEncoding)
			} else if inputFile != "" {
				err = io.ReadConfig(inputEncoding, inputFile, &inputCommand)
			} else {
				err = errors.New(fmt.Sprint("Invalid input source"))
			}
		}
		if err == nil {
			scheduler, err = LoadSchedulerFrom(configPath, encoding, true)
			if err != nil {
				return err
			}
			err = scheduler.Load()
			if err != nil {
				return err
			}
			err = scheduler.UpdateAndPersist(inputCommand, listIndex)
			LogResponse(err, fmt.Sprintf("Update command at index %v", listIndex), inputCommand)
			err = nil
		}
	}
	return err
}


func executeListCommand(parseArgs bool, configList ...model.CommandConfig) error {
	var err error
	var scheduler model.Scheduler
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	if parseArgs {
		err = parse(getListCommandArgsParser())
		if err != nil {
			return err
		}
	}
	if len(configList) == 0 && (configPath == ""  || encoding.String() == "") {
		err = errors.New(fmt.Sprint("Invalid parameters"))
	} else {
		var list = make([]model.CommandConfig, 0)
		if len(configList) == 0 {
			scheduler, err = LoadSchedulerFrom(configPath, encoding, true)
			if err != nil {
				return err
			}
			err = scheduler.Load()
			if err != nil {
				return err
			}
			list = scheduler.Planned()
		} else {
			list = configList
		}
		//TODO: Implement query and filter features
		if details {
			var newList = make([]interface{}, 0)
			for idx, r := range list {
				newList = append(newList, struct{
					Line		int					`yaml:"line,omitempty" json:"line,omitempty" xml:"line,omitempty"`
					Command		model.CommandConfig `yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
				}{
					idx,
					r,
				})
			}
			LogListResponse("Planned Tasks", newList)
		} else {
			var newList = make([]interface{}, 0)
			for idx, r := range list {
				cmd := r.Command
				newList = append(newList, struct{
					Line		int					`yaml:"line,omitempty" json:"line,omitempty" xml:"line,omitempty"`
					Command		model.CommandValue `yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
				}{
					idx,
					cmd,
				})
			}
			LogListResponse("Planned Tasks", newList)
		}
	}
	return err
}

func executeActiveCommand(parseArgs bool, configList ...model.Execution) error {
	var err error
	var scheduler model.Scheduler
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	if parseArgs {
		err = parse(getActiveCommandArgsParser())
	}
	if err != nil {
		return err
	}
	if len(configList) == 0 && (configPath == ""  || encoding.String() == "") {
		err = errors.New(fmt.Sprint("Invalid parameters"))
	} else {
		var list = make([]model.Execution, 0)
		if len(configList) == 0 {
			scheduler, err = LoadSchedulerFrom(configPath, encoding, true)
			if err != nil {
				return err
			}
			err = scheduler.Load()
			if err != nil {
				return err
			}
			list = scheduler.Running()
		} else {
			list = configList
		}
		if details {
			var newList = make([]interface{}, 0)
			for idx, r := range list {

				newList = append(newList, struct{
					Line		int					 `yaml:"line,omitempty" json:"line,omitempty" xml:"line,omitempty"`
					Uuid		string				 `yaml:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
					LastExec	time.Time			 `yaml:"lastExecution,omitempty" json:"lastExecution,omitempty" xml:"last-execution,omitempty"`
					Scheduled	bool				 `yaml:"isScheduled,omitempty" json:"isScheduled,omitempty" xml:"is-scheduled,omitempty"`
					NextExec	time.Time			 `yaml:"nextExecution,omitempty" json:"nextExecution,omitempty" xml:"next-execution,omitempty"`
					NoRuns		int					 `yaml:"numberOfExecutions,omitempty" json:"numberOfExecutions,omitempty" xml:"number-of-execution,omitempty"`
					Command		model.CommandConfig  `yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
				}{
					idx,
					r.UUID,
					r.Last,
					r.Scheduled,
					r.Next,
					r.Times,
					r.Command,
				})
			}
			LogListResponse("Active Tasks", newList)
		} else {
			var newList = make([]interface{}, 0)
			for idx, r := range list {
				cmd := r.Command.Command
				newList = append(newList, struct{
					Line		int					`yaml:"line,omitempty" json:"line,omitempty" xml:"line,omitempty"`
					Uuid		string				 `yaml:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
					LastExec	time.Time			 `yaml:"lastExecution,omitempty" json:"lastExecution,omitempty" xml:"last-execution,omitempty"`
					Scheduled	bool				 `yaml:"isScheduled,omitempty" json:"isScheduled,omitempty" xml:"is-scheduled,omitempty"`
					Command		model.CommandValue `yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
				}{
					idx,
					r.UUID,
					r.Last,
					r.Scheduled,
					cmd,
				})
			}
			LogListResponse("Active Tasks", newList)
		}
	}
	return err
}

func executeNextCommand(parseArgs bool, configList ...model.Execution) error {
	var err error
	var scheduler model.Scheduler
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	if parseArgs {
		err = parse(getNextCommandArgsParser())
		if err != nil {
			return err
		}
	}
	if len(configList) == 0 && (configPath == ""  || encoding.String() == "") {
		err = errors.New(fmt.Sprint("Invalid parameters"))
	} else {
		var list = make([]model.Execution, 0)
		if len(configList) == 0 {
			scheduler, err = LoadSchedulerFrom(configPath, encoding, true)
			if err != nil {
				return err
			}
			err = scheduler.Load()
			if err != nil {
				return err
			}
			list = scheduler.NextRunningTasks()
		} else {
			list = configList
		}
		if details {
			var newList = make([]interface{}, 0)
			for idx, r := range list {

				newList = append(newList, struct{
					Line		int					 `yaml:"line,omitempty" json:"line,omitempty" xml:"line,omitempty"`
					Uuid		string				 `yaml:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
					LastExec	time.Time			 `yaml:"lastExecution,omitempty" json:"lastExecution,omitempty" xml:"last-execution,omitempty"`
					Scheduled	bool				 `yaml:"isScheduled,omitempty" json:"isScheduled,omitempty" xml:"is-scheduled,omitempty"`
					NextExec	time.Time			 `yaml:"nextExecution,omitempty" json:"nextExecution,omitempty" xml:"next-execution,omitempty"`
					NoRuns		int					 `yaml:"numberOfExecutions,omitempty" json:"numberOfExecutions,omitempty" xml:"number-of-execution,omitempty"`
					Command		model.CommandConfig  `yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
				}{
					idx,
					r.UUID,
					r.Last,
					r.Scheduled,
					r.Next,
					r.Times,
					r.Command,
				})
			}
			LogListResponse("Next Execution Tasks", newList)
		} else {
			var newList = make([]interface{}, 0)
			for idx, r := range list {
				cmd := r.Command.Command
				newList = append(newList, struct{
					Line		int					`yaml:"line,omitempty" json:"line,omitempty" xml:"line,omitempty"`
					Uuid		string				 `yaml:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
					LastExec	time.Time			 `yaml:"lastExecution,omitempty" json:"lastExecution,omitempty" xml:"last-execution,omitempty"`
					Scheduled	bool				 `yaml:"isScheduled,omitempty" json:"isScheduled,omitempty" xml:"is-scheduled,omitempty"`
					Command		model.CommandValue `yaml:"command,omitempty" json:"command,omitempty" xml:"command,omitempty"`
				}{
					idx,
					r.UUID,
					r.Last,
					r.Scheduled,
					cmd,
				})
			}
			LogListResponse("Next Execution Tasks", newList)
		}
	}
	return err
}
