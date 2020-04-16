package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goft-cloud/http-proxy/config"
	log "github.com/goft-cloud/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	noticeLevels = []log.Level{
		log.DebugLevel,
		log.NoticeLevel,
		log.InfoLevel,
		log.TraceLevel,
	}

	errorLevels = []log.Level{
		log.ErrorLevel,
		log.WarnLevel,
		log.FatalLevel,
		log.PanicLevel,
	}

	loggerInstance = newLogger()

	lc = &loggerConfig{}
)

func Init() (err error) {
	_ = config.DecodeKey("log", lc)

	fmt.Println(lc)

	var (
		formatter    log.Formatter
		noticeFile   string
		errorFile    string
		noticeWriter io.Writer
		errorWriter  io.Writer
	)

	formatter = &log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}

	path := "."
	noticeFile = lc.ErrorFile
	errorFile = lc.NoticeFile

	if false == filepath.IsAbs(noticeFile) {
		noticeFile = path + string(os.PathSeparator) + noticeFile
	}
	if false == filepath.IsAbs(errorFile) {
		errorFile = path + string(os.PathSeparator) + errorFile
	}

	// 输出文件
	noticeWriter, err = os.OpenFile(noticeFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	errorWriter, err = os.OpenFile(errorFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	loggerInstance.setLevel(log.NoticeLevel)

	// Writers
	loggerInstance.setFormatter(formatter)
	loggerInstance.setNoticeWriter(noticeWriter)
	loggerInstance.setErrorWriter(errorWriter)

	return loggerInstance.Init()
}

// 获取log实例
func Logger() *log.Logger {
	return loggerInstance.Logger()
}

func newLogger() *logger {
	return &logger{
		logger: log.New(),
		lock:   &sync.RWMutex{},
	}
}

type loggerConfig struct {
	NoticeFile string `toml:"notice_file"`
	ErrorFile  string `toml:"error_file"`
}

type logger struct {
	noticeWriter io.Writer
	errorWriter  io.Writer
	logger       *log.Logger
	isInit       bool
	formatter    log.Formatter
	level        log.Level
	lock         *sync.RWMutex
}

// 写入level
func (l *logger) setLevel(level log.Level) {
	if false == l.isInit {
		l.level = level
	}
}

func inLevel(level log.Level, levelList []log.Level) bool {
	for _, v := range levelList {
		if v == level {
			return true
		}
	}

	return false
}

// 设置format
func (l *logger) setFormatter(formatter log.Formatter) {
	l.formatter = formatter
}

// 设置notice的writer
func (l *logger) setNoticeWriter(writer io.Writer) {
	if false == l.isInit {
		l.noticeWriter = writer
	}
}

// 设置error的writer
func (l *logger) setErrorWriter(writer io.Writer) {
	if false == l.isInit {
		l.errorWriter = writer
	}
}

// 获取Logger实例
func (l logger) Logger() *log.Logger {
	return l.logger
}

func (l *logger) Levels() []log.Level {
	return log.AllLevels
}

func (l *logger) Fire(entry *log.Entry) error {
	fmt.Println(entry.Message)
	var (
		err error
		msg []byte
	)

	l.lock.Lock()
	defer l.lock.Unlock()

	if false == l.isInit {
		return nil
	}

	msg, err = l.formatter.Format(entry)
	if err != nil {
		return err
	}

	if inLevel(entry.Level, noticeLevels) && l.noticeWriter != nil {
		_, err = l.noticeWriter.Write(msg)
		return err
	}

	if inLevel(entry.Level, errorLevels) && l.errorWriter != nil {
		_, err = l.errorWriter.Write(msg)
		return err
	}

	return nil
}

// 初始化logger
func (l *logger) Init() error {
	l.logger.SetFormatter(l.formatter)
	l.logger.SetLevel(l.level)
	l.logger.AddHook(l)
	l.logger.SetOutput(ioutil.Discard)

	l.isInit = true

	return nil
}

// Error
func Error(context *gin.Context, args ...interface{}) {
	GetEntry(context).Error(args...)
}

// Warn
func Warn(context *gin.Context, args ...interface{}) {
	GetEntry(context).Warn(args...)
}

// Info
func Info(context *gin.Context, args ...interface{}) {
	GetEntry(context).Info(args...)
}

// Trace
func Trace(context *gin.Context, args ...interface{}) {
	GetEntry(context).Trace(args...)
}

// Debug
func Debug(context *gin.Context, args ...interface{}) {
	GetEntry(context).Debug(args...)
}

func GetEntry(context *gin.Context) *log.Entry {
	var (
		traceId  interface{}
		spanId   interface{}
		parentId interface{}
	)
	traceId, _ = context.Get("traceid")
	spanId, _ = context.Get("spanid")
	parentId, _ = context.Get("parentid")

	return loggerInstance.Logger().WithField("traceid", traceId).
		WithField("spanid", spanId).
		WithField("parentid", parentId).
		WithField("application", config.App.Name)
}
