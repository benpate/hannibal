package pub

type QueueTask struct {
	run func() error
}

func NewQueueTask(run func() error) QueueTask {
	return QueueTask{
		run: run,
	}
}

func (task QueueTask) Run() error {
	return task.run()
}
