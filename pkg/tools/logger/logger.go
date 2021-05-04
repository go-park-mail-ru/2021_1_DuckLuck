package logger

import (
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	logFile *os.File
}

func (l *Logger) InitLogger(pathToLogFile, logLevel string) error {
	file, err := os.OpenFile(pathToLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer l.logFile.Close()

	log.SetOutput(file)
	l.logFile = file

	log.SetFormatter(&nested.Formatter{
		HideKeys:    false,
		NoColors:    true,
		FieldsOrder: []string{"requireId", "urlPath"},
	})

	switch logLevel {
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

func HttpAccessLogStart(urlPath, remoteAddr, method, requireId string) {
	log.WithFields(log.Fields{
		"urlPath":    urlPath,
		"requireId":  requireId,
		"method":     method,
		"remoteAddr": remoteAddr,
	}).Info("Start of request processing")
}

func HttpAccessLogEnd(urlPath, remoteAddr, method, requireId string, startTime time.Time) {
	log.WithFields(log.Fields{
		"urlPath":     urlPath,
		"requireId":   requireId,
		"method":      method,
		"remoteAddr":  remoteAddr,
		"reqDuration": time.Since(startTime),
	}).Info("End of request processing")
}

func GrpcAccessLogStart(fullMethode, requireId, req, md string) {
	log.WithFields(log.Fields{
		"call":      fullMethode,
		"requireId": requireId,
		"req":       req,
		"md":        md,
	}).Info("Start of request processing")
}

func GrpcAccessLogEnd(fullMethode, requireId, reply string, startTime time.Time) {
	log.WithFields(log.Fields{
		"call":        fullMethode,
		"requireId":   requireId,
		"reply":       reply,
		"reqDuration": time.Since(startTime),
	}).Info("End of request processing")
}

func LogInfo(packageName, functionName, msg, requireId string) {
	log.WithFields(log.Fields{
		"requireId": requireId,
		"package":   packageName,
		"function":  functionName,
	}).Info(msg)
}

func LogError(packageName, functionName, requireId string, err error) {
	log.WithFields(log.Fields{
		"requireId": requireId,
		"package":   packageName,
		"function":  functionName,
	}).Error(err)
}
