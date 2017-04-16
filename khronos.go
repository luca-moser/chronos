package khronos

type ScheduledFunc func(...interface{}) (interface{}, error)

// creates a new task which executes under the given schedule
func NewScheduledTask(action ScheduledFunc, schedule TaskSchedule) ScheduledTask {
	task := ScheduledTask{}
	task.action = action
	task.schedule = schedule
	task.executeSignal = make(chan struct{})
	task.abortSignal = make(chan struct{})
	return task
}

// a task which's action executes by the given schedule.
type ScheduledTask struct {
	action        ScheduledFunc
	schedule      TaskSchedule
	executeSignal chan struct{}
	abortSignal   chan struct{}
}

// starts the scheduling (non-blocking)
func (st *ScheduledTask) Start() {
	// should use the New... functions to create a schedule
	if st.schedule.plan == 0 {
		panic("schedule has no plan defined")
	}

	// initialize plan
	scheduleAbortSignal := make(chan struct{}, 1)
	go st.schedule.Init(st.executeSignal, scheduleAbortSignal)

	// listen for signals to process the task
	go func() {
		for {
			select {
			case <-st.executeSignal:
				st.action()
			case <-st.abortSignal:
				scheduleAbortSignal<-struct{}{}
				break
			}
		}
	}()
}
