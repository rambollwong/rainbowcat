package pipeline

import (
	"errors"
)

// Task defines the function signature of a task, which takes an input and returns an output and a boolean.
// If the returned boolean is false, the task will be terminated and the job will be ignored.
type Task func(input any) (output any, ok bool)

// TaskProvider interface defines a Task() method that returns a Task function.
type TaskProvider interface {
	Task() Task
}

// GenericTaskProvider is a function type that takes an input of type I and returns an output of type O.
type GenericTaskProvider[I, O any] func(input I) (output O, ok bool)

// Task method converts a GenericTaskProvider to a TaskProvider.
func (g GenericTaskProvider[I, O]) Task() Task {
	return func(input any) (output any, ok bool) {
		return g(input.(I))
	}
}

// Job struct represents a job to be executed in the pipeline.
// It contains an input, output, a flag indicating if the job is successful, and a channel to signal job completion.
type Job struct {
	Input     any
	Output    any
	Ok        bool
	FinishedC chan struct{}

	tp *taskPipeline
}

// do method submit the job to the pipeline for execution.
func (j *Job) do() {
	select {
	case <-j.tp.ptp.closeC:
		return
	case j.tp.jobC <- j:
		go j.run()
	}
}

// run method executes the task associated with the job and sends the output and error to the appropriate channels.
func (j *Job) run() {
	j.Output, j.Ok = j.tp.jobTask(j.Input)
	select {
	case <-j.tp.ptp.closeC:
	case j.FinishedC <- struct{}{}:
	}
}

// taskPipeline struct represents a single pipeline in the parallel task pipeline.
// It contains the pipeline index, a channel for receiving jobs, the task function for the pipeline,
// and a reference to the parent ParallelTaskPipeline.
type taskPipeline struct {
	index   uint8
	jobC    chan *Job
	jobTask Task

	ptp *ParallelTaskPipeline
}

// loop method continuously listens for incoming jobs and executes them.
// It also handles forwarding jobs to the next pipeline in the sequence.
func (tp *taskPipeline) loop() {
	for {
		select {
		case job := <-tp.jobC:
			select {
			case <-tp.ptp.closeC:
				return
			case <-job.FinishedC:
				if tp.ptp.pipelineCount == tp.index+1 {
					if !tp.ptp.noOutput {
						tp.ptp.outputC <- job.Output
					}
					continue
				}
				if !job.Ok {
					continue
				}
				job.Input = job.Output
				job.Output = nil
				job.FinishedC = make(chan struct{})
				job.tp = tp.ptp.pipelines[tp.index+1]
				job.do()
			}
		case <-tp.ptp.closeC:
			return
		}
	}
}

// ParallelTaskPipeline struct represents the entire parallel task pipeline. It contains the count of pipelines,
// an array of pipeline instances, and a channel for closing the pipeline.
type ParallelTaskPipeline struct {
	pipelineCount uint8
	pipelines     []*taskPipeline

	noOutput bool
	outputC  chan any
	closeC   chan struct{}
}

// RunParallelTaskPipeline function initializes and starts the parallel task pipeline.
// It takes the count of pipelines, an array of maximum concurrent quantities for each pipeline,
// and an array of task providers for each pipeline. It returns the initialized ParallelTaskPipeline instance.
// Each pipeline is responsible for executing the same logic that can be executed in parallel.
// The output of the task of the pipeline is used as the input of the task of the next pipeline.
// The output of the task of the last pipeline will be pushed to OutputC or ignored.
func RunParallelTaskPipeline(
	pipelineCount uint8,
	maxConcurrentQuantities []uint8,
	pipelineTaskProviders ...TaskProvider,
) (*ParallelTaskPipeline, error) {
	if pipelineCount == 0 {
		return nil, errors.New("invalid pipeline count")
	}
	if len(maxConcurrentQuantities) != int(pipelineCount) {
		return nil, errors.New("invalid max concurrent quantities")
	}
	if len(pipelineTaskProviders) != int(pipelineCount) {
		return nil, errors.New("invalid pipeline task providers")
	}
	p := &ParallelTaskPipeline{
		pipelineCount: pipelineCount,
		pipelines:     make([]*taskPipeline, pipelineCount),
		noOutput:      false,
		outputC:       make(chan any),
		closeC:        make(chan struct{}),
	}
	for i := uint8(0); i < pipelineCount; i++ {
		tp := &taskPipeline{
			index:   i,
			jobC:    make(chan *Job, maxConcurrentQuantities[i]),
			jobTask: pipelineTaskProviders[i].Task(),
			ptp:     p,
		}
		p.pipelines[i] = tp
		go tp.loop()
	}
	return p, nil
}

// Close method closes the pipeline and stops further execution of jobs.
func (p *ParallelTaskPipeline) Close() {
	close(p.closeC)
}

// PushJob method pushes a job into the pipeline by submitting it to the first pipeline in the sequence.
func (p *ParallelTaskPipeline) PushJob(input any) {
	firstTP := p.pipelines[0]
	job := &Job{
		Input:     input,
		Output:    nil,
		Ok:        false,
		FinishedC: make(chan struct{}),
		tp:        firstTP,
	}
	job.do()
}

// NoOutput sets a flag to indicate that the pipeline should not produce any output.
func (p *ParallelTaskPipeline) NoOutput() *ParallelTaskPipeline {
	p.noOutput = true
	return p
}

// OutputC returns a channel to receive the output from the pipeline.
// If the pipeline is configured to produce no output, it returns nil.
// Otherwise, it returns the outputC channel used to send the output.
func (p *ParallelTaskPipeline) OutputC() <-chan any {
	if p.noOutput {
		return nil
	}
	return p.outputC
}
