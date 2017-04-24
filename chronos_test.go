package chronos

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestOnceAtDateSchedule(t *testing.T) {
	t.Parallel()
	now := time.Now()
	plan := NewOnceAtDatePlan(now.Add(time.Duration(2) * time.Second))

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(3) * time.Second)
	if executionCount != 1 {
		t.Fatalf("expected execution count to be 1 but was %d", executionCount)
	}
}

func TestOnceAfterDuration(t *testing.T) {
	t.Parallel()
	plan := NewOnceAfterDuration(time.Duration(2) * time.Second)

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(3) * time.Second)
	if executionCount != 1 {
		t.Fatalf("expected execution count to be 1 but was %d", executionCount)
	}
}

func TestDailySchedule(t *testing.T) {
	t.Parallel()
	now := time.Now()
	plan := NewDailySchedulingPlan([]DayTime{
		{now.Hour(), now.Minute(), now.Second() + 3},
		{now.Hour(), now.Minute(), now.Second() + 2},
		{now.Hour(), now.Minute(), now.Second() + 1},
	})

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(4) * time.Second)
	if executionCount != 3 {
		t.Fatalf("expected execution count to be 3 but was %d", executionCount)
	}
}

func TestWeeklySchedule(t *testing.T) {
	t.Parallel()
	now := time.Now()
	plan := NewWeeklySchedulingPlan([]Weekday{
		{Day: now.Weekday(), At: DayTime{now.Hour(), now.Minute(), now.Second() + 3}},
		{Day: now.Weekday(), At: DayTime{now.Hour(), now.Minute(), now.Second() + 2}},
		{Day: now.Weekday(), At: DayTime{now.Hour(), now.Minute(), now.Second() + 1}},
	})

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(4) * time.Second)
	if executionCount != 3 {
		t.Fatalf("expected execution count to be 3 but was %d", executionCount)
	}
}

func TestMonthlySchedule(t *testing.T) {
	t.Parallel()
	now := time.Now()
	plan := NewMonthlySchedulingPlan([]MonthDay{
		{Day: uint(now.Day()), At: DayTime{now.Hour(), now.Minute(), now.Second() + 3}},
		{Day: uint(now.Day()), At: DayTime{now.Hour(), now.Minute(), now.Second() + 2}},
		{Day: uint(now.Day()), At: DayTime{now.Hour(), now.Minute(), now.Second() + 1}},
	})

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, plan)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(4) * time.Second)
	if executionCount != 3 {
		t.Fatalf("expected execution count to be 3 but was %d", executionCount)
	}
}

func TestFilePersistence(t *testing.T) {
	t.Parallel()
	plan := NewOnceAfterDuration(time.Duration(4) * time.Second)

	storage := NewFileStorage("./", "schedule.json")
	if err := storage.SaveSchedule(&plan); err != nil {
		t.Fatal(err)
	}

	defer func() {
		os.Remove(path.Join("./", "schedule.json"))
	}()

	planFromDisk, err := storage.LoadSchedule()
	if err != nil {
		t.Fatal(err)
	}

	executionCount := 0
	task := NewScheduledTask(func() {
		executionCount++
	}, *planFromDisk)
	defer task.Stop()
	task.Start()
	<-time.After(time.Duration(5) * time.Second)
	if executionCount != 1 {
		t.Fatalf("expected execution count to be 1 but was %d", executionCount)
	}
}
