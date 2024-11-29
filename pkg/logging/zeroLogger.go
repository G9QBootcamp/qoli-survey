package logging

import (
	"os"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type zeroLogger struct {
	cfg    *config.Config
	logger *zerolog.Logger
}

var zeroSyncLogger *zerolog.Logger

var zeroLogLevels map[string]zerolog.Level = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"warn":  zerolog.WarnLevel,
	"info":  zerolog.InfoLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

func (l *zeroLogger) getLogLevel() zerolog.Level {
	level, exisits := zeroLogLevels[l.cfg.Logging.Level]
	if !exisits {
		return zerolog.DebugLevel
	}
	return level
}

func newZeroLog(cfg *config.Config) (logger *zeroLogger) {
	logger = &zeroLogger{cfg: cfg}

	logger.Init(cfg)
	return logger
}

func (l *zeroLogger) Init(cfg *config.Config) {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		file, err := os.OpenFile(cfg.Logging.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic("could not open log file")
		}
		logger := zerolog.New(file).With().Timestamp().Str("AppName", "qoli-survey").Str("Logger", "zerolog").Logger()
		zerolog.SetGlobalLevel(l.getLogLevel())

		zeroSyncLogger = &logger

	})
	l.logger = zeroSyncLogger

}

func (l *zeroLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.Debug().Str("Category", string(cat)).Str("SubCategory", string(sub)).Fields(prepareZeroLogKeys(extra)).Msg(msg)
}
func (l *zeroLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debug().Msgf(template, args...)
}

func (l *zeroLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.Warn().Str("Category", string(cat)).Str("SubCategory", string(sub)).Fields(prepareZeroLogKeys(extra)).Msg(msg)
}
func (l *zeroLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warn().Msgf(template, args...)
}

func (l *zeroLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	extra[TraceId] = generateTraceID()

	l.logger.Error().Str("Category", string(cat)).Str("SubCategory", string(sub)).Fields(prepareZeroLogKeys(extra)).Msg(msg)
}
func (l *zeroLogger) Errorf(template string, args ...interface{}) {
	l.logger.Error().Msgf(template, args...)
}

func (l *zeroLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	extra[TraceId] = generateTraceID()

	l.logger.Fatal().Str("Category", string(cat)).Str("SubCategory", string(sub)).Fields(prepareZeroLogKeys(extra)).Msg(msg)
}
func (l *zeroLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal().Msgf(template, args...)
}

func (l *zeroLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.Info().Str("Category", string(cat)).Str("SubCategory", string(sub)).Fields(prepareZeroLogKeys(extra)).Msg(msg)
}
func (l *zeroLogger) Infof(template string, args ...interface{}) {
	l.logger.Info().Msgf(template, args...)
}

func prepareZeroLogKeys(extra map[ExtraKey]interface{}) map[string]interface{} {
	params := map[string]interface{}{}

	for k, v := range extra {
		params[string(k)] = v
	}
	return params
}
