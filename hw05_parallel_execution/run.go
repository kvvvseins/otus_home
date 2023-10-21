package hw05parallelexecution

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
)

const MaxGoroutines = "1000"

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrEmptyTasks          = errors.New("на выполнение пришло 0 задач")
	ErrEmptyNGoroutines    = errors.New("нужно указать хотя бы 1 поток")
	ErrMaxNGoroutines      = errors.New("максимально можно указать только " + MaxGoroutines + " потоков")
)

type Task func() error

type Process struct {
	cntUncompletedTasks int64
	maxUncompletedTasks int64
	wg                  *sync.WaitGroup
	mu                  *sync.Mutex
}

func (p *Process) incrementError() {
	p.mu.Lock()
	defer p.mu.Unlock()

	atomic.AddInt64(&p.cntUncompletedTasks, 1)
}

func (p *Process) isMaxErrors() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.maxUncompletedTasks == 0 || p.cntUncompletedTasks >= p.maxUncompletedTasks
}

type Workers struct {
	tasksChan chan Task
	max       int64
	cnt       int64
	mu        *sync.Mutex
}

func (w *Workers) increment() {
	w.mu.Lock()
	defer w.mu.Unlock()

	atomic.AddInt64(&w.cnt, 1)
}

func (w *Workers) decrement() {
	w.mu.Lock()
	defer w.mu.Unlock()

	atomic.AddInt64(&w.cnt, -1)
}

func (w *Workers) isMax() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.cnt >= w.max
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if err := checkRequest(tasks, n); err != nil {
		return err
	}

	cntTasks := int64(len(tasks))

	result := make(chan error, cntTasks)
	defer close(result)

	process := &Process{mu: &sync.Mutex{}, maxUncompletedTasks: int64(m), wg: &sync.WaitGroup{}}
	defer process.wg.Wait()

	go func() {
		workers := &Workers{
			tasksChan: make(chan Task, n),
			max:       int64(n),
			mu:        &sync.Mutex{},
		}
		defer close(workers.tasksChan)

		for _, task := range tasks {
			if ok := produceWorkers(task, workers, process, result); !ok {
				break
			}
		}
	}()

	return consume(result, process, cntTasks)
}

func produceWorkers(
	task Task,
	workers *Workers,
	process *Process,
	result chan<- error,
) bool {
	for {
		select {
		case taskWorker := <-workers.tasksChan:
			go func() {
				err := taskWorker()
				if err != nil {
					process.incrementError()
				}

				result <- err
				workers.decrement()
				process.wg.Done()
			}()

			return true
		default:
			if process.isMaxErrors() {
				return false
			}

			if !workers.isMax() {
				workers.tasksChan <- task

				workers.increment()
				process.wg.Add(1)
			}
		}
	}
}

func consume(
	result <-chan error,
	process *Process,
	cntTasks int64,
) error {
	completedTasks := int64(0)

	for {
		select {
		case <-result:
			atomic.AddInt64(&completedTasks, 1)
		default:
			if process.isMaxErrors() {
				return ErrErrorsLimitExceeded
			}

			if completedTasks == cntTasks {
				return nil
			}
		}
	}
}

func checkRequest(tasks []Task, n int) error {
	if len(tasks) == 0 {
		return ErrEmptyTasks
	}

	if n == 0 {
		return ErrEmptyNGoroutines
	}

	if maxGoroutines, _ := strconv.Atoi(MaxGoroutines); n >= maxGoroutines {
		return ErrMaxNGoroutines
	}

	return nil
}
