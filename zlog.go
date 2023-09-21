package zlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/mattn/go-colorable"
)

// Base attributes
const (
	FmtReset int = iota
	FmtBold
	FmtFaint
	FmtItalic
	FmtUnderline
	FmtBlinkSlow
	FmtBlinkRapid
	FmtReverseVideo
	FmtConcealed
	FmtCrossedOut
)

// Foreground colors
const (
	FgBlack int = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity colors
const (
	FgHiBlack int = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background colors
const (
	BgBlack int = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity colors
const (
	BgHiBlack int = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

var (
	HandlerOptions *slog.HandlerOptions
	TimeFormat     = "[15:04:05.000]"

	KeyStatus    = "status"
	KeyDuration  = "duration"
	KeyMethod    = "method"
	KeyPath      = "path"
	KeyDelimiter = "="

	FmtLevelDebug = []int{FgHiMagenta}
	FmtLevelInfo  = []int{FgGreen}
	FmtLevelWarn  = []int{FgYellow}
	FmtLevelError = []int{FgHiRed}

	FmtStatus1XX     = []int{FgGreen}
	FmtStatus2XX     = []int{FgGreen}
	FmtStatus3XX     = []int{FgGreen}
	FmtStatus4XX     = []int{FgYellow}
	FmtStatus5XX     = []int{FgHiRed}
	FmtStatusUnknown = []int{FgHiRed}

	FmtTime     = []int{FgHiBlack}
	FmtDuration = []int{FgCyan, FmtItalic}
	FmtPath     = []int{FgHiCyan}
	FmtMessage  = []int{FmtReset}

	FmtAttrKey       = []int{FgBlue}
	FmtAttrDelimiter = []int{FgHiBlack}
	FmtAttrValue     = []int{FgYellow}
)

var (
	w   io.Writer
	log *slog.Logger
)

func New() *slog.Logger {
	if w == nil {
		w = colorable.NewColorableStdout()
	}
	if log == nil {
		log = slog.New(&logHandler{sh: slog.NewTextHandler(w, HandlerOptions)})
	}
	return log
}

// Fmt format log with attribute
//
//	for example :
//	  log.Fmt("text", fmtBold, FgRed)
//
//	output (text with bold red foreground) :
//	  \x1b[1;31mtext\x1b[0m
func Fmt(s string, attribute ...int) string {
	format := make([]string, len(attribute))
	for i, v := range attribute {
		format[i] = strconv.Itoa(v)
	}
	return "\x1b[" + strings.Join(format, ";") + "m" + s + "\x1b[0m"
}

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
	level := r.Level.String()
	statusCode := ""
	duration := ""
	method := ""
	path := ""
	switch r.Level {
	case slog.LevelDebug:
		level = Fmt(fmt.Sprintf("%6s", level), FmtLevelDebug...)
	case slog.LevelInfo:
		level = Fmt(fmt.Sprintf("%6s", level), FmtLevelInfo...)
	case slog.LevelWarn:
		level = Fmt(fmt.Sprintf("%6s", level), FmtLevelWarn...)
	case slog.LevelError:
		level = Fmt(fmt.Sprintf("%6s", level), FmtLevelError...)
	}
	attrs := []slog.Attr{}
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == KeyStatus {
			statusCode = a.Value.String()
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
		} else if a.Key == KeyDuration {
			duration = a.Value.String()
			if a.Value.Kind() == slog.KindDuration {
				duration = a.Value.Duration().String()
			}
		} else if a.Key == KeyMethod {
			method = a.Value.String()
			switch method {
			case http.MethodGet:
				method = Fmt(fmt.Sprintf("%7s", method), FgGreen)
			case http.MethodPost:
				method = Fmt(fmt.Sprintf("%7s", method), FgYellow)
			case http.MethodPut:
				method = Fmt(fmt.Sprintf("%7s", method), FgBlue)
			case http.MethodPatch:
				method = Fmt(fmt.Sprintf("%7s", method), FgCyan)
			case http.MethodDelete:
				method = Fmt(fmt.Sprintf("%7s", method), FgRed)
			default:
				method = fmt.Sprintf("%7s", method)
			}
		} else if a.Key == KeyPath {
			path = a.Value.String()
		} else {
			attrs = append(attrs, a)
		}
		return true
	})

	b := &bytes.Buffer{}
	h.write(b, level, " ")
	h.write(b, Fmt(r.Time.Format(TimeFormat), FmtTime...), " ")
	if statusCode != "" {
		h.write(b, statusCode, " ")
	}
	if duration != "" {
		h.write(b, Fmt(fmt.Sprintf("%12s", duration), FmtDuration...), " ")
	}
	if method != "" {
		h.write(b, method, " ")
	}
	if path != "" {
		h.write(b, Fmt(path, FmtPath...), " ")
	}
	h.write(b, Fmt(r.Message, FmtMessage...), " ")
	for _, a := range attrs {
		h.write(b, Fmt(a.Key, FmtAttrKey...), Fmt(KeyDelimiter, FmtAttrDelimiter...), Fmt(a.Value.String(), FmtAttrValue...), " ")
	}
	h.write(b, "\n")
	_, err := b.WriteTo(w)
	return err
}
func (h *logHandler) write(b *bytes.Buffer, args ...string) {
	for _, s := range args {
		b.WriteString(s)
	}
}
