package cron

import (
	"flag"
	"fmt"
)

func Help(command string, fl *flag.FlagSet) {
	switch command {
	case "help":
		fl.Usage()
	case "daemon":
		helpDaemonCommand()
	case "add":
		helpAddCommand()
	case "remove":
		helpRemoveCommand()
	case "update":
		helpUpdateCommand()
	case "list":
		helpListCommand()
	case "active":
		helpActiveCommand()
	case "next":
		helpNextCommand()
	default:
		fmt.Printf("Cannot describe unknown command: <%s>\n", command)
		fmt.Printf("Available commands: %v\n", Commands)
		fl.Usage()
	}
}

func helpDaemonCommand() {
	var fl  = getDaemonCommandArgsParser()
	fmt.Printf("Start scheduler as daemon, in sync mode\n")
	fl.Usage()
}

func helpAddCommand() {
	var fl  = getAddCommandArgsParser()
	fmt.Printf("Add one configuration item from file or in line encoded text\n")
	fl.Usage()
}

func helpRemoveCommand() {
	var fl  = getRemoveCommandArgsParser()
	fmt.Printf("Removes one or more scheduler configuration items by line a range of lines numbers\n")
	fl.Usage()
}

func helpUpdateCommand() {
	var fl  = getUpdateCommandArgsParser()
	fmt.Printf("Update one configuration item at a specific raw number from file or in line encoded text\n")
	fl.Usage()
}

func helpListCommand() {
	var fl  = getListCommandArgsParser()
	fmt.Printf("List all configuration items, in summary (text table) or detail mode (encoding format)\n")
	fl.Usage()
}

func helpActiveCommand() {
	var fl  = getActiveCommandArgsParser()
	fmt.Printf("List all active processes, in summary (text table) or detail mode (encoding format)\n")
	fl.Usage()
}

func helpNextCommand() {
	var fl  = getNextCommandArgsParser()
	fmt.Printf("List next execution processes, in summary (text table) or detail mode (encoding format)\n")
	fl.Usage()
}

