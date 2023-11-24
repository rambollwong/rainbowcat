package filewriter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// RollingPeriod defines the enumeration for file rolling periods
type RollingPeriod string

const (
	RollingPeriodYear   RollingPeriod = "YEAR"
	RollingPeriodMonth  RollingPeriod = "MONTH"
	RollingPeriodDay    RollingPeriod = "DAY"
	RollingPeriodHour   RollingPeriod = "HOUR"
	RollingPeriodMinute RollingPeriod = "MINUTE"
	RollingPeriodSecond RollingPeriod = "SECOND"
)

var (
	TimeFormatYear   = "2006"
	TimeFormatMonth  = "200601"
	TimeFormatDay    = "20060102"
	TimeFormatHour   = "20060102_15"
	TimeFormatMinute = "20060102_15_04"
	TimeFormatSecond = "20060102_15_04_05"
)

// TimeRollingFileWriter is a time-based rolling file writer
type TimeRollingFileWriter struct {
	mu              sync.Mutex
	nextCheckTime   time.Time
	deleteCheckTime time.Time
	file            *os.File

	basePath       string
	baseFilePrefix string
	baseFileExt    string
	maxBackups     int
	rollPeriod     RollingPeriod
}

// NewTimeRollingFileWriter creates a new instance of TimeRollingFileWriter.
//
//	params:
//		- basePath: defines the path to save the files.
//		- baseFileName: defines the base name of the files. When file rotating occurs,
//			files may be renamed according to a specific format.
//		- maxBackups: defines the maximum number of file backups to keep.
//			If there is no limit, set the value to a negative value.
//		- rollPeriod: specify the time rolling period.
func NewTimeRollingFileWriter(
	basePath, baseFileName string,
	maxBackups int,
	rollPeriod RollingPeriod,
) (*TimeRollingFileWriter, error) {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, err
	}
	w := &TimeRollingFileWriter{}
	if maxBackups < 0 {
		maxBackups = 0
	}
	w.basePath = basePath
	w.maxBackups = maxBackups
	w.baseFileExt = filepath.Ext(baseFileName)
	w.baseFilePrefix = strings.TrimSuffix(baseFileName, w.baseFileExt)
	switch rollPeriod {
	case RollingPeriodYear, RollingPeriodMonth, RollingPeriodDay,
		RollingPeriodHour, RollingPeriodMinute, RollingPeriodSecond:
		w.rollPeriod = rollPeriod
	default:
		return nil, errors.New("unsupported roll period")
	}
	if err := w.tryRotate(); err != nil {
		return nil, err
	}
	return w, nil
}

// Close closes the file writer
func (w *TimeRollingFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}

// Write writes data to the file
func (w *TimeRollingFileWriter) Write(bz []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.tryRotate(); err != nil {
		return 0, err
	}
	return w.file.Write(bz)
}

