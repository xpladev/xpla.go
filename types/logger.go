package types

import "github.com/logrusorgru/aurora"

type Logger interface {
	Info(interface{}, ...LMsg)
	Warn(interface{}, ...LMsg)
	Err(interface{}, ...LMsg) error
}

type LMsg struct {
	Key   string
	Value string
}

func (l LMsg) LogKV() string {
	return bb(l.Key+"=") + l.Value
}

func LogMsg(key string, value string) LMsg {
	return LMsg{
		Key:   key,
		Value: value,
	}
}

func bb(str string) string {
	return aurora.BrightBlue(str).String()
}
