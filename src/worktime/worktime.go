package worktime

import (
	"fmt"
	"log"

	"work-time-logging/spreadsheet"
)

type WorkTime struct {
	sheet *spreadsheet.Spreadsheet
}

func New(sheet *spreadsheet.Spreadsheet) *WorkTime {
	return &WorkTime{sheet: sheet}
}

func (this *WorkTime) Get(target string, year, month int) {
	sheetName := fmt.Sprintf("%04d%02d", year, month)
	rows := this.sheet.Get(target, sheetName, "A4", "K40")
	fmt.Printf("%v", rows)
	records, err := ParseWorkTimeRecords(year, month, rows)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", records)
}
