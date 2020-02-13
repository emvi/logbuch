package logbuch

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

const (
	defaultFiles      = 10
	defaultFileSize   = 1024 * 1024 * 5 // 5 MB
	defaultBufferSize = 4096            // 4 KB
)

// NameSchema is an interface to generate log file names.
// If you implement this interface, make sure the Name() method returns unique file names.
type NameSchema interface {
	// Name returns the next file name used to store log data.
	Name() string
}

// RollingFileAppender is a manager for rolling log files.
// It needs to be closed using the Close() method.
type RollingFileAppender struct {
	// Files is the number of files used before rolling over.
	Files int

	// FileSize is the maximum size of a single log file.
	FileSize int

	// FileName is the naming schema used to create the next log file.
	FileName NameSchema

	// FileDir is the output directory for log files.
	FileDir string

	buffer          []byte
	maxBufferSize   int
	currentFile     *os.File
	currentFileSize int
	fileNames       []string
	m               sync.Mutex
}

// NewRollingFileAppender creates a new RollingFileAppender.
// If you pass values below or equal to 0 for files, size or bufferSize, default values will be used.
// The file output directory is created if required and can be left empty to use the current directory.
// The filename schema is required.
func NewRollingFileAppender(files, size, bufferSize int, dir string, filename NameSchema) (*RollingFileAppender, error) {
	if files <= 0 {
		files = defaultFiles
	}

	if size <= 0 {
		size = defaultFileSize
	}

	if bufferSize <= 0 {
		bufferSize = defaultBufferSize
	}

	if filename == nil {
		return nil, errors.New("filename schema must be specified")
	}

	if err := os.MkdirAll(dir, 0774); err != nil {
		return nil, err
	}

	appender := &RollingFileAppender{Files: files,
		FileSize:      size,
		FileName:      filename,
		FileDir:       dir,
		buffer:        make([]byte, 0, bufferSize),
		maxBufferSize: bufferSize,
		fileNames:     make([]string, 0, files)}

	if err := appender.nextFile(); err != nil {
		return nil, err
	}

	return appender, nil
}

// Write writes given data to the rolling log files.
// This might not happen immediately as the RollingFileAppender uses a buffer.
// If you want the data to be persisted, call Flush().
func (appender *RollingFileAppender) Write(p []byte) (n int, err error) {
	appender.m.Lock()
	defer appender.m.Unlock()
	appender.buffer = append(appender.buffer, p...)

	if len(appender.buffer) >= appender.maxBufferSize {
		if err := appender.flush(); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

// Flush writes all log data currently in buffer into the currently active log file.
func (appender *RollingFileAppender) Flush() error {
	appender.m.Lock()
	defer appender.m.Unlock()
	return appender.flush()
}

// Close flushes the log data and closes all open file handlers.
func (appender *RollingFileAppender) Close() error {
	appender.m.Lock()
	defer appender.m.Unlock()

	if err := appender.flush(); err != nil {
		return err
	}

	return appender.currentFile.Close()
}

func (appender *RollingFileAppender) flush() error {
	offset := 0

	for offset < len(appender.buffer) {
		if appender.currentFileSize >= appender.FileSize {
			if err := appender.nextFile(); err != nil {
				return err
			}
		}

		bytes := len(appender.buffer) - offset
		maxBytes := appender.FileSize - appender.currentFileSize

		if bytes > maxBytes {
			bytes = maxBytes
		}

		n, err := appender.currentFile.Write(appender.buffer[offset : offset+bytes])

		if err != nil {
			return err
		}

		appender.currentFileSize += n
		offset += n
	}

	appender.buffer = appender.buffer[:0]
	return nil
}

func (appender *RollingFileAppender) nextFile() error {
	if appender.currentFile != nil {
		if err := appender.currentFile.Close(); err != nil {
			return err
		}
	}

	path := filepath.Join(appender.FileDir, appender.FileName.Name())
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0664)

	if err != nil {
		return err
	}

	appender.currentFile = f
	appender.currentFileSize = 0
	return appender.updateFiles(path)
}

func (appender *RollingFileAppender) updateFiles(path string) error {
	appender.fileNames = append(appender.fileNames, path)

	if len(appender.fileNames) > appender.Files {
		n := len(appender.fileNames) - appender.Files
		filesToDelete := appender.fileNames[:n]
		appender.fileNames = appender.fileNames[n:]

		for _, file := range filesToDelete {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}

	return nil
}
