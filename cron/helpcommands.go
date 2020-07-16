package cron

import (
	"fmt"
)

func Help(command string, subCommand string) {
	switch subCommand {
	case "help":
		helpHelpHelp()
	case "explain":
		helpExplainHelp()
	case "daemon":
		helpDaemonCommand()
	case "once":
		helpOnceCommand()
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
		PrintHelp(DefaultParser(command))
	}
}

func helpHelpHelp() {
	var fl  = DefaultParser(command)
	fmt.Printf("Request Help for %v commands.\n", Commands)
	fl.Usage()
}

func helpExplainHelp() {
	var fl  = DefaultParser(command)
	fmt.Printf("Request Help explain for %v commands.\n", Commands[4:])
	fl.Usage()
}

func helpDaemonCommand() {
	var fl  = getDaemonCommandArgsParser()
	fmt.Printf("Start scheduler as daemon, in sync mode\n")
	fl.Usage()
}

func helpOnceCommand() {
	var fl  = getOnceCommandArgsParser()
	fmt.Printf("Start scheduler as one-shot run, in sync mode\n")
	fl.Usage()
}

func helpAddCommand() {
	var fl  = getAddCommandArgsParser()
	fmt.Printf("Add one configuration item from file or in line encoded text\n")
	PrintHelp(fl)
}

func helpRemoveCommand() {
	var fl  = getRemoveCommandArgsParser()
	fmt.Printf("Removes one or more scheduler configuration items by line a range of lines numbers\n")
	PrintHelp(fl)
}

func helpUpdateCommand() {
	var fl  = getUpdateCommandArgsParser()
	fmt.Printf("Update one configuration item at a specific raw number from file or in line encoded text\n")
	PrintHelp(fl)
}

func helpListCommand() {
	var fl  = getListCommandArgsParser()
	fmt.Printf("List all configuration items, in summary (text table) or detail mode (encoding format)\n")
	PrintHelp(fl)
}

func helpActiveCommand() {
	var fl  = getActiveCommandArgsParser()
	fmt.Printf("List all active processes, in summary (text table) or detail mode (encoding format)\n")
	PrintHelp(fl)
}

func helpNextCommand() {
	var fl  = getNextCommandArgsParser()
	fmt.Printf("List next execution processes, in summary (text table) or detail mode (encoding format)\n")
	PrintHelp(fl)
}

