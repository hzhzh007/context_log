package context_log

import (
	"bytes"
	"fmt"
	"github.com/rs/xid"
	"time"
)

type ContextLog struct {
	buf   *bytes.Buffer
	Uuid  string
	sTime time.Time
	tTime time.Time //临时统计用的时间
}

func NewContextLog(msg string) *ContextLog {
	sc := new(ContextLog)
	sc.buf = bytes.NewBufferString(msg)
	sc.Uuid = xid.New().String()
	sc.sTime = time.Now()
	//sc.tTime
	return sc
}

//just for compatibility
func NewContext(msg string) *ContextLog {
	return NewContextLog(msg)
}

func (sc *ContextLog) StartTimer() {
	sc.tTime = time.Now()
}

func (sc *ContextLog) StopTimer(key string) {
	duration := time.Now().Sub(sc.tTime)
	sc.buf.WriteString(fmt.Sprintf(" %s=%v", key, duration))
}

func (sc *ContextLog) AddNotes(key string, val string) {
	sc.buf.WriteString(fmt.Sprintf(" %s=%s", key, val))
}
func (sc *ContextLog) Flush() {
	duration := time.Now().Sub(sc.sTime)
	Log.Info(fmt.Sprintf("%s=%s cost=%v %s ", "Uuid", sc.Uuid, duration, sc.buf.String()))
}
func (sc *ContextLog) Debug(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Debugf(s, args...)
}
func (sc *ContextLog) Info(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Infof(s, args...)
}
func (sc *ContextLog) Notice(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Noticef(s, args...)
}
func (sc *ContextLog) Warning(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Warningf(s, args...)
}
func (sc *ContextLog) Error(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Errorf(s, args...)
}
func (sc *ContextLog) Critical(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Criticalf(s, args...)
}
