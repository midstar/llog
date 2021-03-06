// Package llog (Level Logger) extends the standard log package with:
//
//    * configurable log levels
//    * log file wrapping if configurable size exceeded
//
// llog is using the "standard" logger in the log package.
//
// Default level is LvlInfo and default output is stderr.
package llog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Level type is used for different debuggnig levels
type Level int

const (
	// LvlTrace debugging - logs of high frequency
	LvlTrace Level = 1
	// LvlDebug debugging - logs of medium frequency
	LvlDebug Level = 2
	// LvlInfo debugging - logs of low frequency
	LvlInfo Level = 3
	// LvlWarn something goes wrong but is not really an error
	LvlWarn Level = 4
	// LvlError an error has occurred
	LvlError Level = 5
	// LvlPanic a non-recoverable error has occurred
	LvlPanic Level = 6
)

// globLevelSet is the current level set. Default is LvlInfo.
var globLevelSet = LvlInfo

// globFileName is the name of file where logging output goes or nil if stderr
var globFileName string

// globFile is the file where logging output goes or nil if stderr
var globFile *os.File

// globMaxSizeKB max size of log file until wrap
var globMaxSizeKB int

// globMutex is a mutex to secure thread safety
var globMutex = &sync.Mutex{}

// init initialize llog.
func init() {
	// Log on format:
	// 2009/01/23 01:23:23 file.go:23: INFO - message
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// SetLevel sets lowest log priority that shall be written to the output.
func SetLevel(level Level) {
	globLevelSet = level
}

// SetFile logs to a file instead of stderr (default). If the file is more
// than maxSizeKB the old file will be backed up and a new log file
// will be written. If an error occurs stderr logging will be kept.
func SetFile(fileName string, maxSizeKB int) error {
	var err error
	globFileName = fileName
	globFile, err = os.OpenFile(globFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		globFile = nil
		return err
	}
	// Don't close the file with intention

	log.SetOutput(globFile)
	globMaxSizeKB = maxSizeKB
	return nil
}

// globCounter is counting to know when log wrap should be checked
var globCounter int

//wrapLogIfNeeded wraps the log if globMaxSizeKB has exceeded. To avoid
//file accesses for every log entry the actual file check is only done
//every 20th log write.
func wrapLogIfNeeded() {
	globMutex.Lock() // For thread safety
	defer globMutex.Unlock()

	if globFile == nil {
		// Not storing to a file
		return
	}

	globCounter++
	if globCounter >= 20 {
		globCounter = 0 // Reset counter
		globFile.Sync()
		info, _ := globFile.Stat()
		if (info.Size() / 1024) >= int64(globMaxSizeKB) {
			// Time to wrap
			log.SetOutput(os.Stderr) // Temporary log to stderr
			globFile.Close() // Close file
			backupFileName := globFileName + ".1"
			os.Remove(backupFileName)           // Remove backup if existing
			os.Rename(globFileName, backupFileName) // Make backup
			SetFile(globFileName, globMaxSizeKB)    // Start over on log
		}
	}
}

func loglevel(level Level, prefix string, format string, v ...interface{}) {
	if level >= globLevelSet {
		wrapLogIfNeeded()
		log.Output(3, fmt.Sprintf(prefix+format, v...))
	}
}

// Trace writes a log on trace level
func Trace(format string, v ...interface{}) {
	loglevel(LvlTrace, "TRACE - ", format, v...)
}

// Debug writes a log on debug level
func Debug(format string, v ...interface{}) {
	loglevel(LvlDebug, "DEBUG - ", format, v...)
}

// Info writes a log on info level
func Info(format string, v ...interface{}) {
	loglevel(LvlInfo, "INFO - ", format, v...)
}

// Warn writes a log on warn level
func Warn(format string, v ...interface{}) {
	loglevel(LvlWarn, "WARN - ", format, v...)
}

// Error writes a log on error level
func Error(format string, v ...interface{}) {
	loglevel(LvlError, "ERROR - ", format, v...)
}

// Panic writes a log on panic level, flush
// the log and calls panic()
func Panic(format string, v ...interface{}) {
	if LvlPanic >= globLevelSet {
		log.Output(2, fmt.Sprintf("PANIC - "+format, v...))
		panic(fmt.Sprintf(format, v...))
	}
}
