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
var subCommand string

var Args []string
func init() {
	Args = os.Args[1:]
	cron.Args = os.Args[2:]
}

func main() {
	if len(Args) == 0 {
		fmt.Printf("Missing command argument, available commands: %v\n", cron.Commands)
		cron.PrintHelp(cron.DefaultParser(command))

	}
	command = strings.ToLower(Args[0])
	if b, _ := utils.ListContains(command, cron.Commands); ! b {
		fmt.Printf("Command '%s' not in list : %v\n", command, cron.Commands)
		cron.PrintHelp(cron.DefaultParser(command))
		fl.Usage()
		os.Exit(2)
	}
	if command == "help" {
		if len(Args) < 2 {
			fmt.Printf("Missing command help command argument, available commands: %v\n", cron.Commands[1:])
			cron.PrintHelp(cron.DefaultParser(command))

		}
		subCommand = strings.ToLower(Args[1])
		cron.Args = os.Args[3:]
		cron.Help(command, subCommand)
	} else if command == "explain" {
		if len(Args) < 2 {
			fmt.Printf("Missing command help command argument, available commands: %v\n", cron.Commands[1:])
			cron.PrintHelp(cron.DefaultParser(command))

		}
		subCommand = strings.ToLower(Args[1])
		cron.Args = os.Args[3:]
		cron.Explain(command, subCommand)
	} else {
		err := cron.Exec(command)
		if err != nil {
			cron.LogMany("Error execution command %s, Error: %v", command, err)
		}
	}
}

func testTextSummary() {
	var str = struct {
		Name		string
		Surname		string
		Age			int
		ContactData		interface{}
		Id			int
	}{
		"Fabrizio",
		"Torelli",
		45,
		struct{
			Tel		string
			Email	string
		}{
			"+353834841333",
			"hellgate75@gmail.com",
		},
		1,
	}
	var str2 = struct {
		Name		string
		Surname		string
		Age			int
		ContactData		interface{}
		Id			int
	}{
		"Eleonora",
		"Mori",
		43,
		struct{
			Tel		string
			Email	string
		}{
			"+353838983428",
			"eleonoraleanore77@gmail.com",
		},
		2,
	}
//	data, err := io.EncodeTextFormatSummary(str)
	var in = make([]interface{}, 0)
	in = append(in, str)
	in = append(in, str2)
	data, err := io.EncodeTextFormatSummary(in)
	if err == nil {
		fmt.Println(string(data))
	} else {
		fmt.Printf("Error: %v", err)
	}
	//fmt.Println("Name", io.formatColumn("myNameIsFabrizio"))
}
