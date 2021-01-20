package worktime

import (
	"fmt"

	"golang.org/x/xerrors"

	"work-time-logging/configuration"
	"work-time-logging/spreadsheet"
)

type WorkTime struct {
	sheet  *spreadsheet.Spreadsheet
	config *configuration.Config
}

func New(sheet *spreadsheet.Spreadsheet, config *configuration.Config) *WorkTime {
	return &WorkTime{sheet: sheet, config: config}
}

func (this *WorkTime) getSheetName(year, month int) string {
	return fmt.Sprintf("%04d%02d", year, month)
}

func (this *WorkTime) getPeriodCellAddress(recordIndex, periodIndex int, startOrEnd string) (string, error) {
	row := 4 + recordIndex

	col := int('C') + 2*periodIndex
	if startOrEnd == "end" {
		col++
	}
	return fmt.Sprintf("%c%d", rune(col), row), nil
}

func (this *WorkTime) getTravelExpenseCellAddress(recordIndex int) ([]string, error) {
	row := 4 + recordIndex
	return []string{fmt.Sprintf("J%d", row), fmt.Sprintf("K%d", row)}, nil
}

func (this *WorkTime) Get(projectName string, year, month int) (*MonthlyWorkTime, error) {
	spreadsheetId, err := this.config.FindSpreadsheetId(projectName)
	if err != nil {
		return nil, xerrors.Errorf("Unable to find spreadsheet id: %w", err)
	}
	rows, err := this.sheet.Get(spreadsheetId, this.getSheetName(year, month), "A4", "K40")
	if err != nil {
		return nil, xerrors.Errorf("Unable to get sheet data: %w", err)
	}
	monthlyWorkTime, err := parseMonthlyWorkTime(year, month, rows)
	if err != nil {
		return nil, xerrors.Errorf("Unable to parse work time data: %w", err)
	}
	return monthlyWorkTime, nil
}

func (this *WorkTime) SetStart(projectName string, date *Date, time *Time) error {
	monthlyWorkTime, err := this.Get(projectName, date.Year, date.Month)
	if err != nil {
		return err
	}

	recordIndex := -1
	var record *WorkTimeRecord
	for i, rcd := range monthlyWorkTime.Records {
		if date.Equal(rcd.Date) {
			recordIndex = i
			record = &rcd
			break
		}
	}
	if recordIndex == -1 {
		return fmt.Errorf("specified date not found")
	}

	periodIndex := -1
	for i, period := range record.Periods {
		if period.IsEndEmpty() {
			if !period.IsEmpty() {
				return fmt.Errorf("already started")
			}
			periodIndex = i
			break
		}
	}
	if periodIndex == -1 {
		return fmt.Errorf("empty period not found")
	}

	addr, err := this.getPeriodCellAddress(recordIndex, periodIndex, "start")
	if err != nil {
		return err
	}

	spreadsheetId, err := this.config.FindSpreadsheetId(projectName)
	if err != nil {
		return xerrors.Errorf("Unable to find spreadsheet id: %w", err)
	}

	this.sheet.Update(spreadsheetId, this.getSheetName(date.Year, date.Month), addr,
		fmt.Sprintf("%2d:%02d", time.Hour, time.Minute))

	return nil
}

func (this *WorkTime) SetEnd(projectName string, date *Date, time *Time) error {
	monthlyWorkTime, err := this.Get(projectName, date.Year, date.Month)
	if err != nil {
		return err
	}

	recordIndex := -1
	var record *WorkTimeRecord
	for i, rcd := range monthlyWorkTime.Records {
		if date.Equal(rcd.Date) {
			recordIndex = i
			record = &rcd
			break
		}
	}
	if recordIndex == -1 {
		return fmt.Errorf("specified date not found")
	}

	periodIndex := -1
	for i, period := range record.Periods {
		if period.IsEndEmpty() {
			if period.IsEmpty() {
				return fmt.Errorf("not started")
			}
			periodIndex = i
			break
		}
	}
	if periodIndex == -1 {
		return fmt.Errorf("empty period not found")
	}

	addr, err := this.getPeriodCellAddress(recordIndex, periodIndex, "end")
	if err != nil {
		return err
	}

	spreadsheetId, err := this.config.FindSpreadsheetId(projectName)
	if err != nil {
		return xerrors.Errorf("Unable to find spreadsheet id: %w", err)
	}

	this.sheet.Update(spreadsheetId, this.getSheetName(date.Year, date.Month), addr,
		fmt.Sprintf("%2d:%02d", time.Hour, time.Minute))

	return nil
}

func (this *WorkTime) SetTravelExpense(projectName string, date *Date, expense int, note string) error {
	monthlyWorkTime, err := this.Get(projectName, date.Year, date.Month)
	if err != nil {
		return err
	}

	recordIndex := -1
	for i, rcd := range monthlyWorkTime.Records {
		if date.Equal(rcd.Date) {
			recordIndex = i
			break
		}
	}
	if recordIndex == -1 {
		return fmt.Errorf("specified date not found")
	}

	addrList, err := this.getTravelExpenseCellAddress(recordIndex)
	if err != nil {
		return err
	}

	spreadsheetId, err := this.config.FindSpreadsheetId(projectName)
	if err != nil {
		return xerrors.Errorf("Unable to find spreadsheet id: %w", err)
	}

	this.sheet.Update(spreadsheetId, this.getSheetName(date.Year, date.Month),
		addrList[0],
		note)
	this.sheet.Update(spreadsheetId, this.getSheetName(date.Year, date.Month),
		addrList[1],
		expense)

	return nil
}
