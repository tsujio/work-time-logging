package worktime

import (
	"fmt"

	"work-time-logging/spreadsheet"
)

type WorkTime struct {
	sheet *spreadsheet.Spreadsheet
}

func New(sheet *spreadsheet.Spreadsheet) *WorkTime {
	return &WorkTime{sheet: sheet}
}

func (this *WorkTime) Get(target string, year, month int) ([]WorkTimeRecord, error) {
	sheetName := fmt.Sprintf("%04d%02d", year, month)
	rows := this.sheet.Get(target, sheetName, "A4", "K40")
	records, err := ParseWorkTimeRecords(year, month, rows)
	if err != nil {
		return nil, err
	}
	return records, nil
}
