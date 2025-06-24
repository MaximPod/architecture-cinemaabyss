// This file is safe to edit. Once it exists it will not be overwritten
package slogger

import (
	"fmt"
	"log/slog"
	"os"
)

// LogConfig - конфигурация логгера
type LogConfig struct {
	// отладочный режим вкл/выкл (default: false)
	DebugOn bool `yaml:"debug_on"`
	// формат логгера text/json/p_text (default: json)
	HandlerFormat string `yaml:"handler_format"`
	// установить логгер как глобальный (default: false)
	SetAsDefaultLogger bool `yaml:"set_as_default_logger"`
}

// InitLogger - Инициализация логгера slog
//
// @param cfg - конфигурация логгера
// @return *slog.Logger - инициализированный логгер slog
func InitLogger(cfg LogConfig) *slog.Logger {

	// set default log format
	slogOpt := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}

			// Add stackTrace to error
			switch a.Value.Kind() {
			case slog.KindAny:
				switch v := a.Value.Any().(type) {
				case error:
					a.Value = slog.StringValue(fmt.Sprintf("%+v", v))
				}
			}

			return a
		},
	}

	// set level as debug
	if cfg.DebugOn {
		slogOpt.Level = slog.LevelDebug
	}

	// set log format
	var logger *slog.Logger
	switch cfg.HandlerFormat {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slogOpt))
	case "p_text":
		logger = slog.New(NewPrettyTextHandler(os.Stdout, &slogOpt))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slogOpt))
	}

	// set default logger
	if cfg.SetAsDefaultLogger {
		slog.SetDefault(logger)
	}

	return logger
}
