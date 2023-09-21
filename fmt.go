package zlog

import (
	"bytes"
	"strconv"
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

// Fmt formats a string with the provided attributes and colors and returns the formatted string.
//
//	for example :
//	  zlog.Fmt("text", zlog.FmtBold, zlog.FgRed)
//
//	output :
//	  \x1b[1;31mtext\x1b[0m
func Fmt(s string, attrs ...int) string {
	if attrLen := len(attrs); attrLen > 0 {
		b := bytes.Buffer{}
		b.WriteString("\x1b[")
		for i, attr := range attrs {
			b.WriteString(strconv.Itoa(attr))
			if i < attrLen-1 {
				b.WriteRune(';')
			}
		}
		b.WriteRune('m')
		b.WriteString(s)
		b.WriteString("\x1b[0m")
		return b.String()
	}
	return s
}
