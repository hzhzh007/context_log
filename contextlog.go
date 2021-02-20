package context_log

import (
	"bytes"
	"fmt"
	"github.com/rs/xid"
	"time"
)

type ContextLog struct {
	buf        *bytes.Buffer
	Uuid       string
	sTime      time.Time
	tTime      time.Time //临时统计用的时间
	Flushed    bool
	SubContext []*ContextLog
}

func NewContext(msg string) *ContextLog {
	sc := new(ContextLog)
	sc.buf = bytes.NewBufferString(msg)
	sc.Uuid = xid.New().String()
	sc.sTime = time.Now()
	sc.Flushed = false
	sc.SubContext = make([]*ContextLog, 0)
	//sc.tTime
	return sc
}

func (sc *ContextLog) NewSubContext(msg string) *ContextLog {
	subSc := NewContext(msg)
	subSc.Uuid = sc.Uuid + "_" + subSc.Uuid
	sc.SubContext = append(sc.SubContext, subSc)
	return subSc
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
	sc.Flushed = true
	Log.Info(fmt.Sprintf("%s=%s cost=%v %s ", "Uuid", sc.Uuid, duration, sc.buf.String()))
	for _, subSc := range sc.SubContext {
		if !subSc.Flushed {
			subSc.Flush()
		}
	}
}
func (sc *ContextLog) Debug(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Debug(varArgsInsert(s, args...))
}
func (sc *ContextLog) Info(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Info(varArgsInsert(s, args...))
}
func (sc *ContextLog) Notice(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Notice(varArgsInsert(s, args...))
}
func (sc *ContextLog) Warning(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Warning(varArgsInsert(s, args...))
}
func (sc *ContextLog) Error(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Error(varArgsInsert(s, args...))
}

func (sc *ContextLog) Critical(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Critical(varArgsInsert(s, args...))
}
func varArgsInsert(s interface{}, args ...interface{}) []interface{} {
	t := make([]interface{}, 0, len(args)+1)
	t = append(t, s)
	t = append(t, args)
	return t
}
