package queue

// Queue is an interface that wraps the Push method.
type Queue interface {

	// Push adds a new task to the Queue.  The queue is expected to
	// process the task asynchronously, and return immediately.
	Push(task Task)
}

// Task is an interface that wraps the Run method.
// Tasks are fed into the queue and processed asynchronously.
type Task interface {
	Run() error
}
