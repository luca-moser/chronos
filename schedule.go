package khronos

import "time"

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
	days []MonthDay
	// weekdays (used for weekly scheduling)
	// - mondays at 9:45, 15:15 (multiple times on the given weekday)
	weekdays []Weekday
	// daily (used for daily scheduling)
	// - everyday at 14:00, 19:00
	dailyTimes []time

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
		case <-abort:
			break
		}
		// next execution time is reached, execute task
		exe<-struct{}{}
	}
}

// computes the duration until the next execution should happen by the given plan
func( ts *TaskSchedule) nextExecutionIn() time.Duration {
	now := time.Now()
	switch ts.plan {
	case INTERVAL_ONCE_IN:
		return ts.onceAfterDuration
	case INTERVAL_ONCE_DATE:
		return time.Until(ts.onceAtDate)
	case INTERVAL_EVERY_DAY:
		_ = now.Day()

		// what is the next day in our schedule?

	case INTERVAL_EVERY_WEEK:
	case INTERVAL_EVERY_MONTH:
	}
	return 0
}

// a day in a month and a specific time on that day
type MonthDay struct {
	day uint
	at  time.Time
}

// a weekday and a specific time on that day
type Weekday struct {
	day time.Weekday
	at  time.Time
}

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

func NewMonthlySchedulingPlan(days []MonthDay) TaskSchedule {
	taskSchedule := TaskSchedule{}
	taskSchedule.days = days
	taskSchedule.plan = INTERVAL_EVERY_MONTH
	return taskSchedule
}

func NewWeeklySchedulingPlan(weekdays []Weekday) TaskSchedule {
	taskSchedule := TaskSchedule{}
	taskSchedule.weekdays = weekdays
	taskSchedule.plan = INTERVAL_EVERY_WEEK
	return taskSchedule
}

func NewDailySchedulingPlan(times []time) TaskSchedule {
	taskSchedule := TaskSchedule{}
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
