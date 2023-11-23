package filewriter

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSizeRollingFileWriter_Write(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "filewriter_test")
	if err != nil {
		t.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a SizeRollingFileWriter instance
	filePath := filepath.Join(tempDir, "test.log")
	maxBackups := 3
	fileSizeLimit := int64(100)
	writer, err := NewSizeRollingFileWriter(tempDir, "test.log", maxBackups, fileSizeLimit)
	if err != nil {
		t.Fatal("Failed to create SizeRollingFileWriter:", err)
	}
	defer writer.Close()

	// Write data and verify file content
	data := []byte("Hello, World!")
	_, err = writer.Write(data)
	if err != nil {
		t.Fatal("Error writing to file:", err)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("Error reading file content:", err)
	}

	if !bytes.Equal(fileContent, data) {
		t.Error("File content does not match the written data")
	}

	// Write data exceeding the file size limit and verify file rolling
	backupFilePath := filepath.Join(tempDir, "test.1.log")

	// Write data exceeding the file size limit
	largeData := make([]byte, fileSizeLimit+1)
	_, err = writer.Write(largeData)
	if err != nil {
		t.Fatal("Error writing large data to file:", err)
	}

	// Verify if the backup file exists
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		t.Fatal("Backup file does not exist")
	}

	// Verify if the new original file contains the written data
	newFileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("Error reading new file content:", err)
	}

	if !bytes.Equal(newFileContent, largeData) {
		t.Fatal("New file content does not match the written large data")
	}

	for i := 0; i < maxBackups+1; i++ {
		_, _ = writer.Write(data)
		_, _ = writer.Write(largeData)
	}

	// Verify the number of backup files matches the expected value
	backupFiles, err := filepath.Glob(filepath.Join(tempDir, "*.log"))
	if err != nil {
		t.Fatal("Error globbing backup files:", err)
	}

	if len(backupFiles) != maxBackups+1 {
		t.Fatalf("Expected %d backup files, got %d", maxBackups, len(backupFiles))
	}
}
