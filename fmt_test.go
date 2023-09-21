package zlog

import "testing"

var testCases = []struct {
	text     string
	attrs    []int
	expected string
}{
	{"Hello, World!", []int{FmtBold, FgRed}, "\x1b[1;31mHello, World!\x1b[0m"},
	{"Testing 123", []int{FmtUnderline, FgBlue}, "\x1b[4;34mTesting 123\x1b[0m"},
	{"Plain Text", []int{}, "Plain Text"},
	{"12345", []int{FmtFaint, FgYellow}, "\x1b[2;33m12345\x1b[0m"},
	{"Important Message", []int{FmtItalic, FgHiGreen}, "\x1b[3;92mImportant Message\x1b[0m"},
	{"Custom Format", []int{FmtBlinkRapid, BgRed}, "\x1b[6;41mCustom Format\x1b[0m"},
	{"Multi Attributes", []int{FmtBold, FgBlue, FmtUnderline}, "\x1b[1;34;4mMulti Attributes\x1b[0m"},
	{"Background Color", []int{BgCyan, FmtBold}, "\x1b[46;1mBackground Color\x1b[0m"},
	{"Combined Attributes", []int{FmtItalic, FgMagenta, BgHiYellow, FmtReverseVideo}, "\x1b[3;35;103;7mCombined Attributes\x1b[0m"},
	{"Empty Attributes", []int{}, "Empty Attributes"},
}

func TestFmt(t *testing.T) {
	for _, tc := range testCases {
		if formated := Fmt(tc.text, tc.attrs...); tc.expected != formated {
			t.Errorf("Expected formated [%v], got [%v]", tc.expected, formated)
		}
	}
}

func BenchmarkFmt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			Fmt(tc.text, tc.attrs...)
		}
	}
}
