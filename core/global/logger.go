package global

import (
	"be-better/utils"
	"bytes"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
	"strings"
)

func Logger() *logrus.Logger {

	var logPath = GlobalConfig.System.LogPath
	logPathExists, _ := utils.PathExists(logPath)
	if !logPathExists {
		os.Mkdir(logPath, os.ModePerm)
	}

	var log = logrus.New()
	log.SetFormatter(formatter(false))
	hook := NewSyslogHook()
	log.Hooks.Add(hook)

	log.SetLevel(logrus.DebugLevel)

	return log
}

type SyslogHook struct {
	IsHook bool
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewSyslogHook() *SyslogHook {
	return &SyslogHook{}
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	var log = hook.getLogger(entry)
	log.SetLevel(logrus.DebugLevel)
	log.Log(entry.Level, line)

	return nil
}

var fileNameMap = map[string]*rotatelogs.RotateLogs{}

func (hook *SyslogHook) getLogger(entry *logrus.Entry) *logrus.Logger {
	var log = logrus.New()
	log.SetFormatter(formatter(true))

	var filename = "Info"
	switch entry.Level {
	case logrus.PanicLevel:
	case logrus.FatalLevel:
	case logrus.ErrorLevel:
		filename = "Error"
		break
	case logrus.WarnLevel:
	case logrus.InfoLevel:
	case logrus.DebugLevel, logrus.TraceLevel:
		filename = "Info"
	}

	writer, _ := rotatelogs.New(GlobalConfig.System.LogPath + "/beBetter" + filename + "_%Y%m%d.log")
	log.SetOutput(writer)
	return log
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func formatter(isHook bool) *SyslogHook {
	fmtter := &SyslogHook{
		IsHook: isHook,
	}
	return fmtter
}

// Format an log entry
func (f *SyslogHook) Format(entry *logrus.Entry) ([]byte, error) {
	timestampFormat := "2006-01-02 15:04:05.999"

	// output buffer
	b := &bytes.Buffer{}

	if !f.IsHook {
		// write time
		b.WriteString("[")
		b.WriteString(entry.Time.Format(timestampFormat))
		b.WriteString("]")

		// write level
		b.WriteString(" [")
		level := strings.ToUpper(entry.Level.String())
		b.WriteString(level)
		b.WriteString("]")
	}

	// write fields
	if !f.IsHook {
		b.WriteString(" ")
	}
	f.writeFields(b, entry)
	b.WriteString(entry.Message)

	if !f.IsHook {
		b.WriteString("\n")
	}

	return b.Bytes(), nil
}

func (f SyslogHook) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		fmt.Fprintf(
			b,
			" (%s:%d %s)",
			entry.Caller.File,
			entry.Caller.Line,
			entry.Caller.Function,
		)
	}
}

func (f SyslogHook) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *SyslogHook) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
	b.WriteString(" ")
}
