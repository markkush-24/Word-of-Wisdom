package zap

import (
	"fmt" //nolint:goimports
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os" //nolint:goimports
	"word_of_wisdom/config"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(cfg core.LoggerConfig) *ZapLogger {
	logLevel := getLoggerLevel(cfg.Level)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	stderrSyncer := zapcore.Lock(zapcore.AddSync(os.Stderr))

	zapLogger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), stderrSyncer, zap.NewAtomicLevelAt(logLevel)),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.LevelOf(zap.ErrorLevel)),
		zap.AddCallerSkip(2),
	).Sugar()

	if err := zapLogger.Sync(); err != nil {
		zapLogger.Error(fmt.Sprintf("Failed to sync logger: %v", err))
	}

	return &ZapLogger{logger: zapLogger}
}

func (z *ZapLogger) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}

func (z *ZapLogger) Debugf(template string, args ...interface{}) {
	z.logger.Debugf(template, args...)
}

func (z *ZapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z *ZapLogger) Infof(template string, args ...interface{}) {
	z.logger.Infof(template, args...)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

func (z *ZapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Warnf(template, args...)
}

func (z *ZapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z *ZapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Errorf(template, args...)
}

func (z *ZapLogger) DPanic(args ...interface{}) {
	z.logger.DPanic(args...)
}

func (z *ZapLogger) DPanicf(template string, args ...interface{}) {
	z.logger.DPanicf(template, args...)
}

func (z *ZapLogger) Panic(args ...interface{}) {
	z.logger.Panic(args...)
}

func (z *ZapLogger) Panicf(template string, args ...interface{}) {
	z.logger.Panicf(template, args...)
}

func (z *ZapLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}

func (z *ZapLogger) Fatalf(template string, args ...interface{}) {
	z.logger.Fatalf(template, args...)
}

func getLoggerLevel(lv string) zapcore.Level {
	levelMap := map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
	level, exists := levelMap[lv]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}
