package tools

import (
	"os"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	logFile *os.File
}

func (l *Logger) InitLogger() error {
	file, err := os.OpenFile(configs.PathToLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return errors.ErrOpenFile
	}
	defer l.logFile.Close()

	log.SetOutput(file)
	l.logFile = file

	log.SetFormatter(&nested.Formatter{
		HideKeys:    false,
		NoColors:    true,
		FieldsOrder: []string{"requireId", "urlPath"},
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

	return nil
}

func AccessLogStart(urlPath, remoteAddr, method, requireId string) {
	log.WithFields(log.Fields{
		"urlPath":    urlPath,
		"requireId":  requireId,
		"method":     method,
		"remoteAddr": remoteAddr,
	}).Info("Start of request processing")
}

func AccessLogEnd(urlPath, remoteAddr, method, requireId string, startTime time.Time) {
	log.WithFields(log.Fields{
		"urlPath":     urlPath,
		"requireId":   requireId,
		"method":      method,
		"remoteAddr":  remoteAddr,
		"reqDuration": time.Since(startTime),
	}).Info("End of request processing")
}

func LogInfo(urlPath, packageName, functionName, msg, requireId string) {
	log.WithFields(log.Fields{
		"urlPath":   urlPath,
		"requireId": requireId,
		"package":   packageName,
		"function":  functionName,
	}).Info(msg)
}

func LogError(urlPath, packageName, functionName, requireId string, err error) {
	log.WithFields(log.Fields{
		"urlPath":   urlPath,
		"requireId": requireId,
		"package":   packageName,
		"function":  functionName,
	}).Error(err)
}
