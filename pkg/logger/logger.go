package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLog() {
	Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
}
