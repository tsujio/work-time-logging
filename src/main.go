package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"strconv"
	"path/filepath"

	"work-time-logging/configuration"
	"work-time-logging/spreadsheet"
	"work-time-logging/worktime"
)

type showCmdArgs struct {
	projectName string
}

type startCmdArgs struct {
	projectName string
	time        string
}

type endCmdArgs struct {
	projectName string
	time        string
}

type travelCmdArgs struct {
	projectName string
	expense int
	note string
}

type linkCmdArgs struct {
	projectName string
}

func doShow(args *showCmdArgs, config *configuration.Config) {
	s := spreadsheet.New(config)
	w := worktime.New(s, config)

	now := time.Now()

	monthlyWorkTime, err := w.Get(args.projectName, now.Year(), int(now.Month()))
	if err != nil {
		log.Fatalf("%+v", err)
	}

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

	formatTravelExpense := func(t *worktime.TravelExpense) string {
		if t == nil {
			return ""
		}
		return fmt.Sprintf("%s  %d 円", t.Note, t.Expense)
	}

	formatDuration := func(d time.Duration, blankIfZero bool) string {
		h := int(d.Hours())
		m := int((d - d.Truncate(time.Hour)).Minutes())
		if blankIfZero && h == 0 && m == 0 {
			return ""
		}
		return fmt.Sprintf("%d:%02d", h, m)
	}

	for _, record := range monthlyWorkTime.Records {
		fmt.Printf("%2d/%2d  %11s  %11s  %11s    %5s | %s\n",
			record.Date.Month,
			record.Date.Day,
			formatPeriod(record.Periods[0]),
			formatPeriod(record.Periods[1]),
			formatPeriod(record.Periods[2]),
			formatDuration(record.GetDuration(), true),
			formatTravelExpense(record.TravelExpense))
	}

	fmt.Println("-------------------------------------------------------------")
	fmt.Printf("Total:                                         %6s\n",
		formatDuration(monthlyWorkTime.GetDuration(), false))
}

func doStart(args *startCmdArgs, config *configuration.Config) {
	s := spreadsheet.New(config)
	w := worktime.New(s, config)

	now := time.Now()

	var t *worktime.Time
	var err error
	if args.time != "" {
		t, err = worktime.ParseHHMM(args.time)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		t = &worktime.Time{now.Hour(), now.Minute()}
	}
	t = t.RoundTime()

	err = w.SetStart(args.projectName,
		&worktime.Date{now.Year(), int(now.Month()), now.Day()},
		t)
	if err != nil {
		log.Fatal(err)
	}
}

func doEnd(args *endCmdArgs, config *configuration.Config) {
	s := spreadsheet.New(config)
	w := worktime.New(s, config)

	now := time.Now()

	var t *worktime.Time
	var err error
	if args.time != "" {
		t, err = worktime.ParseHHMM(args.time)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		t = &worktime.Time{now.Hour(), now.Minute()}
	}
	t = t.RoundTime()

	err = w.SetEnd(args.projectName,
		&worktime.Date{now.Year(), int(now.Month()), now.Day()},
		t)
	if err != nil {
		log.Fatal(err)
	}
}

func doTravel(args *travelCmdArgs, config *configuration.Config) {
	s := spreadsheet.New(config)
	w := worktime.New(s, config)

	now := time.Now()

	err := w.SetTravelExpense(args.projectName,
		&worktime.Date{now.Year(), int(now.Month()), now.Day()},
		args.expense,
		args.note)
	if err != nil {
		log.Fatal(err)
	}
}

func doLink(args *linkCmdArgs, config *configuration.Config) {
	spreadsheetId, err := config.FindSpreadsheetId(args.projectName)
	if err != nil {
		log.Fatalf("%+w", err)
	}

	s := spreadsheet.New(config)
	link := s.GetSpreadsheetLink(spreadsheetId)

	fmt.Println(link)
}

func main() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	configDir := filepath.Join(filepath.Dir(executable), ".work-time-logging")
	config := configuration.Load(configDir)

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	endCmd := flag.NewFlagSet("end", flag.ExitOnError)
	travelCmd := flag.NewFlagSet("travel", flag.ExitOnError)
	linkCmd := flag.NewFlagSet("link", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s COMMAND [ARGS]", os.Args[0])
	}

	switch os.Args[1] {
	case "show":
		showCmd.Parse(os.Args[2:])
		args := showCmdArgs{
			projectName: showCmd.Arg(0),
		}
		doShow(&args, config)
	case "start":
		var args startCmdArgs
		startCmd.StringVar(&args.time, "time", "", "HH:MM")
		startCmd.Parse(os.Args[2:])
		args.projectName = startCmd.Arg(0)
		doStart(&args, config)
	case "end":
		var args endCmdArgs
		endCmd.StringVar(&args.time, "time", "", "HH:MM")
		endCmd.Parse(os.Args[2:])
		args.projectName = endCmd.Arg(0)
		doEnd(&args, config)
	case "travel":
		var args travelCmdArgs
		travelCmd.Parse(os.Args[2:])
		args.projectName = travelCmd.Arg(0)
		expense, err := strconv.Atoi(travelCmd.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		args.expense = expense
		args.note = travelCmd.Arg(2)
		doTravel(&args, config)
	case "link":
		var args linkCmdArgs
		linkCmd.Parse(os.Args[2:])
		args.projectName = linkCmd.Arg(0)
		doLink(&args, config)
	default:
		log.Fatalf("Invalid command: %s", os.Args[1])
	}
}
