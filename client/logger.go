package client

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	"github.com/logrusorgru/aurora"
)

const (
	// Default verbose is 0.
	VerboseDefault     = 0
	VerboseDetails     = 1
	VerboseImplication = 2
)

var _ types.Logger = &logger{}

type logger struct {
	verbose int
}

func newLogger(verbose int) types.Logger {
	return &logger{
		verbose: verbose,
	}
}

func (l *logger) Info(log interface{}, msgs ...types.LMsg) {
	var print string

	switch {
	case l.verbose == VerboseDetails:
		print = strings.Join([]string{g("INF"), genLogs(log, msgs...), logTime()}, " ")

	case l.verbose == VerboseImplication:
		print = genLogs(log, msgs...)

	default:
		print = strings.Join([]string{g("INF"), genLogs(log, msgs...)}, " ")
	}

	fmt.Println(print)
}

func (l *logger) Warn(log interface{}, msgs ...types.LMsg) {
	var print string

	switch {
	case l.verbose == VerboseDetails:
		print = strings.Join([]string{y("WAR"), genLogs(log, msgs...), logTime()}, " ")

	case l.verbose == VerboseImplication:
		print = genLogs(log, msgs...)

	default:
		print = strings.Join([]string{y("WAR"), genLogs(log, msgs...)}, " ")
	}

	fmt.Println(print)
}

func (l *logger) Err(log interface{}, msgs ...types.LMsg) error {
	switch {
	case l.verbose == VerboseDetails:
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		stackTraceMsg := types.LMsg{
			Key:   "occured",
			Value: string(buf[:stackSize]),
		}

		logLevelMsg := types.LMsg{
			Key:   "verbose",
			Value: util.FromIntToString(VerboseDetails),
		}

		print := strings.Join([]string{
			r("ERR"),
			genLogs(log, msgs...),
			logTime(),
			logLevelMsg.LogKV(),
			stackTraceMsg.LogKV(),
		}, " ")

		return fmt.Errorf(print)

	case l.verbose == VerboseImplication:
		return fmt.Errorf(genLogs(log, msgs...))

	default:
		_, file, line, _ := runtime.Caller(1)
		callerMsg := types.LMsg{
			Key:   "occured",
			Value: file + ":" + util.FromIntToString(line),
		}

		print := strings.Join([]string{
			r("ERR"),
			genLogs(log, msgs...),
			logTime(),
			callerMsg.LogKV(),
		}, " ")
		return fmt.Errorf(print)
	}
}

func genLogs(log interface{}, msgs ...types.LMsg) string {
	print := interfaceLog(log, "")
	if len(msgs) != 0 {
		for _, msg := range msgs {
			print = print + " " + msg.LogKV()
		}
	}

	return print
}

func logTime() string {
	timeMsg := types.LMsg{
		Key:   "time",
		Value: time.Now().Format("2006-01-02 15:04:05"),
	}
	return timeMsg.LogKV()
}

func g(str string) string {
	return aurora.Green(str).String()
}

func r(str string) string {
	return aurora.Red(str).String()
}

func y(str string) string {
	return aurora.Yellow(str).String()
}

func interfaceLog(value interface{}, defaultValue string) string {
	s := fmt.Sprintf("%v", value)
	str := strings.TrimSpace(s)
	if str == "" {
		return defaultValue
	} else {
		return str
	}
}
