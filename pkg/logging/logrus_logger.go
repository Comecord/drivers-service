package logging

import (
	"drivers-service/config"
	"fmt"
	"github.com/google/uuid"
	"github.com/mattn/go-colorable"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var logrusLogger *logrus.Logger

type logrusLoggerWrapper struct {
	cfg *config.Config
}

type CustomFormatter struct {
	logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Message = "\033[1;34m" + entry.Message + "\033[0m" // Пример для синего цвета
	return f.TextFormatter.Format(entry)
}

func newLogrusLogger(cfg *config.Config) *logrusLoggerWrapper {
	logger := &logrusLoggerWrapper{cfg: cfg}
	logger.Init()
	return logger
}

func (l *logrusLoggerWrapper) getLogLevel() logrus.Level {
	switch l.cfg.Logger.Level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error", "fatal":
		return logrus.ErrorLevel
	default:
		return logrus.DebugLevel
	}
}

func (l *logrusLoggerWrapper) Init() {
	if logrusLogger != nil {
		return
	}

	logrusLogger = logrus.New()
	logrusLogger.SetLevel(l.getLogLevel())

	logrusLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		ForceColors:      true,
		TimestampFormat:  time.RFC3339,
		QuoteEmptyFields: true,
	})

	logrusLogger.SetOutput(colorable.NewColorableStdout())

}

func (l *logrusLoggerWrapper) InitFile() {
	if logrusLogger != nil {
		return
	}

	fileName := fmt.Sprintf("%s%s-%s.%s", l.cfg.Logger.FilePath, time.Now().Format("2024-01-02"), uuid.New(), "json")
	logrusLogger = logrus.New()
	logrusLogger.SetLevel(l.getLogLevel())

	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	logrusLogger.SetOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1, // MB
		MaxBackups: 5,
		MaxAge:     20, // days
		LocalTime:  true,
		Compress:   true,
	})
}

func (l *logrusLoggerWrapper) register() {

	if logrusLogger != nil {
		return
	}

	const (
		fileFlag       = os.O_RDWR | os.O_CREATE | os.O_APPEND
		filePermission = 0666
	)

	logrusLogger = logrus.New()

	fileName := l.getFileName()
	file, err := os.OpenFile(fileName, fileFlag, filePermission)
	if err != nil {
		l.Fatalf("error opening file: %v", err)
	}

	chmodErr := file.Chmod(filePermission)
	if chmodErr != nil {
		l.Fatalf("error chmod file: %v", chmodErr)
	}

	writer := io.MultiWriter(os.Stdout, file)

	logrusLogger.SetOutput(writer)

	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
}

func (l *logrusLoggerWrapper) getFileName() string {
	date := l.getCurrentDate()
	return "./logs/golang-" + date + ".log"
}

func (l *logrusLoggerWrapper) getCurrentDate() string {
	timezone := l.cfg.TimeZone

	location, err := time.LoadLocation(timezone)
	if err != nil {
		l.Fatalf(err.Error())
	}

	timeNow := time.Now()
	return timeNow.In(location).Format(time.DateOnly)
}

func (l *logrusLoggerWrapper) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogrusInfo(cat, sub, extra)
	logrusLogger.WithFields(params).Debug(msg)
}

func (l *logrusLoggerWrapper) Debugf(template string, args ...interface{}) {
	logrusLogger.Debugf(template, args...)
}

func (l *logrusLoggerWrapper) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogrusInfo(cat, sub, extra)
	logrusLogger.WithFields(params).Info(msg)
}

func (l *logrusLoggerWrapper) Infof(template string, args ...interface{}) {
	logrusLogger.Infof(template, args...)
}

func (l *logrusLoggerWrapper) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogrusInfo(cat, sub, extra)
	logrusLogger.WithFields(params).Warn(msg)
}

func (l *logrusLoggerWrapper) Warnf(template string, args ...interface{}) {
	logrusLogger.Warnf(template, args...)
}

func (l *logrusLoggerWrapper) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogrusInfo(cat, sub, extra)
	logrusLogger.WithFields(params).Error(msg)
}

func (l *logrusLoggerWrapper) Errorf(template string, args ...interface{}) {
	logrusLogger.Errorf(template, args...)
}

func (l *logrusLoggerWrapper) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogrusInfo(cat, sub, extra)
	logrusLogger.WithFields(params).Fatal(msg)
}

func (l *logrusLoggerWrapper) Fatalf(template string, args ...interface{}) {
	logrusLogger.Fatalf(template, args...)
}

// Replace Category, SubCategory, ExtraKey with your actual types used in prepareLogInfo
func prepareLogrusInfo(cat Category, sub SubCategory, extra map[ExtraKey]interface{}) logrus.Fields {
	fields := logrus.Fields{}
	fields["Category"] = cat
	fields["SubCategory"] = sub

	for k, v := range extra {
		fields[string(k)] = v
	}
	return fields
}
