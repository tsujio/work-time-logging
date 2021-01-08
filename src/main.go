package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/sheets/v4"

	"work-time-logging/spreadsheet"
)

func doShow() {
	client := spreadsheet.GetAPIClient(".")

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := ""
	readRange := ""
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			//			fmt.Printf("%s, %s\n", row[0], row[1])
			fmt.Println(row)
		}
	}
}

func main() {
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s COMMAND [ARGS]", os.Args[0])
	}

	switch os.Args[1] {
	case "show":
		showCmd.Parse(os.Args[1:])
		doShow()
	default:
		log.Fatalf("Invalid command: %s", os.Args[1])
	}
}
