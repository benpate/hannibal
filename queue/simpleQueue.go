package queue

// Simple is a perfectly workable in-memory queue that uses a buffered channel
// to store tasks and multiple goroutines as task runners.
//
// Some day it may make sense to use a ring buffer, but today is not that day
// https://bravenewgeek.com/so-you-wanna-go-fast/
// https://github.com/Workiva/go-datastructures/blob/master/queue/ring.go
type SimpleQueue struct {
	tasks   chan Task
	done    chan struct{}
	options []RetryOption
}

// NewSimpleQueu returns a fully initialized SimpleQueue, which will process
// tasks in the background until the Close() method is called.
func NewSimpleQueue(workers int, maxLength int, options ...RetryOption) Queue {

	// Create the queue
	result := SimpleQueue{
		tasks:   make(chan Task, maxLength),
		done:    make(chan struct{}),
		options: options,
	}

	// Spin up workers to receive tasks
	for counter := 0; counter < workers; counter++ {
		go result.worker()
	}

	// Done
	return &result
}

// Push adds a task to the queue
func (q *SimpleQueue) Push(task Task) {
	go func() {
		q.tasks <- task
	}()
}

// worker processes tasks as they are received
func (q *SimpleQueue) worker() {

	for {
		select {

		case <-q.done:
			return

		case task := <-q.tasks:
			if err := task.Run(); err != nil {
				go Retry(task.Run, q.options...)
			}
		}
	}
}

// Close shuts down all of the queue workers.
func (q *SimpleQueue) Close() {
	close(q.done)
}
