package cron

import (
	"flag"
	"fmt"
	"github.com/hellgate75/go-cron/io"
)
var command string
var configPath string
var encoding io.Encoding
var encodingString string

var inputFormat string
var inputFile string
var inputText string
var outputFormat string

var (
	listIndex = 1
	listFrom = 1
	listTo = 1
)

var details bool
var query string
var filter string
var filterFile string

var silent bool

var nativeGobInFile bool
var nativeGobOutFormat bool

var Args = make([]string, 0)

func DefaultParser(commandString string) *flag.FlagSet {
	var fl = flag.NewFlagSet("go-cron", flag.PanicOnError)
	command = commandString
	var defaultEncoding = io.DefaultEncodingFormat
	fl.StringVar(&encodingString,"format", defaultEncoding.String(), fmt.Sprintf("File encoding format (available: %s)", io.EncodingList))
	defaultFile, _ := io.GetDefaultConfigFile(defaultEncoding)
	fl.BoolVar(&silent, "silent", false, "Execute silent command output (explain)")
	fl.StringVar(&configPath,"path", defaultFile, "Configuration file location")
	return fl
}


func printExplainHelp() {
	fmt.Printf("System Scheduler command, for argument details please type command: help <command-name>\n")
	fmt.Printf("Or for command input sample: explain <command-name> -in-form <my-encoding> -out-form <my-encoding> -native-in true|false -native-out true|false\n")
	fmt.Printf("List of available command names: %s\n", Commands)
	fmt.Println()
}


func PrintHelp(fl *flag.FlagSet) {
	fmt.Printf("go-cron <command> -arg0=value0 -arg1=value1  ... -argn=valuen\n")
	printExplainHelp()
	fl.PrintDefaults()
}

func parse(fl *flag.FlagSet) error {
	var err = fl.Parse(Args)
	if err != nil {
		return err
	}
	encoding = io.EncodingFromValue(encodingString)
	return nil
}

func getDaemonCommandArgsParser() *flag.FlagSet {
	var fl  = DefaultParser("daemon")
	return fl
}

func getAddCommandArgsParser() *flag.FlagSet {
	var fl  = DefaultParser("add")
	fl.StringVar(&inputFormat,"in-form", io.DefaultEncodingFormatString, fmt.Sprintf("Input encoding format (available: %s)", io.EncodingList))
	fl.StringVar(&inputFile, "in-file", "", "Input file absolute path")
	fl.StringVar(&inputText,"in-text", "", "Input text value")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormatString, fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	fl.BoolVar(&nativeGobInFile, "native-in", false, "Use native Gob file for input")
	fl.BoolVar(&nativeGobOutFormat, "native-out", false, "Use native Gob format for output")
	return fl
}

func getRemoveCommandArgsParser() *flag.FlagSet {
	var fl  = DefaultParser("remove")
	fl.IntVar(&listIndex,"index", listIndex, "Output list raw line number to be deleted")
	fl.IntVar(&listFrom,"from", listFrom, "Output list raw first line number to be deleted")
	fl.IntVar(&listTo,"to", listTo, "Output list raw last line number to be deleted")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormatString, fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	fl.BoolVar(&nativeGobOutFormat, "native-out", false, "Use native Gob format for output")
	return fl
}

func getUpdateCommandArgsParser() *flag.FlagSet {
	var fl  = DefaultParser("update")
	fl.StringVar(&inputFormat,"in-form", io.DefaultEncodingFormatString, fmt.Sprintf("Input encoding format (available: %s)", io.EncodingList))
	fl.StringVar(&inputFile, "in-file", "", "Input file absolute path")
	fl.StringVar(&inputText,"in-text", "", "Input text value")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormatString, fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	fl.BoolVar(&nativeGobInFile, "native-in", false, "Use native Gob file for input")
	fl.BoolVar(&nativeGobOutFormat, "native-out", false, "Use native Gob format for output")
	fl.IntVar(&listIndex,"index", listIndex, "Output list raw line number of change/replacement")
	return fl
}

func getListCommandArgsParser() *flag.FlagSet {
	var fl  = DefaultParser("list")
	fl.BoolVar(&details, "details", false, "Show details for each scheduler next execution processes, in the requested encoding format")
	fl.StringVar(&query, "query", "", "Comma separated <column name>=<value> keys")
	fl.StringVar(&filter, "filter", "", "Go style template output filter template text")
	fl.StringVar(&filterFile, "filter-file", "", "Go style template output filter template template")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormatString, fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	fl.BoolVar(&nativeGobOutFormat, "native-out", false, "Use native Gob format for output")
	return fl
}

func getActiveCommandArgsParser() *flag.FlagSet {
	var fl  = DefaultParser("active")
	fl.BoolVar(&details, "details", false, "Show details for each scheduler next execution processes, in the requested encoding format")
	fl.StringVar(&query, "query", "", "Comma separated <column name>=<value> keys")
	fl.StringVar(&filter, "filter", "", "Go style template output filter template text")
	fl.StringVar(&filterFile, "filter-file", "", "Go style template output filter template template")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormatString, fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	fl.BoolVar(&nativeGobOutFormat, "native-out", false, "Use native Gob format for output")
	return fl
}

func getNextCommandArgsParser() *flag.FlagSet {
	var fl  = DefaultParser("next")
	fl.BoolVar(&details, "details", false, "Show details for each scheduler next execution processes, in the requested encoding format")
	fl.StringVar(&query, "query", "", "Comma separated <column name>=<value> keys")
	fl.StringVar(&filter, "filter", "", "Go style template output filter template text")
	fl.StringVar(&filterFile, "filter-file", "", "Go style template output filter template template")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormatString, fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	fl.BoolVar(&nativeGobOutFormat, "native-out", false, "Use native Gob format for output")
	return fl
}

