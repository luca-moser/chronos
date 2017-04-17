package chronos

import "time"

// creates a new month day (day must be >0 and <=31)
func NewMonthDay(day uint, at time.Time) MonthDay {
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
	day uint
	at  time.Time
}

type MonthDaysSorted []MonthDay

func (mds MonthDaysSorted) Len() int      { return len(mds) }
func (mds MonthDaysSorted) Swap(i, j int) { mds[i], mds[j] = mds[j], mds[i] }
func (mds MonthDaysSorted) Less(i, j int) bool {
	this := mds[i]
	next := mds[j]

	if this.day < next.day {
		return true
	}

	if this.day > next.day {
		return false
	}

	return this.at.Before(next.at)
}

// a weekday and a specific time on that day
type Weekday struct {
	day time.Weekday
	at  time.Time
}

type WeekdaysSorted []Weekday

func (wds WeekdaysSorted) Len() int      { return len(wds) }
func (wds WeekdaysSorted) Swap(i, j int) { wds[i], wds[j] = wds[j], wds[i] }
func (wds WeekdaysSorted) Less(i, j int) bool {
	this := wds[i]
	next := wds[j]

	if this.day < next.day {
		return true
	}

	if this.day > next.day {
		return false
	}

	return this.at.Before(next.at)
}

type TimesSorted []time.Time

func (ts TimesSorted) Len() int      { return len(ts) }
func (ts TimesSorted) Swap(i, j int) { ts[i], ts[j] = ts[j], ts[i] }
func (ts TimesSorted) Less(i, j int) bool {
	return ts[i].Before(ts[j])
}