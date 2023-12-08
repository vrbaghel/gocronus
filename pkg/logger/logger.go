package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/** creates and returns a new zap logger instance */
func NewLogger(level int) *zap.Logger {
	path := getCWD()
	createLogDir(path)

	encoder := getEncoder()
	writerSyncer := getLogWriter(path)
	logLevel := parseLoglevel(level)

	return zap.New(
		zapcore.NewCore(
			encoder,
			writerSyncer,
			logLevel,
		),
		zap.AddCaller(),
	)
}

/** parses & returns zapcode loglevel */
func parseLoglevel(level int) zapcore.Level {
	LOG_LEVEL := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
	}
	return LOG_LEVEL[level]
}

/** returns path of current working directory */
func getCWD() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalln("logger : unable to get cwd :: exiting program...")
	}
	return path
}

/** creates logging directory if it doesn't exist */
func createLogDir(path string) {
	if _, err := os.Stat(filepath.Join(path, LOG_DIRECTORY)); os.IsNotExist(err) {
		log.Println("logger : creating log directory", LOG_DIRECTORY)
		errCreatingDir := os.Mkdir(LOG_DIRECTORY, os.ModePerm)
		if errCreatingDir != nil {
			log.Fatalln("logger : failed to create log directory", LOG_DIRECTORY)
		}
	}
}

/** configure logger output via encoder
 * use console encoder for readability
 */
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

/** leverage lumberjack for log rotation */
func getLogWriter(path string) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(path, LOG_DIRECTORY, LOG_FILE),
		MaxSize:    LOG_FILE_MAX_SIZE,
		MaxBackups: LOG_FILE_MAX_NUMBER,
		MaxAge:     LOG_FILE_MAX_AGE,
		Compress:   LOG_FILE_ENABLE_COMPRESSION,
	}

	return zapcore.AddSync(lumberjackLogger)
}
