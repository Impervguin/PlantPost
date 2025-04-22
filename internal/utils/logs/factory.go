package logs

import (
	"os"
	"time"
)

var _ LogFileFactory = &EveryDayFileFactory{}

type EveryDayFileFactory struct {
	currentFile *os.File
	lastDay     time.Time
	filePrefix  string
	fileSuffix  string
}

func NewEveryDayFileFactory(filePrefix string, fileSuffix string) *EveryDayFileFactory {
	return &EveryDayFileFactory{
		currentFile: nil,
		lastDay:     time.Time{},
		filePrefix:  filePrefix,
		fileSuffix:  fileSuffix,
	}
}

func (fa *EveryDayFileFactory) updFile() error {
	fname := time.Now().Format(time.DateOnly)
	fname = fname + fa.fileSuffix
	f, err := os.Create(fa.filePrefix + fname)
	if err != nil {
		return err
	}
	if fa.currentFile != nil {
		fa.currentFile.Close()
	}
	fa.currentFile = f

	return nil
}

func (f *EveryDayFileFactory) GetLogFile() *os.File {
	if time.Since(f.lastDay) > 24*time.Hour {
		f.updFile()
		f.lastDay = time.Now().UTC().Truncate(24 * time.Hour)
	}
	return f.currentFile
}

func (f *EveryDayFileFactory) Write(p []byte) (n int, err error) {
	return f.GetLogFile().Write(p)
}

func (f *EveryDayFileFactory) Sync() error {
	return f.GetLogFile().Sync()
}
