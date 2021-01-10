package spreadsheet

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/sheets/v4"

	"work-time-logging/configuration"
)

type Spreadsheet struct {
	api    *http.Client
	config *configuration.Config
}

func New(config *configuration.Config) *Spreadsheet {
	api := GetAPIClient(config.Dir)
	return &Spreadsheet{api: api, config: config}
}

func (this *Spreadsheet) findSpreadsheetId(name string) string {
	for _, sheet := range this.config.Spreadsheets {
		if sheet.Name == name {
			return sheet.Id
		}
	}
	log.Fatalf("Spreadsheet not found: %s", name)
	return ""
}

func (this *Spreadsheet) Get(spreadsheetName, sheetName, leftUpper, rightBottom string) [][]interface{} {
	srv, err := sheets.New(this.api)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := this.findSpreadsheetId(spreadsheetName)
	readRange := fmt.Sprintf("%s!%s:%s", sheetName, leftUpper, rightBottom)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	return resp.Values
}

func (this *Spreadsheet) Update(spreadsheetName, sheetName, address string, value interface{}) {
	srv, err := sheets.New(this.api)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := this.findSpreadsheetId(spreadsheetName)
	readRange := fmt.Sprintf("%s!%s", sheetName, address)
	vr := sheets.ValueRange{Values: [][]interface{}{[]interface{}{value}}}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, readRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
}
