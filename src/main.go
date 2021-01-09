package main

import (
	"flag"
	"log"
	"os"

	"work-time-logging/spreadsheet"
	"work-time-logging/worktime"
	"work-time-logging/configuration"
)

type ShowCmdArgs struct {
	SheetName string
}

func doShow(args *ShowCmdArgs, config *configuration.Config) {
	s := spreadsheet.New(config)
	w := worktime.New(s)
	w.Get(args.SheetName, 2021, 1)
}

func main() {
	config := configuration.Load(".")

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s COMMAND [ARGS]", os.Args[0])
	}

	switch os.Args[1] {
	case "show":
		var args ShowCmdArgs
		showCmd.Parse(os.Args[2:])
		args.SheetName = showCmd.Arg(0)
		doShow(&args, config)
	default:
		log.Fatalf("Invalid command: %s", os.Args[1])
	}
}
