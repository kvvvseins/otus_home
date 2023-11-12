package hw06pipelineexecution

import (
	"os"
	"os/signal"
	"syscall"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		result := make(Bi)
		defer close(result)

		return result
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	out := recursiveStaging(in, done, sigChan, stages)

	result := make(Bi)

	go func() {
		defer close(result)

		for {
			select {
			case v, ok := <-out:
				if !ok {
					return
				}

				result <- v
			case <-done:
				return
			case <-sigChan:
				return
			default:
				continue
			}
		}
	}()

	return result
}

func recursiveStaging(in In, done In, sigChan chan os.Signal, stages []Stage) Out {
	lenStages := len(stages)

	for i := range stages {
		pipelineIn := make(Bi)

		go func() {
			defer close(pipelineIn)

			for {
				select {
				case val, ok := <-in:
					if !ok {
						return
					}

					pipelineIn <- val
				case <-done:
					return
				case <-sigChan:
					return
				}
			}
		}()

		return recursiveStaging(stages[i](pipelineIn), done, sigChan, stages[i+1:lenStages])
	}

	return in
}
