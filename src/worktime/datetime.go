package worktime

import (
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

type Time struct {
	Hour, Minute int
}

func (this *Time) RoundTime() *Time {
	diff := 10 - this.Minute % 10
	m := this.Minute + diff
	h := this.Hour
	if m >= 60 {
		h++
		m = 0
	}
	return &Time{Hour: h, Minute: m}
}
