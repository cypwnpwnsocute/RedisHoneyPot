// @Title  log.go
// @Description A highly interactive honeypot supporting redis protocol
// @Author  Cy 2021.04.08
package main

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"sync"
	"time"
)

const logDir = "logs"
const logName = "log"
const MaxAge = time.Hour * 24 * 7
const RotationTime = time.Hour * 24

var log *logrus.Logger
var once sync.Once

func init() {
	log = GetLoggerInstance()
}

func GetLoggerInstance() *logrus.Logger {
	once.Do(func() {
		log = configLocalFilesystemLogger(logDir, logName, MaxAge, RotationTime)
	})

	return log
}

func configLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *logrus.Logger {
	filePath, err := createDir(logPath)
	baseLogPath := path.Join(filePath, logFileName)
	CheckError(err)
	logFilePath := baseLogPath + ".%Y%m%d%H%M"

	writer, err := rotatelogs.New(
		logFilePath,
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	pathMap := lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}
	lfHook := lfshook.NewHook(pathMap, &logrus.TextFormatter{})
	var Log *logrus.Logger
	Log = logrus.New()
	Log.SetReportCaller(true)
	Log.Hooks.Add(lfHook)

	return Log
}

func CheckError(err error) {
	if err != nil {
		log.Error(err)
	}
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func createDir(path string) (directPath string, err error) {
	absPath, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}

	direct := absPath + "/" + path + "/"

	haveDirect := CheckFileIsExist(direct)
	if !haveDirect {
		err = os.Mkdir(direct, os.ModePerm)
	}

	return direct, err
}
