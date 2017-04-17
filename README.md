### Chronos - scheduled tasks for golang

Chronos is a library for scheduling tasks under
a given scheduling plan. It supports scheduling of monthly, weekly, daily
tasks or at a specific date or after a given duration. It is similar to the cron service on linux.

## Usage


Daily:
```go
// each day on 4:30, 10:45, 12:30
plan := NewDailySchedulingPlan([]DayTime{
   DayTime{4, 20, 0}, 
   DayTime{10, 45, 0}, 
   DayTime{12, 30},
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
      Weekday{day: time.Monday, at: DayTime{12, 0, 0},
      Weekday{day: time.Friday, at: DayTime{9, 0, 0},
      Weekday{day: time.Friday, at: DayTime{15, 0, 0},
      Weekday{day: time.Sunday, at: DayTime{19, 0, 0},
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
      MonthDay{day: 12, at: DayTime{16, 30, 0},
      MonthDay{day: 20, at: DayTime{18, 0, 0},
      MonthDay{day: 20, at: DayTime{21, 0, 0},      
      MonthDay{day: 31, at: DayTime{9, 0, 0},
})

task := NewScheduledTask(func() {
      ...
}, plan)
defer task.Stop()

// start the task
task.Start()
```