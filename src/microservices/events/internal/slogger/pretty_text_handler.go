// This file is safe to edit. Once it exists it will not be overwritten
package slogger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// PrettyTextHandler is a slog.Handler that formats log records as formated text by errors.
type PrettyTextHandler struct {
	slog.Handler
	out io.Writer
}

// Создание нового обработчика текстового логирования с форматированием
//
// @param out - место назначения для вывода логирования
// @param opts - опции обработчика
// @return *PrettyTextHandler - новый обработчик текстового логирования с форматированием
func NewPrettyTextHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyTextHandler {
	h := PrettyTextHandler{
		Handler: slog.NewTextHandler(out, opts),
		out:     out,
	}

	return &h
}

// Handle - обработчик логирования
func (h *PrettyTextHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelError, slog.LevelInfo, slog.LevelDebug:
		h.FormatTextOutput(&r)
	default:
		h.Handler.Handle(ctx, r)
	}

	return nil
}

// FormatTextOutput - выполняет форматированный вывод данных записи в out handler-а
func (h *PrettyTextHandler) FormatTextOutput(r *slog.Record) error {
	fmt.Fprintf(h.out, "level=%s msg=%s\n", r.Level.String(), strings.Trim(r.Message, "\""))

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.out, "\t%s=%+v\n", a.Key, a.Value.Any())
		return true
	})

	return nil
}
