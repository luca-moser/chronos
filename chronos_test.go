package chronos

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestDailySchedule(t *testing.T) {
	now := time.Now()
	plan := NewDailySchedulingPlan([]DayTime{
		DayTime{now.Hour(), now.Minute(), now.Second() + 4},
		DayTime{now.Hour(), now.Minute(), now.Second() + 2},
		DayTime{now.Hour(), now.Minute(), now.Second() + 1},
	})

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(6) * time.Second)
	if executionCount != 3 {
		t.Fatalf("expected execution count to be 3 but was %d", executionCount)
	}
}

func TestWeeklySchedule(t *testing.T) {
	now := time.Now()
	plan := NewWeeklySchedulingPlan([]Weekday{
		Weekday{day: now.Weekday(), at: DayTime{now.Hour(), now.Minute(), now.Second() + 4}},
		Weekday{day: now.Weekday(), at: DayTime{now.Hour(), now.Minute(), now.Second() + 2}},
		Weekday{day: now.Weekday(), at: DayTime{now.Hour(), now.Minute(), now.Second() + 1}},
	})

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(6) * time.Second)
	if executionCount != 3 {
		t.Fatalf("expected execution count to be 3 but was %d", executionCount)
	}
}

func TestMonthlySchedule(t *testing.T) {
	now := time.Now()
	plan := NewMonthlySchedulingPlan([]MonthDay{
		MonthDay{day: uint(now.Day()), at: DayTime{now.Hour(), now.Minute(), now.Second() + 4}},
		MonthDay{day: uint(now.Day()), at: DayTime{now.Hour(), now.Minute(), now.Second() + 2}},
		MonthDay{day: uint(now.Day()), at: DayTime{now.Hour(), now.Minute(), now.Second() + 1}},
	})

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(6) * time.Second)
	if executionCount != 3 {
		t.Fatalf("expected execution count to be 3 but was %d", executionCount)
	}
}
