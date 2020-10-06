package bitflyergo

import (
	"log"
)

// Logger is logger.
var Logger *log.Logger

func logln(v ...interface{}) {
	if Logger == nil {
		log.Println(v...)
		return
	}
	Logger.Println(v...)
}

func logf(format string, v ...interface{}) {
	if Logger == nil {
		log.Printf(format, v...)
		return
	}
	Logger.Printf(format, v...)
}
