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
	dailyTimes         []time.Time
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

	for {
		nextExecutionSignal := time.NewTimer(ts.nextExecutionIn())
		select {
		case <-nextExecutionSignal.C:
			if ts.plan == INTERVAL_ONCE_IN || ts.plan == INTERVAL_ONCE_DATE {
				break
			}
		case <-abort:
			break
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
		next := time.Date(now.Year(), now.Month(), now.Day(), nextTime.Hour(), nextTime.Minute(), nextTime.Second(), 0, time.UTC)
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
		next = time.Date(now.Year(), now.Month(), next.Day(), at.Hour(), at.Minute(), at.Second(), 0, time.UTC)
		return time.Until(next)

	case INTERVAL_EVERY_MONTH:
		todayNum := now.Day()
		var next time.Time
		nextDay := ts.nextDay()
		at := nextDay.at
		nextDayNum := nextDay.day

		// TODO: check leap year and 30/31 case

		if int(nextDayNum) < todayNum {
			// advance one month
			next = now.AddDate(0, 1, 0)
			next = time.Date(next.Year(), next.Month(), int(nextDayNum), at.Hour(), at.Minute(), at.Second(), 0, time.UTC)
		} else {
			next = time.Date(now.Year(), now.Month(), int(nextDayNum), at.Hour(), at.Minute(), at.Second(), 0, time.UTC)
		}
		return time.Until(next)
	}
	return 0
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
func (ts *TaskSchedule) nextDailyTime() time.Time {
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

func NewDailySchedulingPlan(times []time.Time) TaskSchedule {
	if len(times) == 0 {
		panic("times slice is empty")
	}
	taskSchedule := TaskSchedule{}
	sort.Sort(TimesSorted(times))
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
