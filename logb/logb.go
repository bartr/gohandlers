package logb

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

// TODO - should we replace with Gorilla logger?
// TODO - add Apache log file support

var reqLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

// SetupAppLog - setup log multi writer
func SetupAppLog(logPath string, logPrefix string, logSuffix string) error {

	logFile := BuildFullLogName(logPath, logPrefix, logSuffix)

	// prepend date and time to log entries
	log.SetFlags(log.Ldate | log.Ltime)

	// open the log file
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	// setup a multiwriter to log to file and stdout
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	return nil
}

// BuildFullLogName - build the full log file name
// app services sets the WEBSITE_ROLE_INSTANCE_ID environment variable
//   since we're writing to the CIFS share, we need to differentiate log file names
//   in case there are multiple instances running
func BuildFullLogName(logPath string, logPrefix string, logSuffix string) string {
	if !strings.HasSuffix(logPath, "/") {
		logPath += "/"
	}

	fileName := logPath + logPrefix

	// use instance ID to differentiate log files between instances in App Services
	if iid := os.Getenv("WEBSITE_ROLE_INSTANCE_ID"); iid != "" {
		fileName += "_" + strings.TrimSpace(iid)
	}

	return fileName + logSuffix
}

// SetLogFile - initialize the log file and add multi writer
func SetLogFile(logPath string, logPrefix string, logSuffix string) error {
	logFile := BuildFullLogName(logPath, logPrefix, logSuffix)

	if logFile == "" {
		return errors.New("ERROR: logbpath cannot be blank")
	}

	// open the logfile
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	// setup the multi writer
	wrt := io.MultiWriter(os.Stdout, f)
	reqLog.SetOutput(wrt)

	return nil
}
