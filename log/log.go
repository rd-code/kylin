package log

import (
	"fmt"
	"io"
	"os"
)

type Level int

const (
	Info = iota
	Warn
	Error
)

//定义日志的接口
type Logger interface {
	//设置日志等级
	Level(level Level)
	Info(format string, args ...interface{}) error
	Warn(format string, args ...interface{}) error
	Error(format string, args ...interface{}) error
}

var Default = New(os.Stdout)
//生成日志写入类型
func New(writer io.Writer) Logger {
	return &loggerImpl{
		writer: writer,
	}
}

//默认的日志实现类
type loggerImpl struct {
	writer io.Writer
	level  Level
}

func (l *loggerImpl) Level(level Level) {
	l.level = level
}

//具体的写入实现
func (l *loggerImpl) write(format string, args ...interface{}) (err error) {
	_, err = fmt.Fprintf(l.writer, format+"\n", args...)
	return err
}

func (l *loggerImpl) Info(format string, args ...interface{}) (err error) {
	if l.level > Info {
		return
	}
	return l.write(format, args...)
}

func (l *loggerImpl) Warn(format string, args ...interface{}) (err error) {
	if l.level > Warn {
		return
	}
	return l.write(format, args...)
}

func (l *loggerImpl) Error(format string, args ...interface{}) (err error) {
	if l.level > Error {
		return
	}
	return l.write(format, args...)
}
