package spreadsheet

import (
	"fmt"
	"net/http"

	"google.golang.org/api/sheets/v4"
	"golang.org/x/xerrors"

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

func (this *Spreadsheet) findSpreadsheetId(name string) (string, error) {
	for _, sheet := range this.config.Spreadsheets {
		if sheet.Name == name {
			return sheet.Id, nil
		}
	}
	return "", xerrors.Errorf("Spreadsheet not found: %s", name)
}

func (this *Spreadsheet) Get(spreadsheetName, sheetName, leftUpper, rightBottom string) ([][]interface{}, error) {
	srv, err := sheets.New(this.api)
	if err != nil {
		return nil, xerrors.Errorf("Unable to retrieve Sheets client: %w", err)
	}

	spreadsheetId, err := this.findSpreadsheetId(spreadsheetName)
	if err != nil {
		return nil, xerrors.Errorf("Unable to find sheet id: %w", err)
	}
	readRange := fmt.Sprintf("%s!%s:%s", sheetName, leftUpper, rightBottom)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, xerrors.Errorf("Unable to retrieve data from sheet: %w", err)
	}

	return resp.Values, nil
}

func (this *Spreadsheet) Update(spreadsheetName, sheetName, address string, value interface{}) error {
	srv, err := sheets.New(this.api)
	if err != nil {
		return xerrors.Errorf("Unable to retrieve Sheets client: %w", err)
	}

	spreadsheetId, err := this.findSpreadsheetId(spreadsheetName)
	if err != nil {
		return xerrors.Errorf("Unable to find sheet id: %w", err)
	}
	readRange := fmt.Sprintf("%s!%s", sheetName, address)
	vr := sheets.ValueRange{Values: [][]interface{}{[]interface{}{value}}}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, readRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return xerrors.Errorf("Unable to retrieve data from sheet: %w", err)
	}

	return nil
}
