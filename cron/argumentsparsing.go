package cron

import (
	"flag"
	"fmt"
	"github.com/hellgate75/go-cron/io"
	"os"
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

func defaultParser(command string) *flag.FlagSet {
	var fl = flag.NewFlagSet("go-cron", flag.ContinueOnError)

	fl.StringVar(&command, "command", command, fmt.Sprintf("Requested %s command", command))
	var defaultEncoding = io.DefaultEncodingFormat
	fl.StringVar(&encodingString,"format", defaultEncoding.String(), fmt.Sprintf("File encoding format (available: %s)", io.EncodingList))
	defaultFile, _ := io.GetDefaultConfigFile(defaultEncoding)
	fl.StringVar(&configPath,"path", defaultFile, "Configuration file location")
	return fl
}

func parse(fl *flag.FlagSet) error {
	var err = fl.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	encoding = io.EncodingFromValue(encodingString)
	return nil
}

func getDaemonCommandArgsParser() *flag.FlagSet {
	var fl  = defaultParser("daemon")
	return fl
}

func getAddCommandArgsParser() *flag.FlagSet {
	var fl  = defaultParser("add")
	fl.StringVar(&inputFormat,"in-form", io.DefaultEncodingFormat.String(), fmt.Sprintf("Input encoding format (available: %s)", io.EncodingList))
	fl.StringVar(&inputFile, "in-file", "", "Input file absolute path")
	fl.StringVar(&inputText,"in-text", "", "Input text value")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormat.String(), fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	return fl
}

func getRemoveCommandArgsParser() *flag.FlagSet {
	var fl  = defaultParser("remove")
	fl.IntVar(&listIndex,"index", listIndex, "Output list raw line number to be deleted")
	fl.IntVar(&listFrom,"from", listFrom, "Output list raw first line number to be deleted")
	fl.IntVar(&listTo,"to", listTo, "Output list raw last line number to be deleted")
	return fl
}

func getUpdateCommandArgsParser() *flag.FlagSet {
	var fl  = defaultParser("update")
	fl.StringVar(&inputFormat,"in-form", io.DefaultEncodingFormat.String(), fmt.Sprintf("Input encoding format (available: %s)", io.EncodingList))
	fl.StringVar(&inputFile, "in-file", "", "Input file absolute path")
	fl.StringVar(&inputText,"in-text", "", "Input text value")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormat.String(), fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	fl.IntVar(&listIndex,"index", listIndex, "Output list raw line number of change/replacement")
	return fl
}

func getListCommandArgsParser() *flag.FlagSet {
	var fl  = defaultParser("list")
	fl.BoolVar(&details, "details", details, "Show details for each scheduler configuration item, in the requested encoding format")
	fl.StringVar(&query, "query", query, "Comma separated <column name>=<value> keys")
	fl.StringVar(&filter, "filter", filter, "Go style template output filter template text")
	fl.StringVar(&filterFile, "filter-file", filterFile, "Go style template output filter template template")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormat.String(), fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	return fl
}

func getActiveCommandArgsParser() *flag.FlagSet {
	var fl  = defaultParser("active")
	fl.BoolVar(&details, "details", details, "Show details for each scheduler active processes, in the requested encoding format")
	fl.StringVar(&query, "query", query, "Comma separated <column name>=<value> keys")
	fl.StringVar(&filter, "filter", filter, "Go style template output filter template text")
	fl.StringVar(&filterFile, "filter-file", filterFile, "Go style template output filter template template")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormat.String(), fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	return fl
}

func getNextCommandArgsParser() *flag.FlagSet {
	var fl  = defaultParser("next")
	fl.BoolVar(&details, "details", details, "Show details for each scheduler next execution processes, in the requested encoding format")
	fl.StringVar(&query, "query", query, "Comma separated <column name>=<value> keys")
	fl.StringVar(&filter, "filter", filter, "Go style template output filter template text")
	fl.StringVar(&filterFile, "filter-file", filterFile, "Go style template output filter template template")
	fl.StringVar(&outputFormat, "out-form", io.DefaultEncodingFormat.String(), fmt.Sprintf("Output encoding format (available: text, %s)", io.EncodingList))
	return fl
}