// tryRotate attempts to perform file rotation
func (w *TimeRollingFileWriter) tryRotate() error {
	var (
		nextCheckTime   time.Time
		deleteCheckTime time.Time
		now             = time.Now()
		timeFormat      string
	)

	if time.Now().Before(w.nextCheckTime) {
		return nil
	}

	if w.file != nil {
		_ = w.file.Close()
	}

	switch w.rollPeriod {
	case RollingPeriodYear:
		nextCheckTime = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
		deleteCheckTime = time.Date(nextCheckTime.Year()-w.maxBackups, 1, 1, 0, 0, 0, 0, now.Location())
		timeFormat = TimeFormatYear

	case RollingPeriodMonth:
		nextCheckTime = time.Date(
			now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location(),
		).AddDate(0, 1, 0)
		deleteCheckTime = nextCheckTime.AddDate(0, -w.maxBackups, 0)
		timeFormat = TimeFormatMonth

	case RollingPeriodDay:
		nextCheckTime = time.Date(
			now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
		).AddDate(0, 0, 1)
		deleteCheckTime = nextCheckTime.AddDate(0, 0, -w.maxBackups)
		timeFormat = TimeFormatDay

	case RollingPeriodHour:
		nextCheckTime = time.Date(
			now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location(),
		).Add(time.Hour)
		deleteCheckTime = nextCheckTime.Add(-time.Duration(w.maxBackups) * time.Hour)
		timeFormat = TimeFormatHour

	case RollingPeriodMinute:
		nextCheckTime = time.Date(
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location(),
		).Add(time.Minute)
		deleteCheckTime = nextCheckTime.Add(-time.Duration(w.maxBackups) * time.Minute)
		timeFormat = TimeFormatMinute

	case RollingPeriodSecond:
		nextCheckTime = time.Date(
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location(),
		).Add(time.Second)
		deleteCheckTime = nextCheckTime.Add(-time.Duration(w.maxBackups) * time.Second)
		timeFormat = TimeFormatSecond

	default:
		return errors.New("unsupported roll period")
	}

	// Open the new file
	fileName := fmt.Sprintf("%s.%s%s", w.baseFilePrefix, now.Format(timeFormat), w.baseFileExt)
	file, err := os.OpenFile(filepath.Join(w.basePath, fileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	w.file = file

	// Set the next check time and delete check time
	w.nextCheckTime = nextCheckTime
	w.deleteCheckTime = deleteCheckTime

	if w.maxBackups >= 0 {
		// Try to delete old files
		go w.tryDeleteOldFiles()
	}

	return nil
}

// tryDeleteOldFiles tries to delete old files based on the delete check time
func (w *TimeRollingFileWriter) tryDeleteOldFiles() {
	files, err := filepath.Glob(filepath.Join(w.basePath, "*"+w.baseFileExt))
	if err != nil {
		fmt.Println("error while globbing files:", err)
		return
	}
	fileCount := len(files)
	if fileCount <= w.maxBackups {
		return
	}
	sort.Slice(files, func(i, j int) bool {
		indexTimeI, err := w.getFileIndexTime(files[i])
		if err != nil {
			return false
		}
		indexTimeJ, err := w.getFileIndexTime(files[j])
		if err != nil {
			return false
		}
		return indexTimeI.After(indexTimeJ)
	})
	for _, file := range files {
		fileTime, err := w.getFileIndexTime(file)
		if err != nil {
			fmt.Println("error while getting file index time: " + err.Error())
			fileCount--
			continue
		}
		// Check if the file is older than the delete check time
		if fileTime.Before(w.deleteCheckTime) {
			err = os.Remove(file)
			if err != nil {
				fmt.Println("failed to remove old file:", err)
			}
			fileCount--
		}
		if fileCount <= w.maxBackups {
			return
		}
	}
}

// getFileIndexTime extracts the index time from the given file name.
// It parses the file name based on the rolling period and returns the corresponding time value.
func (w *TimeRollingFileWriter) getFileIndexTime(file string) (time.Time, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return time.Time{}, err
	}
	fileName := fileInfo.Name()
	fileName = strings.TrimSuffix(fileName, w.baseFileExt)
	fileDate := strings.TrimPrefix(fileName, w.baseFilePrefix+".")
	var fileTime time.Time
	switch w.rollPeriod {
	case RollingPeriodYear:
		fileTime, err = time.ParseInLocation(TimeFormatYear, fileDate, w.deleteCheckTime.Location())
	case RollingPeriodMonth:
		fileTime, err = time.ParseInLocation(TimeFormatMonth, fileDate, w.deleteCheckTime.Location())
	case RollingPeriodDay:
		fileTime, err = time.ParseInLocation(TimeFormatDay, fileDate, w.deleteCheckTime.Location())
	case RollingPeriodHour:
		fileTime, err = time.ParseInLocation(TimeFormatHour, fileDate, w.deleteCheckTime.Location())
	case RollingPeriodMinute:
		fileTime, err = time.ParseInLocation(TimeFormatMinute, fileDate, w.deleteCheckTime.Location())
	case RollingPeriodSecond:
		fileTime, err = time.ParseInLocation(TimeFormatSecond, fileDate, w.deleteCheckTime.Location())
	default:
		panic("bug found! unexpected roll period value found")
	}
	return fileTime, err
}
