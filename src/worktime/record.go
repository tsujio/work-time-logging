package worktime

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

type MonthlyWorkTime struct {
	Year, Month int
	Records     []WorkTimeRecord
}

func (this *MonthlyWorkTime) GetDuration() time.Duration {
	var sum time.Duration
	for _, r := range this.Records {
		sum += r.GetDuration()
	}
	return sum
}

type WorkTimeRecord struct {
	Date          *Date
	Periods       []Period
	TravelExpense *TravelExpense
}

func (this *WorkTimeRecord) GetDuration() time.Duration {
	var sum time.Duration
	for _, p := range this.Periods {
		sum += p.GetDuration()
	}
	return sum
}

type Period struct {
	Start time.Time
	End   time.Time
}

func (this *Period) IsEmpty() bool {
	return this.Start.Equal(time.Time{}) && this.IsEndEmpty()
}

func (this *Period) IsEndEmpty() bool {
	return this.End.Equal(time.Time{})
}

func (this *Period) GetDuration() time.Duration {
	if this.IsEmpty() || this.IsEndEmpty() {
		return 0
	} else {
		return this.End.Sub(this.Start)
	}
}

type TravelExpense struct {
	Expense int
	Note    string
}

func parseTravelExpense(expense, note string) (*TravelExpense, error) {
	e, err := strconv.Atoi(strings.TrimSpace(expense))
	if err != nil {
		return nil, err
	}
	return &TravelExpense{Expense: e, Note: note}, nil
}

func parseDate(year, month int, value string) (*Date, error) {
	slice := strings.Split(value, "/")
	if len(slice) != 2 {
		return nil, fmt.Errorf("Invalid date: %s", value)
	}
	m, err := strconv.Atoi(slice[0])
	if err != nil {
		return nil, fmt.Errorf("Invalid month: %s", slice[0])
	}
	if m != month {
		return nil, fmt.Errorf("Month mismatch: %d != %d", m, month)
	}
	d, err := strconv.Atoi(slice[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid date: %s", slice[1])
	}
	return &Date{Year: year, Month: m, Day: d}, nil
}

func parseTime(date *Date, value string) (time.Time, error) {
	slice := strings.Split(value, ":")
	if len(slice) != 2 {
		return time.Time{}, fmt.Errorf("Invalid time: %s", value)
	}
	h, err := strconv.Atoi(slice[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid hour: %s", slice[0])
	}
	m, err := strconv.Atoi(slice[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid minute: %s", slice[1])
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year, time.Month(date.Month), date.Day, h, m, 0, 0, jst), nil
}

func parsePeriod(date *Date, start, end string) (*Period, error) {
	if start == "" {
		if end == "" {
			return &Period{Start: time.Time{}, End: time.Time{}}, nil
		} else {
			return nil, fmt.Errorf("%v's start is empty but end is present", date)
		}
	} else {
		s, err := parseTime(date, start)
		if err != nil {
			return nil, err
		}

		if end == "" {
			return &Period{Start: s, End: time.Time{}}, nil
		} else {
			e, err := parseTime(date, end)
			if err != nil {
				return nil, err
			}

			if s.Equal(e) || s.After(e) {
				e = e.AddDate(0, 0, 1)
			}
			if s.After(e) {
				log.Fatalf("Invalid period: %v - %v", s, e)
			}

			return &Period{Start: s, End: e}, nil
		}
	}
}

func parseDuration(d string) (time.Duration, error) {
	slice := strings.Split(d, ":")
	if len(slice) != 2 {
		return 0, xerrors.Errorf("Invalid duration: %s", d)
	}
	h, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, xerrors.Errorf("Invalid hours: %s", slice[0])
	}
	m, err := strconv.Atoi(slice[1])
	if err != nil {
		return 0, xerrors.Errorf("Invalid minutes: %s", slice[1])
	}
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}

func parseWorkTimeRecord(year, month int, row []string) (*WorkTimeRecord, error) {
	if len(row) < 9 {
		return nil, xerrors.Errorf("Invalid row: %v", row)
	}

	date, err := parseDate(year, month, string(row[0]))
	if err != nil {
		return nil, xerrors.Errorf("Unable to parse date: %w", err)
	}

	var periods []Period
	for _, i := range []int{2, 4, 6} {
		p, err := parsePeriod(date, string(row[i]), string(row[i+1]))
		if err != nil {
			return nil, xerrors.Errorf("Unable to parse period: %w", err)
		}
		periods = append(periods, *p)
	}

	var travelExpense *TravelExpense
	if len(row) > 10 {
		travelExpense, err = parseTravelExpense(row[10], row[9])
		if err != nil {
			return nil, xerrors.Errorf("Unable to parse travel expense: %w", err)
		}
	} else {
		travelExpense = nil
	}

	record := &WorkTimeRecord{
		Date:          date,
		Periods:       periods,
		TravelExpense: travelExpense,
	}

	// Validate duration
	sumActual := record.GetDuration()
	sumGiven, err := parseDuration(row[8])
	if err != nil {
		return nil, xerrors.Errorf("Unable to parse duration: %w", err)
	}
	if sumActual != sumGiven {
		return nil, xerrors.Errorf("Duration mismatch: date=%v, given=%v, actual=%v", date, sumGiven, sumActual)
	}

	return record, nil
}

func parseMonthlyWorkTime(year, month int, rows [][]interface{}) (*MonthlyWorkTime, error) {
	var records []WorkTimeRecord
	for i, rawRow := range rows {
		var row []string
		for _, v := range rawRow {
			if s, ok := v.(string); ok {
				row = append(row, s)
			} else {
				return nil, xerrors.Errorf("Invalid row: %v", rawRow)
			}
		}

		record, err := parseWorkTimeRecord(year, month, row)
		if err != nil {
			return nil, xerrors.Errorf("Unable to parse work time record: %w", err)
		}

		records = append(records, *record)

		if record.Date.IsLastDayOfMonth() {
			// Validate duration
			var sumActual time.Duration
			for _, r := range records {
				sumActual += r.GetDuration()
			}
			s, ok := rows[i+1][8].(string)
			if !ok {
				return nil, xerrors.Errorf("Invalid row: %v", rows[i+1])
			}
			sumGiven, err := parseDuration(s)
			if err != nil {
				return nil, xerrors.Errorf("Unable to parse total duration: %w", err)
			}
			if sumActual != sumGiven {
				return nil, xerrors.Errorf("Total duration mismatch: given=%v, actual=%v", sumActual, sumGiven)
			}

			break
		}
	}
	return &MonthlyWorkTime{Year: year, Month: month, Records: records}, nil
}
