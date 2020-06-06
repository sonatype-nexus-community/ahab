package logger

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
)

const DefaultLogFilename = "ahab.combined.log"

// DefaultLogFile can be overridden to use a different file name for upstream consumers
var DefaultLogFile = DefaultLogFilename

func GetLogger(level logrus.Level) (*logrus.Logger, error) {
	logger := logrus.New()

	logger.Level = level
	logger.Formatter = &logrus.JSONFormatter{}
	location, err := LogFileLocation()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(location, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	logger.Out = file

	return logger, nil
}

// LogFileLocation will return the location on disk of the log file
func LogFileLocation() (result string, err error) {
	result, _ = os.UserHomeDir()
	err = os.MkdirAll(path.Join(result, types.OssIndexDirName), os.ModePerm)
	if err != nil {
		return
	}
	result = path.Join(result, types.OssIndexDirName, DefaultLogFile)
	return
}

func PrintErrorAndLogLocation(err error) {
	fmt.Println()
	fmt.Println("Uh oh, an error occurred, if this persists try rerunning with -v, -vv, or -vvv to get more information in the logs")
	fmt.Printf("Error: %v\n", err)
	location, _ := LogFileLocation()
	fmt.Printf("Check log file at %s for more information\n", location)
}
