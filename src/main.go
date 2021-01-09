package main

import (
	"flag"
	"log"
	"os"
	"fmt"

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
	records, err := w.Get(args.SheetName, 2021, 1)
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		formatPeriod := func(p worktime.Period) string {
			if p.IsEmpty() {
				return ""
			}
			startStr := fmt.Sprintf("%2d:%02d", p.Start.Hour(), p.Start.Minute())
			if p.IsEndEmpty() {
				return startStr + "-"
			}
			endHour := p.End.Hour()
			if p.Start.Day() != p.End.Day() {
				endHour += 24
			}
			endStr := fmt.Sprintf("%2d:%02d", endHour, p.End.Minute())
			return startStr + "-" + endStr
		}

		fmt.Printf("%2d/%2d  %11s  %11s  %11s\n",
			record.Date.Month,
			record.Date.Day,
			formatPeriod(record.Periods[0]),
			formatPeriod(record.Periods[1]),
			formatPeriod(record.Periods[2]))
	}
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
