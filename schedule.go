package chronos

import (
	"sort"
	"time"
)

// intervals
type ScheduleInterval byte

const (
	// schedule once at the given date
	INTERVAL_ONCE_DATE ScheduleInterval = iota + 1
	// schedule once after the given duration
	INTERVAL_ONCE_IN
	// schedule every month at the given day
	INTERVAL_EVERY_MONTH
	// schedule every week at the given weekday
	INTERVAL_EVERY_WEEK
	// schedule every day at the given time
	INTERVAL_EVERY_DAY
)

/*
	represents a schedule which uses the interval constants
	to define the plans for a given task.
*/
type TaskSchedule struct {

	// specific days of a month, i.e. the 8th, 16th, 20th (used in monthly scheduling)
	// - 8th at 08:00, 12:00... (multiple times on the given days)
	days         []MonthDay
	nextDayIndex int

	// weekdays (used for weekly scheduling)
	// - mondays at 9:45, 15:15 (multiple times on the given weekday)
	weekdays         []Weekday
	nextWeekdayIndex int

	// daily (used for daily scheduling)
	// - everyday at 14:00, 19:00
	dailyTimes         []DayTime
	nextDailyTimeIndex int

	// once at specific date
	onceAtDate time.Time

	// once after duration
	onceAfterDuration time.Duration

	// how are we scheduled?
	plan ScheduleInterval

	// next execution time
	nextExecutionOn time.Time
}

func (ts *TaskSchedule) Init(exe chan<- struct{}, abort <-chan struct{}) {

exit:
	for {
		nextDurationToWait := ts.nextExecutionIn()
		nextExecutionSignal := time.NewTimer(nextDurationToWait)
		select {
		case <-nextExecutionSignal.C:
			if ts.plan == INTERVAL_ONCE_IN || ts.plan == INTERVAL_ONCE_DATE {
				exe <- struct{}{}
				break exit
			}
		case <-abort:
			break exit
		}
		// next execution time is reached, execute task
		exe <- struct{}{}
	}
}

// computes the duration until the next execution should happen by the given plan
func (ts *TaskSchedule) nextExecutionIn() time.Duration {
	now := time.Now()
	switch ts.plan {
	case INTERVAL_ONCE_IN:
		return ts.onceAfterDuration
	case INTERVAL_ONCE_DATE:
		return time.Until(ts.onceAtDate)

	case INTERVAL_EVERY_DAY:
		nextTime := ts.nextDailyTime()
		next := time.Date(now.Year(), now.Month(), now.Day(), nextTime.hour, nextTime.minute, nextTime.second, 0, time.Local)
		if next.Before(time.Now()) {
			// the next time is on the next day
			next = next.AddDate(0, 0, 1)
		}
		return time.Until(next)

	case INTERVAL_EVERY_WEEK:
		todayWeekday := now.Weekday()

		// next
		var next time.Time
		nextWeekday := ts.nextWeekday()
		at := nextWeekday.at
		nextWeekdayNum := nextWeekday.day

		if nextWeekdayNum < todayWeekday {
			// advance a week
			inWeekAdv := 7 - todayWeekday
			next = now.AddDate(0, 0, int(inWeekAdv+nextWeekdayNum))
		} else {
			// same week
			next = now.AddDate(0, 0, int(nextWeekdayNum-todayWeekday))
		}
		next = time.Date(now.Year(), now.Month(), next.Day(), at.hour, at.minute, at.second, 0, time.Local)
		if next.Before(time.Now()) {
			// the next time is in one week
			next = next.AddDate(0, 0, 7)
		}
		return time.Until(next)

	case INTERVAL_EVERY_MONTH:
		todayNum := now.Day()
		var next time.Time
		nextDay := ts.nextDay()
		at := nextDay.at
		nextDayNum := nextDay.day

		if int(nextDayNum) < todayNum {
			// advance one month
			next = addMonth(now)

			// auto. shrink to the correct day of the next month
			nextDay := shrinkDay(nextDayNum, next)
			next = time.Date(next.Year(), next.Month(), int(nextDay), at.hour, at.minute, at.second, 0, time.Local)
		} else {
			nextDay := shrinkDay(nextDayNum, now)
			next = time.Date(now.Year(), now.Month(), int(nextDay), at.hour, at.minute, at.second, 0, time.Local)
		}
		if next.Before(time.Now()) {
			// the next time is in one month
			nextMonthDate := addMonth(next)
			nextDay := shrinkDay(nextDayNum, nextMonthDate)
			next = time.Date(nextMonthDate.Year(), nextMonthDate.Month(), int(nextDay), at.hour, at.minute, at.second, 0, time.Local)
		}
		return time.Until(next)
	}
	return 0
}

