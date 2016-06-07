package context_log

import (
	"bytes"
	"fmt"
	"github.com/rs/xid"
	"time"
)

type ServerContext struct {
	buf   *bytes.Buffer
	Uuid  string
	sTime time.Time
	tTime time.Time //临时统计用的时间
}

func NewContext(msg string) *ServerContext {
	sc := new(ServerContext)
	sc.buf = bytes.NewBufferString(msg)
	sc.Uuid = xid.New().String()
	sc.sTime = time.Now()
	//sc.tTime
	return sc
}

func (sc *ServerContext) StartTimer() {
	sc.tTime = time.Now()
}

func (sc *ServerContext) StopTimer(key string) {
	duration := time.Now().Sub(sc.tTime)
	sc.buf.WriteString(fmt.Sprintf(" %s=%v", key, duration))
}

func (sc *ServerContext) AddNotes(key string, val string) {
	sc.buf.WriteString(fmt.Sprintf(" %s=%s", key, val))
}
func (sc *ServerContext) Flush() {
	duration := time.Now().Sub(sc.sTime)
	Log.Info(fmt.Sprintf("%s=%s cost=%v %s ", "Uuid", sc.Uuid, duration, sc.buf.String()))
}
func (sc *ServerContext) Debug(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Debug(s, args...)
}
func (sc *ServerContext) Info(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Info(s, args...)
}
func (sc *ServerContext) Notice(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Notice(s, args...)
}
func (sc *ServerContext) Warning(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Warning(s, args...)
}
func (sc *ServerContext) Error(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Error(s, args...)
}
func (sc *ServerContext) Critical(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.Uuid, format)
	Log.Critical(s, args...)
}
