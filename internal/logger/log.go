package logger

import (
	"fmt"
	"log"
	"time"
)

var (
	Info     *log.Logger
	Error    *log.Logger
	Critical *log.Logger
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	if bytes[len(bytes)-1] == '\n' {
		bytes[len(bytes)-1] = '"'
	}
	return fmt.Print("[" + time.Now().UTC().Format("02/Jan/2006:15:04:05 -0700") + "]" + string(bytes) + "\n")
}

func init() {
	Info = log.New(new(logWriter), " [INFO] \"", log.Lmsgprefix)
	Error = log.New(new(logWriter), " [ERROR] \"", log.Lmsgprefix)
	Critical = log.New(new(logWriter), " [CRITICAL] \"", log.Lmsgprefix)
}
