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

func (this *WorkTime) getSheetName(year, month int) string {
	return fmt.Sprintf("%04d%02d", year, month)
}

func (this *WorkTime) getPeriodCellAddress(recordIndex, periodIndex int, startOrEnd string) (string, error) {
	row := 4 + recordIndex

	col := int('C') + 2 * periodIndex
	if startOrEnd == "end" {
		col++
	}
	return fmt.Sprintf("%c%d", rune(col), row), nil
}

func (this *WorkTime) Get(projectName string, year, month int) (*MonthlyWorkTime, error) {
	rows := this.sheet.Get(projectName, this.getSheetName(year, month), "A4", "K40")
	monthlyWorkTime, err := parseMonthlyWorkTime(year, month, rows)
	if err != nil {
		return nil, err
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

	this.sheet.Update(projectName, this.getSheetName(date.Year, date.Month), addr,
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

	this.sheet.Update(projectName, this.getSheetName(date.Year, date.Month), addr,
		fmt.Sprintf("%2d:%02d", time.Hour, time.Minute))

	return nil
}
