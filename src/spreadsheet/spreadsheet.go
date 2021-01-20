package spreadsheet

import (
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
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

func (this *Spreadsheet) GetSpreadsheetLink(spreadsheetId string) string {
	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", spreadsheetId)
}

func (this *Spreadsheet) Get(spreadsheetId, sheetName, leftUpper, rightBottom string) ([][]interface{}, error) {
	srv, err := sheets.New(this.api)
	if err != nil {
		return nil, xerrors.Errorf("Unable to retrieve Sheets client: %w", err)
	}

	readRange := fmt.Sprintf("%s!%s:%s", sheetName, leftUpper, rightBottom)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, xerrors.Errorf("Unable to retrieve data from sheet: %w", err)
	}

	return resp.Values, nil
}

func (this *Spreadsheet) Update(spreadsheetId, sheetName, address string, value interface{}) error {
	srv, err := sheets.New(this.api)
	if err != nil {
		return xerrors.Errorf("Unable to retrieve Sheets client: %w", err)
	}

	updateRange := fmt.Sprintf("%s!%s", sheetName, address)
	vr := sheets.ValueRange{Values: [][]interface{}{[]interface{}{value}}}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, updateRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return xerrors.Errorf("Unable to retrieve data from sheet: %w", err)
	}

	return nil
}
