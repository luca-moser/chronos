package chronos

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestOnceAtDateSchedule(t *testing.T) {
	now := time.Now()
	plan := NewOnceAtDatePlan(now.Add(time.Duration(3) * time.Second))

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(5) * time.Second)
	if executionCount != 1 {
		t.Fatalf("expected execution count to be 1 but was %d", executionCount)
	}
}

func TestOnceAfterDuration(t *testing.T) {
	plan := NewOnceAfterDuration(time.Duration(4) * time.Second)

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(6) * time.Second)
	if executionCount != 1 {
		t.Fatalf("expected execution count to be 1 but was %d", executionCount)
	}
}

func TestDailySchedule(t *testing.T) {
	now := time.Now()
	plan := NewDailySchedulingPlan([]DayTime{
		{now.Hour(), now.Minute(), now.Second() + 4},
		{now.Hour(), now.Minute(), now.Second() + 2},
		{now.Hour(), now.Minute(), now.Second() + 1},
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
		{day: now.Weekday(), at: DayTime{now.Hour(), now.Minute(), now.Second() + 4}},
		{day: now.Weekday(), at: DayTime{now.Hour(), now.Minute(), now.Second() + 2}},
		{day: now.Weekday(), at: DayTime{now.Hour(), now.Minute(), now.Second() + 1}},
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
		{day: uint(now.Day()), at: DayTime{now.Hour(), now.Minute(), now.Second() + 4}},
		{day: uint(now.Day()), at: DayTime{now.Hour(), now.Minute(), now.Second() + 2}},
		{day: uint(now.Day()), at: DayTime{now.Hour(), now.Minute(), now.Second() + 1}},
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
