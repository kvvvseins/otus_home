package hw06pipelineexecution

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

	out := recursiveStaging(in, done, stages)

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
			default:
				continue
			}
		}
	}()

	return result
}

func recursiveStaging(in In, done In, stages []Stage) Out {
	lenStages := len(stages)

	for i := range stages {
		select {
		case <-done:
			return in
		default:
			return recursiveStaging(stages[i](in), done, stages[i+1:lenStages])
		}
	}

	return in
}
