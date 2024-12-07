package database

import (
    "github.com/sirupsen/logrus"
    "os"
    "path/filepath"
    "time"
    "gorm.io/gorm/logger"
    "context"
)

// Logger 实现了 logrus 的接口，用于将日志写入文件
type GormLogger struct {
    logger *logrus.Logger
}

// NewLogger 创建并配置一个新的 logrus Logger 实例
func NewGormLogger() *GormLogger {
    logger := logrus.New()
    logger.SetLevel(logrus.DebugLevel) // 设置日志级别

    // 设置日志格式为 JSON
    logger.SetFormatter(&logrus.TextFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    })

    // 设置日志输出到文件
    logFilePath := filepath.Join("./logs", "gorm.log")
    if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
        os.MkdirAll("logs", os.ModePerm)
        os.Create(logFilePath)
    }
    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        logger.Fatal("打开日志文件失败：", err)
        return nil
    }
    logger.SetOutput(file)

    return &GormLogger{logger: logger}
}

// LogMode 实现了 GORM 的 LogMode 方法
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
    switch level {
    case logger.Silent:
        l.logger.SetLevel(logrus.PanicLevel)
    case logger.Error:
        l.logger.SetLevel(logrus.ErrorLevel)
    case logger.Warn:
        l.logger.SetLevel(logrus.WarnLevel)
    case logger.Info:
        l.logger.SetLevel(logrus.InfoLevel)
    }
    return l
}

// Info 实现了 GORM 的 Info 方法
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
    l.logger.WithTime(time.Now()).Infof(msg, data...)
}

// Warn 实现了 GORM 的 Warn 方法
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
    l.logger.WithTime(time.Now()).Warnf(msg, data...)
}

// Error 实现了 GORM 的 Error 方法
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
    l.logger.WithTime(time.Now()).Errorf(msg, data...)
}

// Trace 实现了 GORM 的 Trace 方法
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    elapsed := time.Since(begin)
    sql, rows := fc()
    if err != nil {
        l.logger.WithTime(time.Now()).WithError(err).Errorf("Execute SQL: %s, rows affected: %d, duration: %v", sql, rows, elapsed)
    } else {
        l.logger.WithTime(time.Now()).Infof("Execute SQL: %s, rows affected: %d, duration: %v", sql, rows, elapsed)
    }
}
