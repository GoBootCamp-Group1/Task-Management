package log

import (
	"io"
	"log"
	"os"
)

var (
	InfoLog    *log.Logger
	TraceLog   *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
)

func init() {
	logFile, err := os.OpenFile("appLogs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	InfoLog = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	TraceLog = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
