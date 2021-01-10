package worktime

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type WorkTimeRecord struct {
	Date    Date
	Periods []Period
}

type Date struct {
	Year, Month, Day int
}

func (this *Date) IsLastDayOfMonth() bool {
	d := time.Date(this.Year, time.Month(this.Month), this.Day, 0, 0, 0, 0, time.UTC)
	return d.Month() != d.AddDate(0, 0, 1).Month()
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
	jst, err := time.LoadLocation("Asia/Tokyo")
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

func ParseWorkTimeRecords(year, month int, rows [][]interface{}) ([]WorkTimeRecord, error) {
	var records []WorkTimeRecord
	for _, row := range rows {
		var record []string
		for _, v := range row {
			if s, ok := v.(string); ok {
				record = append(record, s)
			} else {
				return nil, fmt.Errorf("")
			}
		}

		if len(record) < 8 {
			return nil, fmt.Errorf("Invalid record: %v", record)
		}

		date, err := parseDate(year, month, string(record[0]))
		if err != nil {
			return nil, err
		}

		var periods []Period
		for _, i := range []int{2, 4, 6} {
			p, err := parsePeriod(date, string(record[i]), string(record[i+1]))
			if err != nil {
				return nil, err
			}
			periods = append(periods, *p)
		}

		records = append(records, WorkTimeRecord{Date: *date, Periods: periods})

		if date.IsLastDayOfMonth() {
			break
		}
	}
	return records, nil
}
