package main

import (
	"bytes"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
	"sync"
)

var LOG_LEVEL = 5
var LOG_LEVEL_LOCK sync.Mutex

const (
	DEBUG_LOG   = int(5)
	INFO_LOG    = int(4)
	WARNING_LOG = int(3)
	ERROR_LOG   = int(2)
	FATAL_LOG   = int(1)
	MONITOR_LOG = int(0)
)

var LOG_LEVEL_VALUE = map[string]int{
	"DEBUG":   5,
	"INFO":    4,
	"WARNING": 3,
	"ERROR":   2,
	"FATAL":   1,
	"MONITOR": 0,
}

var std *log.Logger

var Facility = syslog.LOG_LOCAL5

func SetFlags(flag int, prefix string) *log.Logger {
	std = log.New(os.Stderr, prefix, log.LstdFlags)
	std.SetFlags(flag)

	return std
}

// func NewSyslog(flag int, prefix string) *log.Logger {
//  var err error
//  std, err = syslog.NewLogger(syslog.LOG_DEBUG|Facility, flag)
//  if err != nil {
//      fmt.Println("NewSyslog error %s", err.Error())
//      return nil
//  }

//  std.SetPrefix(prefix)
//  return std
// }

func SetLogLevelString(level string) error {
	if LOG_LEVEL_VALUE[strings.ToUpper(level)] < FATAL_LOG {
		return fmt.Errorf("SetLogLevelString Error: Unsupport level:%s", level)
	}
	LOG_LEVEL_LOCK.Lock()
	LOG_LEVEL = LOG_LEVEL_VALUE[level]
	LOG_LEVEL_LOCK.Unlock()

	return nil
}

func SetLogLevelInt(level int) error {
	if level < FATAL_LOG || level > DEBUG_LOG {
		return fmt.Errorf("SetLogLevelString Error: Unsupport level:%d", level)
	}
	LOG_LEVEL_LOCK.Lock()
	LOG_LEVEL = level
	LOG_LEVEL_LOCK.Unlock()

	return nil
}

func LogPrintf(level int, format string, v ...interface{}) {

	var buffer bytes.Buffer
	if level <= LOG_LEVEL {
		switch level {
		case DEBUG_LOG:
			buffer.WriteString("[DEBUG]:")
		case INFO_LOG:
			buffer.WriteString("[INFO]:")
		case WARNING_LOG:
			buffer.WriteString("[WARNING]:")
		case ERROR_LOG:
			buffer.WriteString("[ERROR]:")
		case FATAL_LOG:
			buffer.WriteString("[FATAL]:")
		case MONITOR_LOG:
			buffer.WriteString("[MONITOR]:")
		}
		buffer.WriteString(format)
		std.Output(3, fmt.Sprintf(buffer.String(), v...))
	}

}

func LogDebug(format string, v ...interface{}) {
	LogPrintf(DEBUG_LOG, format, v...)
}

func LogInfo(format string, v ...interface{}) {
	LogPrintf(INFO_LOG, format, v...)
}

func LogWarning(format string, v ...interface{}) {
	LogPrintf(WARNING_LOG, format, v...)
}

func LogError(format string, v ...interface{}) {
	LogPrintf(ERROR_LOG, format, v...)
}
func LogFatal(format string, v ...interface{}) {
	LogPrintf(FATAL_LOG, format, v...)
}

func LogMonitor(format string, v ...interface{}) {
	LogPrintf(MONITOR_LOG, format, v...)
}
