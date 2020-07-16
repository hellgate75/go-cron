package cron

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hellgate75/go-cron/io"
	"github.com/hellgate75/go-cron/model"
	"time"
)

func Explain(command string, subCommand string) {
	switch subCommand {
	case "help":
		PrintHelp(DefaultParser("help"))
	case "add":
		explainAddCommand()
	case "remove":
		explainRemoveCommand()
	case "update":
		explainUpdateCommand()
	case "list":
		explainListCommand()
	case "active":
		explainActiveCommand()
	case "next":
		explainNextCommand()
	default:
		fmt.Printf("Cannot explain unknown data command: <%s>\n", command)
		fmt.Printf("Available commands: %v\n", Commands[2:])
		PrintHelp(DefaultParser(command))
	}
}

func explainAddCommand() {
	_ = parse(getAddCommandArgsParser())
	p := 10 * time.Hour
	c := model.CommandConfig{
		Command: "myCommand myArg1 myArg2 ...",
		Repeat: 5,
		Since: time.Now(),
		Period: p.String(),
		OnDemand: false,
	}
	if ! silent {
		fmt.Printf("Add one configuration item: \n")
	}
	if nativeGobInFile {
		b, _ := io.EncodeGobValue(&c)
		LogText(string(b))
	} else {
		var inputEncoding = io.EncodingFromValue(inputFormat)
		b, _ := io.EncodeValue(&c, inputEncoding)
		LogText(string(b))
	}
}

func explainRemoveCommand() {
	_ = parse(getRemoveCommandArgsParser())
	if ! silent {
		fmt.Printf("Delete one or more configuration items, including the index or initial and end indexes from the configuration list: \n")
	}
	fmt.Printf("No prototype data \n")
}

func explainUpdateCommand() {
	_ = parse(getUpdateCommandArgsParser())
	p := 10 * time.Hour
	c := model.CommandConfig{
		Command: "myCommand myArg1 myArg2 ...",
		Repeat: 5,
		Since: time.Now(),
		Period: p.String(),
		OnDemand: false,
	}
	if ! silent {
		fmt.Printf("Update one configuration item, including the index from the configuration list: \n")
	}
	if nativeGobInFile {
		b, _ := io.EncodeGobValue(&c)
		LogText(string(b))
	} else {
		var inputEncoding = io.EncodingFromValue(inputFormat)
		b, _ := io.EncodeValue(&c, inputEncoding)
		LogText(string(b))
	}
}

func explainListCommand() {
	_ = parse(getListCommandArgsParser())
	p := 10 * time.Hour
	c := model.CommandConfig{
		Command: "myCommand myArg1 myArg2 ...",
		Repeat: 5,
		Since: time.Now(),
		Period: p.String(),
		OnDemand: false,
	}
	if ! silent {
		fmt.Printf("List all configuration items, in summary (text table) or detail mode (encoding format)\n")
		fmt.Printf("Output sample:\n")
	}
	_ = executeListCommand(false, c)
}


func explainActiveCommand() {
	_ = parse(getActiveCommandArgsParser())
	p := 20 * time.Minute
	p2 := 1 * time.Minute
	c := model.CommandConfig{
		Command: "myCommand myArg1 myArg2 ...",
		Repeat: 5,
		Since: time.Now().Add(-p2),
		Period: p.String(),
		OnDemand: false,
	}
	e := model.Execution{
		Command: c,
		Times: 5,
		Scheduled: false,
		UUID: uuid.New().String(),
		Last: time.Now().Add(- p - p2),
		Next: time.Now().Add(- p2),
		Map: make(map[string]interface{}),
	}

	if ! silent {
		fmt.Printf("List next execution processes, in summary (text table) or detail mode (encoding format)\n")
		fmt.Printf("Output sample:\n")
	}
	_ = executeActiveCommand(false, e)
}


func explainNextCommand() {
	_ = parse(getNextCommandArgsParser())
	p := 20 * time.Minute
	c := model.CommandConfig{
		Command: "myCommand myArg1 myArg2 ...",
		Repeat: 5,
		Since: time.Now(),
		Period: p.String(),
		OnDemand: false,
	}
	e := model.Execution{
		Command: c,
		Times: 1,
		Scheduled: true,
		UUID: uuid.New().String(),
		Last: time.Now().Add(- p),
		Next: time.Now(),
		Map: make(map[string]interface{}),
	}
	if ! silent {
		fmt.Printf("List all active processes, in summary (text table) or detail mode (encoding format)\n")
		fmt.Printf("Output sample:\n")
	}
	_ = executeNextCommand(false, e)
}

