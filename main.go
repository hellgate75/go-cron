package main

import (
	"flag"
	"fmt"
	"github.com/hellgate75/go-cron/cron"
	"github.com/hellgate75/go-cron/io"
	"github.com/hellgate75/go-cron/utils"
	"os"
	"strings"
)

var fl *flag.FlagSet
var command string
var help string
var configPath string
var encoding io.Encoding
var encodingString string

func init() {
	fl = flag.NewFlagSet("go-cron", flag.ContinueOnError)
	fl.StringVar(&command, "command", "", fmt.Sprintf("Requested Command (available: %v)", cron.Commands))
	fl.StringVar(&help, "help", "", fmt.Sprintf("Help for Command (available: %v)", cron.Commands[1:]))
	var defaultEncoding = io.DefaultEncodingFormat
	fl.StringVar(&encodingString, "format", defaultEncoding.String(), fmt.Sprintf("File encoding format (available: %s)", io.EncodingList))
	defaultFile, _ := io.GetDefaultConfigFile(defaultEncoding)
	fl.StringVar(&configPath, "path", defaultFile, "Configuration file location")
}

func printExplainHelp() {
	fmt.Printf("System Scheduler command, for argument details please type command: help <command-name>\n")
	fmt.Printf("List of available command names: %s\n", cron.Commands)
	fmt.Println()
}

func main() {
	var err = fl.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Error occured during parse: %v\n", err)
		printExplainHelp()
		fl.Usage()
		os.Exit(1)
	}
	if b, _ := utils.ListContains(command, cron.Commands); ! b {
		fmt.Printf("Command '%s' not in list : %v\n", command, cron.Commands)
		printExplainHelp()
		fl.Usage()
		os.Exit(2)
	}
	if encoding = io.EncodingFromValue(encodingString); encoding != io.EncodingUnknown {
		var commandArg = strings.ToLower(command)
		if "help" == commandArg {
			if "help" == strings.ToLower(help) {
				printExplainHelp()
			}
			cron.Help(help, fl)
		} else {
			err = cron.Exec(commandArg, configPath, encoding)
			if err != nil {
				fmt.Printf("Error during execution of the command <%s>\n", commandArg)
			} else {
				fmt.Printf("Command %s execution completed!!\n", commandArg)
			}
		}
	} else {
		fmt.Printf("Encoding '%s' unknown\n", encodingString)
		printExplainHelp()
		fl.Usage()
		os.Exit(3)
	}
}