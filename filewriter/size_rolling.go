package filewriter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// SizeRollingFileWriter is a file writer with rolling based on file size.
type SizeRollingFileWriter struct {
	mu          sync.Mutex
	file        *os.File
	currentSize int64

	basePath       string
	baseFilePrefix string
	baseFileExt    string
	maxBackups     int
	fileSizeLimit  int64
}

// NewSizeRollingFileWriter creates a new SizeRollingFileWriter instance with the given parameters.
//
//	params:
//		- basePath: defines the path to save the files.
//		- baseFileName: defines the base name of the files. When file rotating occurs,
//			files may be renamed according to a specific format.
//		- maxBackups: defines the maximum number of file backups to keep.
//			If there is no limit, set the value to a negative value.
//	 	- fileSizeLimit: defines the maximum size of each file in bytes.
//	 		When maxBackups is not a negative value, if the current file size reaches the upper limit,
//	 		rotation will be triggered.
func NewSizeRollingFileWriter(
	basePath, baseFileName string,
	maxBackups int,
	fileSizeLimit int64,
) (*SizeRollingFileWriter, error) {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, err
	}
	w := &SizeRollingFileWriter{}
	if maxBackups < 0 {
		maxBackups = 0
	}
	w.basePath = basePath
	w.maxBackups = maxBackups
	w.baseFileExt = filepath.Ext(baseFileName)
	w.baseFilePrefix = strings.TrimSuffix(baseFileName, w.baseFileExt)
	w.fileSizeLimit = fileSizeLimit
	if err := w.openFile(); err != nil {
		return nil, err
	}
	if err := w.tryRotate(0); err != nil {
		return nil, err
	}
	return w, nil
}

// Close closes the file writer.
func (w *SizeRollingFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}

// Write writes data to the file.
func (w *SizeRollingFileWriter) Write(bz []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.tryRotate(int64(len(bz))); err != nil {
		return 0, err
	}
	n, err = w.file.Write(bz)
	if err != nil {
		return n, err
	}
	w.currentSize += int64(n)
	return
}

// tryRotate checks if the current file size exceeds the limit and performs log rotation if necessary.
func (w *SizeRollingFileWriter) tryRotate(bytesLength int64) error {
	if w.currentSize == 0 || w.currentSize+bytesLength <= w.fileSizeLimit {
		return nil
	}

	files, err := filepath.Glob(filepath.Join(w.basePath, "*"+w.baseFileExt))
	if err != nil {
		return errors.New("error while globbing files: " + err.Error())
	}
	fileCount := len(files)
	sort.Slice(files, func(i, j int) bool {
		return w.getFileIndex(files[i]) > w.getFileIndex(files[j])
	})
	for _, file := range files {
		fileIndexInt := w.getFileIndex(file)
		if fileIndexInt == 0 {
			fileCount--
			continue
		}
		if fileCount > w.maxBackups && fileIndexInt > w.maxBackups-1 {
			err = os.Remove(file)
			if err != nil {
				return errors.New("error while removing file: " + err.Error())
			}
			fileCount--
			continue
		}
		newFileName := fmt.Sprintf("%s.%d%s", w.baseFilePrefix, fileIndexInt+1, w.baseFileExt)
		err = os.Rename(file, filepath.Join(w.basePath, newFileName))
		if err != nil {
			return err
		}
	}

	if w.file != nil {
		_ = w.file.Close()
		newFileName := fmt.Sprintf("%s.1%s", w.baseFilePrefix, w.baseFileExt)
		err = os.Rename(
			w.file.Name(),
			filepath.Join(w.basePath, newFileName),
		)
		if err != nil {
			return err
		}
	}

	return w.openFile()
}

// openFile opens the current log file for writing.
func (w *SizeRollingFileWriter) openFile() error {
	file, err := os.OpenFile(filepath.Join(w.basePath, w.baseFilePrefix+w.baseFileExt), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	w.file = file

	info, err := file.Stat()
	if err != nil {
		return err
	}
	w.currentSize = info.Size()
	return nil
}

// getFileIndex extracts the index number from the file name.
func (w *SizeRollingFileWriter) getFileIndex(file string) int {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return 0
	}
	fileName := fileInfo.Name()
	fileName = strings.TrimSuffix(fileName, w.baseFileExt)
	fileIndex := strings.TrimPrefix(fileName, w.baseFilePrefix+".")
	fileIndexInt, err := strconv.Atoi(fileIndex)
	if err != nil {
		return 0
	}
	return fileIndexInt
}