// advances one month:
// unlike the built in time.AddDate function it actually advances one month
// instead of just adding 30 or whatever days to the given date.
func addMonth(date time.Time) time.Time {
	if date.Month() == 12 {
		return time.Date(date.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
	}
	return time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.Local)
}

// shrinks the given day correctly to the given month's last day
// with respect to leap years
func shrinkDay(day uint, date time.Time) uint {
	month := date.Month()
	year := date.Year()

	if month == time.February && (day == 30 || day == 31) {
		if isLeapYear(year) {
			return 29
		}
		return 28
	}

	if day == 31 {
		switch month {
		case time.April:
			fallthrough
		case time.June:
			fallthrough
		case time.September:
			fallthrough
		case time.November:
			return 30
		}
	}

	return day
}

func isLeapYear(year int) bool {
	return year%400 == 0 || year%4 == 0 && year%100 != 0
}

// returns the next day to check the duration for
func (ts *TaskSchedule) nextDay() MonthDay {
	day := ts.days[ts.nextDayIndex]
	if ts.nextDayIndex == len(ts.days)-1 {
		ts.nextDayIndex = 0 // reset
	} else {
		ts.nextDayIndex++
	}
	return day
}

// returns the next weekday to check the duration for
func (ts *TaskSchedule) nextWeekday() Weekday {
	weekday := ts.weekdays[ts.nextWeekdayIndex]
	if ts.nextWeekdayIndex == len(ts.weekdays)-1 {
		ts.nextWeekdayIndex = 0 // reset
	} else {
		ts.nextWeekdayIndex++
	}
	return weekday
}

// returns the next daily time to check the duration for
func (ts *TaskSchedule) nextDailyTime() DayTime {
	dailyTime := ts.dailyTimes[ts.nextDailyTimeIndex]
	if ts.nextDailyTimeIndex == len(ts.dailyTimes)-1 {
		ts.nextDailyTimeIndex = 0 // reset
	} else {
		ts.nextDailyTimeIndex++
	}
	return dailyTime
}

func NewMonthlySchedulingPlan(days []MonthDay) TaskSchedule {
	if len(days) == 0 {
		panic("days slice is empty")
	}
	taskSchedule := TaskSchedule{}
	sort.Sort(MonthDaysSorted(days))
	taskSchedule.days = days
	taskSchedule.plan = INTERVAL_EVERY_MONTH
	return taskSchedule
}

func NewWeeklySchedulingPlan(weekdays []Weekday) TaskSchedule {
	if len(weekdays) == 0 {
		panic("weekdays slice is empty")
	}
	taskSchedule := TaskSchedule{}
	sort.Sort(WeekdaysSorted(weekdays))
	taskSchedule.weekdays = weekdays
	taskSchedule.plan = INTERVAL_EVERY_WEEK
	return taskSchedule
}

func NewDailySchedulingPlan(times []DayTime) TaskSchedule {
	if len(times) == 0 {
		panic("times slice is empty")
	}
	taskSchedule := TaskSchedule{}
	sort.Sort(DayTimesSorted(times))
	taskSchedule.dailyTimes = times
	taskSchedule.plan = INTERVAL_EVERY_DAY
	return taskSchedule
}

func NewOnceAtDatePlan(date time.Time) TaskSchedule {
	taskSchedule := TaskSchedule{}
	taskSchedule.onceAtDate = date
	taskSchedule.plan = INTERVAL_ONCE_DATE
	return taskSchedule
}

func NewOnceAfterDuration(duration time.Duration) TaskSchedule {
	taskSchedule := TaskSchedule{}
	taskSchedule.onceAfterDuration = duration
	taskSchedule.plan = INTERVAL_ONCE_IN
	return taskSchedule
}
