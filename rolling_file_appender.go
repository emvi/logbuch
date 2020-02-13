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

type NameSchema interface {
	Name() string
}

type RollingFileAppender struct {
	Files    int
	FileSize int
	FileName NameSchema
	FileDir  string

	buffer          []byte
	maxBufferSize   int
	currentFile     *os.File
	currentFileSize int
	fileNames       []string
	m               sync.Mutex
}

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

func (appender *RollingFileAppender) Flush() error {
	appender.m.Lock()
	defer appender.m.Unlock()
	return appender.flush()
}

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