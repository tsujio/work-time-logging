package worktime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	Year, Month, Day int
}

func (this *Date) Equal(date *Date) bool {
	return *this == *date
}

func (this *Date) IsLastDayOfMonth() bool {
	d := time.Date(this.Year, time.Month(this.Month), this.Day, 0, 0, 0, 0, time.UTC)
	return d.Month() != d.AddDate(0, 0, 1).Month()
}

func (this *Date) GetWeekday() time.Weekday {
	d := time.Date(this.Year, time.Month(this.Month), this.Day, 0, 0, 0, 0, time.UTC)
	return d.Weekday()
}

type Time struct {
	Hour, Minute int
}

func ParseHHMM(hhmm string) (*Time, error) {
	slice := strings.Split(hhmm, ":")
	if len(slice) != 2 {
		return nil, fmt.Errorf("Invalid time: %s", hhmm)
	}
	h, err := strconv.Atoi(slice[0])
	if err != nil {
		return nil, fmt.Errorf("Invalid hour: %s", slice[0])
	}
	m, err := strconv.Atoi(slice[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid minute: %s", slice[1])
	}
	return &Time{h, m}, nil
}

func (this *Time) RoundTime() *Time {
	diff := (10 - this.Minute%10) % 10
	m := this.Minute + diff
	h := this.Hour
	if m >= 60 {
		h++
		m = 0
	}
	return &Time{Hour: h, Minute: m}
}
