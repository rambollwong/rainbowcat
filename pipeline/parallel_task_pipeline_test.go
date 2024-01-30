package pipeline

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// MockTaskProvider is a mock implementation of TaskProvider for testing purposes.
type MockTaskProvider struct{}

// Task is a mock implementation of the Task() method that returns a dummy task function.
func (m *MockTaskProvider) Task() Task {
	return func(input any) (output any, err error) {
		// Dummy task that simply appends "processed" to the input string
		return fmt.Sprintf("%s processed", input.(string)), nil
	}
}

func TestRunParallelTaskPipeline(t *testing.T) {
	// Define the number of pipelines and the maximum concurrent quantities for each pipeline
	pipelineCount := uint8(3)
	maxConcurrentQuantities := []uint8{3, 2, 3}

	// Create mock task providers for each pipeline
	taskProviders := []TaskProvider{
		&MockTaskProvider{},
		&MockTaskProvider{},
		&MockTaskProvider{},
	}

	// Run the parallel task pipeline
	ptp, err := RunParallelTaskPipeline(pipelineCount, maxConcurrentQuantities, taskProviders...)
	if err != nil {
		t.Errorf("Failed to run parallel task pipeline: %s", err)
	}

	// Push jobs into the pipeline
	jobs := []string{"job1", "job2", "job3", "job4", "job5"}
	for _, job := range jobs {
		ptp.PushJob(job)
	}

	// Close the pipeline
	defer ptp.Close()

	outputC := ptp.OutputC()
	for _, job := range jobs {
		require.Equal(t, fmt.Sprintf("%s processed processed processed", job), (<-outputC).(string))
	}
}
