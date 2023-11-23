package filewriter

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTimeRollingFileWriter_Write(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "filewriter_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a TimeRollingFileWriter instance
	writer, err := NewTimeRollingFileWriter(tempDir, "test.log", 3, RollingPeriodDay)
	if err != nil {
		t.Fatalf("Failed to create TimeRollingFileWriter: %v", err)
	}
	defer writer.Close()

	// Write data
	data := []byte("Hello, World!")
	_, err = writer.Write(data)
	if err != nil {
		t.Fatalf("Failed to write data: %v", err)
	}

	// Ensure the file is created
	files, err := filepath.Glob(filepath.Join(tempDir, "test.*.log"))
	if err != nil {
		t.Fatalf("Failed to glob files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	// Read the file content and verify
	fileContent, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(fileContent) != string(data) {
		t.Fatalf("File content mismatch, expected '%s', got '%s'", string(data), string(fileContent))
	}
}

func TestTimeRollingFileWriter_Rotate(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "filewriter_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a TimeRollingFileWriter instance
	writer, err := NewTimeRollingFileWriter(tempDir, "test.log", 5, RollingPeriodSecond)
	if err != nil {
		t.Fatalf("Failed to create TimeRollingFileWriter: %v", err)
	}
	defer writer.Close()

	for i := 0; i < 10; i++ {
		// Write data
		data := []byte("Hello, World!")
		_, err = writer.Write(data)
		if err != nil {
			t.Fatalf("Failed to write data: %v", err)
		}

		// Wait for 1 Second to trigger file rotation
		time.Sleep(1 * time.Second)
	}

	// Ensure the old file is deleted
	files, err := filepath.Glob(filepath.Join(tempDir, "test.*.log"))
	if err != nil {
		t.Fatalf("Failed to glob files: %v", err)
	}
	if len(files) != 5 {
		t.Fatalf("Expected 5 file, got %d", len(files))
	}
}
