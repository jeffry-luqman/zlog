package zlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"

	"github.com/mattn/go-colorable"
)

// Base config for logger handler, you can override this with your own preferences.
var (
	// HandlerOptions contains options for the logger handler.
	HandlerOptions *slog.HandlerOptions
	// TimeFormat is the format used for logging timestamps.
	TimeFormat = "[15:04:05.000]"
)

// Some keys with specific format, you can override this with your own preferences.
var (
	KeyStatus    = "status"
	KeyDuration  = "duration"
	KeyMethod    = "method"
	KeyPath      = "path"
	KeyDelimiter = "="
	FieldOrder   = []string{
		slog.LevelKey,
		slog.TimeKey,
		KeyStatus,
		KeyDuration,
		KeyMethod,
		KeyPath,
		slog.SourceKey,
		slog.MessageKey,
	}
)

// Format attributes for some keys, you can override this with your own preferences.
var (
	FmtLevelDebug = []int{FgHiMagenta}
	FmtLevelInfo  = []int{FgGreen}
	FmtLevelWarn  = []int{FgHiYellow}
	FmtLevelError = []int{FgHiRed}

	FmtStatus1XX     = []int{FgGreen}
	FmtStatus2XX     = []int{FgGreen}
	FmtStatus3XX     = []int{FgGreen}
	FmtStatus4XX     = []int{FgHiYellow}
	FmtStatus5XX     = []int{FgHiRed}
	FmtStatusUnknown = []int{FgHiRed}

	FmtMethodGet    = []int{FgGreen}
	FmtMethodPost   = []int{FgYellow}
	FmtMethodPut    = []int{FgBlue}
	FmtMethodPatch  = []int{FgCyan}
	FmtMethodDelete = []int{FgRed}
	FmtMethodOther  = []int{FgMagenta}

	FmtTime     = []int{FgHiBlack}
	FmtDuration = []int{FgCyan, FmtItalic}
	FmtPath     = []int{FgHiCyan}
	FmtMessage  = []int{FmtReset}

	FmtAttrKey       = []int{FgBlue}
	FmtAttrDelimiter = []int{FgHiBlack}
	FmtAttrValue     = []int{FgYellow}
)

var (
	Writer io.Writer
	log    *slog.Logger
)

// New creates a new logger instance and returns it.
func New() *slog.Logger {
	if Writer == nil {
		Writer = colorable.NewColorableStdout()
	}
	if log == nil {
		log = slog.New(&logHandler{sh: slog.NewTextHandler(Writer, HandlerOptions)})
	}
	return log
}

// logHandler is a log handler that implements the slog.Handler interface.
type logHandler struct {
	sh slog.Handler
}

func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.sh.Enabled(ctx, level)
}
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{sh: h.sh.WithAttrs(attrs)}
}
func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{sh: h.sh.WithGroup(name)}
}
func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := []slog.Attr{}
	fields := map[string]string{}
	level := r.Level.String()
	switch r.Level {
	case slog.LevelDebug:
		level = Fmt(fmt.Sprintf("%-6s", level), FmtLevelDebug...)
	case slog.LevelInfo:
		level = Fmt(fmt.Sprintf("%-6s", level), FmtLevelInfo...)
	case slog.LevelWarn:
		level = Fmt(fmt.Sprintf("%-6s", level), FmtLevelWarn...)
	case slog.LevelError:
		level = Fmt(fmt.Sprintf("%-6s", level), FmtLevelError...)
	}
	if slices.Contains(FieldOrder, slog.LevelKey) {
		fields[slog.LevelKey] = level
	} else {
		attrs = append(attrs, slog.String(slog.LevelKey, level))
	}
	if slices.Contains(FieldOrder, slog.TimeKey) {
		fields[slog.TimeKey] = Fmt(r.Time.Format(TimeFormat), FmtTime...)
	} else {
		attrs = append(attrs, slog.String(slog.TimeKey, Fmt(r.Time.Format(TimeFormat), FmtTime...)))
	}
	if slices.Contains(FieldOrder, slog.MessageKey) {
		fields[slog.MessageKey] = Fmt(r.Message, FmtMessage...)
	} else {
		attrs = append(attrs, slog.String(slog.MessageKey, Fmt(r.Message, FmtMessage...)))
	}
	// fields[slog.SourceKey] = Fmt(r.Message, FmtMessage...) // todo

	r.Attrs(func(a slog.Attr) bool {
		if a.Key == KeyStatus && slices.Contains(FieldOrder, a.Key) {
			statusCode := a.Value.String()
			if a.Value.Kind() == slog.KindInt64 {
				code := a.Value.Int64()
				if code >= http.StatusInternalServerError {
					statusCode = Fmt(statusCode, FmtStatus5XX...)
				} else if code >= http.StatusBadRequest {
					statusCode = Fmt(statusCode, FmtStatus4XX...)
				} else if code >= http.StatusMultipleChoices {
					statusCode = Fmt(statusCode, FmtStatus3XX...)
				} else if code >= http.StatusOK {
					statusCode = Fmt(statusCode, FmtStatus2XX...)
				} else if code >= http.StatusContinue {
					statusCode = Fmt(statusCode, FmtStatus1XX...)
				} else {
					statusCode = Fmt(statusCode, FgGreen)
				}
			}
			fields[KeyStatus] = statusCode
		} else if a.Key == KeyDuration && slices.Contains(FieldOrder, a.Key) {
			duration := a.Value.String()
			if a.Value.Kind() == slog.KindDuration {
				duration = a.Value.Duration().String()
			}
			fields[KeyDuration] = Fmt(fmt.Sprintf("%12s", duration), FmtDuration...)
		} else if a.Key == KeyMethod && slices.Contains(FieldOrder, a.Key) {
			method := a.Value.String()
			switch method {
			case http.MethodGet:
				method = Fmt(fmt.Sprintf("%7s", method), FmtMethodGet...)
			case http.MethodPost:
				method = Fmt(fmt.Sprintf("%7s", method), FmtMethodPost...)
			case http.MethodPut:
				method = Fmt(fmt.Sprintf("%7s", method), FmtMethodPut...)
			case http.MethodPatch:
				method = Fmt(fmt.Sprintf("%7s", method), FmtMethodPatch...)
			case http.MethodDelete:
				method = Fmt(fmt.Sprintf("%7s", method), FmtMethodDelete...)
			default:
				method = Fmt(fmt.Sprintf("%7s", method), FmtMethodOther...)
			}
			fields[KeyMethod] = method
		} else if a.Key == KeyPath && slices.Contains(FieldOrder, a.Key) {
			fields[KeyPath] = a.Value.String()
		} else {
			attrs = append(attrs, a)
		}
		return true
	})

	b := &bytes.Buffer{}
	for _, key := range FieldOrder {
		val, ok := fields[key]
		if ok {
			h.write(b, val, " ")
		}
	}
	for _, a := range attrs {
		h.write(b, Fmt(a.Key, FmtAttrKey...), Fmt(KeyDelimiter, FmtAttrDelimiter...), Fmt(a.Value.String(), FmtAttrValue...), " ")
	}
	h.write(b, "\n")
	_, err := b.WriteTo(Writer)
	return err
}
func (h *logHandler) write(b *bytes.Buffer, args ...string) {
	for _, s := range args {
		b.WriteString(s)
	}
}
