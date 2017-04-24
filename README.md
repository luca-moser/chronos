### Chronos - scheduled tasks for golang [![Build Status](https://travis-ci.org/luca-moser/chronos.svg?branch=master)](https://travis-ci.org/luca-moser/chronos)

Chronos is a library for scheduling tasks under
a given scheduling plan. It supports scheduling of monthly, weekly, daily
tasks or at a specific date or after a given duration. It is similar to the cron service on linux.

## Usage

Notes:
* The task's `Start()` function is non blocking.
* The defined action is executed in a separate goroutine.

### Schedule

Daily:
```go
// each day on 4:30, 10:45, 12:30
plan := NewDailySchedulingPlan([]DayTime{
   {4, 20, 0}, 
   {10, 45, 0}, 
   {12, 30},
})

task := NewScheduledTask(func() {
      ...
}, plan)
defer task.Stop()

// start the task
task.Start()
```

Weekly:
```go
// on monday at 12, on friday twice at 9 am and 3 pm, once at sunday at 7 pm
plan := NewWeeklySchedulingPlan([]DayTime{
      {Day: time.Monday, At: DayTime{12, 0, 0},
      {Day: time.Friday, At: DayTime{9, 0, 0},
      {Day: time.Friday, At: DayTime{15, 0, 0},
      {Day: time.Sunday, At: DayTime{19, 0, 0},
})

task := NewScheduledTask(func() {
      ...
}, plan)
defer task.Stop()

// start the task
task.Start()
```

Monthly:
```go
// on the 12th at 4:30 pm, 20th at 6 pm and 9 pm, 31th on 9 am
// NOTE: the 31th auto. shrinks to 28, 29 (February with leap year) or 30
plan := NewMonthlySchedulingPlan([]MonthDay{
      {Day: 12, At: DayTime{16, 30, 0},
      {Day: 20, At: DayTime{18, 0, 0},
      {Day: 20, At: DayTime{21, 0, 0},
      {Day: 31, At: DayTime{9, 0, 0},
})

task := NewScheduledTask(func() {
      ...
}, plan)
defer task.Stop()

// start the task
task.Start()
```

### Once 

at a given date:
```go
// once on the 29th of November at 3 pm
plan := NewOnceAtDatePlan(time.Date(2017, time.November, 29, 15, 0, 0, 0, time.Local))

task := NewScheduledTask(func() {
      ...
}, plan)
defer task.Stop()

// start the task
task.Start()
```

after a given duration (essentially the same as doing time.After(d time.Duration)):
```go
// once after 5 seconds
plan := NewOnceAfterDuration(time.Duration(5) * time.Second)

task := NewScheduledTask(func() {
      ...
}, plan)
defer task.Stop()

// start the task
task.Start()
```