package chronos

import "time"

// creates a new month day (day must be >0 and <=31)
func NewMonthDay(day uint, at DayTime) MonthDay {
	if day > 31 {
		panic("month day can't be greater than 31")
	}
	// we use the actual real day number
	if day == 0 {
		panic("month day can't be 0")
	}
	return MonthDay{day, at}
}

// a day in a month and a specific time on that day
type MonthDay struct {
	Day uint    `json:"day"`
	At DayTime `json:"at"`
}

type MonthDaysSorted []MonthDay

func (mds MonthDaysSorted) Len() int      { return len(mds) }
func (mds MonthDaysSorted) Swap(i, j int) { mds[i], mds[j] = mds[j], mds[i] }
func (mds MonthDaysSorted) Less(i, j int) bool {
	this := mds[i]
	next := mds[j]

	if this.Day < next.Day {
		return true
	}

	if this.Day > next.Day {
		return false
	}

	return this.At.AsTime().Before(next.At.AsTime())
}

// a weekday and a specific time on that day
type Weekday struct {
	Day time.Weekday `json:"weekday"`
	At DayTime      `json:"at"`
}

type WeekdaysSorted []Weekday

func (wds WeekdaysSorted) Len() int      { return len(wds) }
func (wds WeekdaysSorted) Swap(i, j int) { wds[i], wds[j] = wds[j], wds[i] }
func (wds WeekdaysSorted) Less(i, j int) bool {
	this := wds[i]
	next := wds[j]

	if this.Day < next.Day {
		return true
	}

	if this.Day > next.Day {
		return false
	}

	return this.At.AsTime().Before(next.At.AsTime())
}

// represents a 24 hour day
type DayTime struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}

func (dt *DayTime) AsTime() time.Time {
	return time.Date(0, 0, 0, dt.Hour, dt.Minute, dt.Second, 0, time.Local)
}

type DayTimesSorted []DayTime

func (dts DayTimesSorted) Len() int      { return len(dts) }
func (dts DayTimesSorted) Swap(i, j int) { dts[i], dts[j] = dts[j], dts[i] }
func (dts DayTimesSorted) Less(i, j int) bool {
	return dts[i].AsTime().Before(dts[j].AsTime())
}
