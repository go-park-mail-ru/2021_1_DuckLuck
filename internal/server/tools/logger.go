package tools

import (
	"os"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	logFile *os.File
}

func (l *Logger) init() {
	file, err := os.OpenFile(configs.PathToLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("can't open file for log")
	}
	defer l.logFile.Close()

	log.SetOutput(file)
	l.logFile = file

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	switch configs.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func AccessLogStart(urlPath, remoteAddr, method, requireId string, startTime time.Time) {
	log.WithFields(log.Fields{
		"urlPath":    urlPath,
		"requireId":  requireId,
		"method":     method,
		"remoteAddr": remoteAddr,
		"startTime":  startTime,
	}).Info("[START] " + urlPath)
}

func AccessLogEnd(urlPath, remoteAddr, method, requireId string, startTime time.Time) {
	log.WithFields(log.Fields{
		"urlPath":     urlPath,
		"requireId":   requireId,
		"method":      method,
		"remoteAddr":  remoteAddr,
		"reqDuration": time.Since(startTime),
	}).Info("[END] " + urlPath)
}

func LogInfo(packageName, functionName, urlPath, msg, requireId string) {
	log.WithFields(log.Fields{
		"urlPath":   urlPath,
		"requireId": requireId,
		"package":   packageName,
		"function":  functionName,
	}).Info(msg)
}

func LogError(packageName, functionName, urlPath, requireId string, err error) {
	log.WithFields(log.Fields{
		"urlPath":   urlPath,
		"requireId": requireId,
		"package":   packageName,
		"function":  functionName,
	}).Error(err)
}
