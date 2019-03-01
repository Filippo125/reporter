package main

import (
	"fmt"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

type PlainTextFormatter struct {
}

func (f *PlainTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	message := entry.Message
	time := entry.Time
	level := entry.Level
	return []byte(fmt.Sprintf("%s [%s] %s\n", time.Format("2006-01-02T15:04:05"), strings.ToUpper(level.String()), message)), nil
}
