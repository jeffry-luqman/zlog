package zlog

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

// This code is modified from https://github.com/golang/go/blob/54f78cf8f1b8deea787803aeff5fb6150d7fac8f/src/log/slog/logger_test.go#L27
func TestZlog(t *testing.T) {
	ctx := context.Background()
	w = &bytes.Buffer{}
	TimeFormat = "Z"

	l := New()

	check := func(want []byte) {
		t.Helper()
		x, ok := w.(*bytes.Buffer)
		if ok {
			t.Helper()
			got := x.Bytes()
			if !bytes.Equal(got, want) {
				t.Errorf("\ngot  %s\nwant %s", string(got), string(want))
			}
			x.Reset()
		}
	}

	l.Info("msg", "a", 1, "b", 2)
	check([]byte{27, 91, 51, 50, 109, 73, 78, 70, 79, 32, 32, 27, 91, 48, 109, 32, 27, 91, 57, 48, 109, 90, 27, 91, 48, 109, 32, 27, 91, 48, 109, 109, 115, 103, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 97, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 49, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 98, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 50, 27, 91, 48, 109, 32, 10})

	// By default, debug messages are not printed.
	l.Debug("bg", slog.Int("a", 1), "b", 2)
	check([]byte{})

	l.Warn("w", slog.Duration("dur", 3*time.Second))
	check([]byte{27, 91, 57, 51, 109, 87, 65, 82, 78, 32, 32, 27, 91, 48, 109, 32, 27, 91, 57, 48, 109, 90, 27, 91, 48, 109, 32, 27, 91, 48, 109, 119, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 100, 117, 114, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 51, 115, 27, 91, 48, 109, 32, 10})

	l.Error("bad", "a", 1)
	check([]byte{27, 91, 57, 49, 109, 69, 82, 82, 79, 82, 32, 27, 91, 48, 109, 32, 27, 91, 57, 48, 109, 90, 27, 91, 48, 109, 32, 27, 91, 48, 109, 98, 97, 100, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 97, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 49, 27, 91, 48, 109, 32, 10})

	l.Log(ctx, slog.LevelWarn+1, "w", slog.Int("a", 1), slog.String("b", "two"))
	check([]byte{87, 65, 82, 78, 43, 49, 32, 27, 91, 57, 48, 109, 90, 27, 91, 48, 109, 32, 27, 91, 48, 109, 119, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 97, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 49, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 98, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 116, 119, 111, 27, 91, 48, 109, 32, 10})

	l.LogAttrs(ctx, slog.LevelInfo+1, "a b c", slog.Int("a", 1), slog.String("b", "two"))
	check([]byte{73, 78, 70, 79, 43, 49, 32, 27, 91, 57, 48, 109, 90, 27, 91, 48, 109, 32, 27, 91, 48, 109, 97, 32, 98, 32, 99, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 97, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 49, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 98, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 116, 119, 111, 27, 91, 48, 109, 32, 10})

	l.Info("info", "a", []slog.Attr{slog.Int("i", 1)})
	check([]byte{27, 91, 51, 50, 109, 73, 78, 70, 79, 32, 32, 27, 91, 48, 109, 32, 27, 91, 57, 48, 109, 90, 27, 91, 48, 109, 32, 27, 91, 48, 109, 105, 110, 102, 111, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 97, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 91, 105, 61, 49, 93, 27, 91, 48, 109, 32, 10})

	l.Info("info", "a", slog.GroupValue(slog.Int("i", 1)))
	check([]byte{27, 91, 51, 50, 109, 73, 78, 70, 79, 32, 32, 27, 91, 48, 109, 32, 27, 91, 57, 48, 109, 90, 27, 91, 48, 109, 32, 27, 91, 48, 109, 105, 110, 102, 111, 27, 91, 48, 109, 32, 27, 91, 51, 52, 109, 97, 27, 91, 48, 109, 27, 91, 57, 48, 109, 61, 27, 91, 48, 109, 27, 91, 51, 51, 109, 91, 105, 61, 49, 93, 27, 91, 48, 109, 32, 10})
}

// This is a simple benchmark. See the benchmarks subdirectory for more extensive ones.
// This code is modified from https://github.com/golang/go/blob/54f78cf8f1b8deea787803aeff5fb6150d7fac8f/src/log/slog/logger_test.go#L522
func BenchmarkSlog(b *testing.B) {
	ctx := context.Background()
	l := slog.New(slog.NewTextHandler(io.Discard, HandlerOptions))
	b.Run("no attrs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(ctx, slog.LevelInfo, "msg")
		}
	})
	b.Run("attrs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
		}
	})
	b.Run("attrs-parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
			}
		})
	})
	b.Run("keys-values", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.Log(ctx, slog.LevelInfo, "msg", "a", 1, "b", "two", "c", true)
		}
	})
	b.Run("WithContext", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(ctx, slog.LevelInfo, "msg2", slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
		}
	})
	b.Run("WithContext-parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
			}
		})
	})
}

// This is a simple benchmark. See the benchmarks subdirectory for more extensive ones.
// This code is modified from https://github.com/golang/go/blob/54f78cf8f1b8deea787803aeff5fb6150d7fac8f/src/log/slog/logger_test.go#L522
func BenchmarkZlog(b *testing.B) {
	w = io.Discard
	ctx := context.Background()
	l := New()
	b.Run("no attrs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(ctx, slog.LevelInfo, "msg")
		}
	})
	b.Run("attrs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
		}
	})
	b.Run("attrs-parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
			}
		})
	})
	b.Run("keys-values", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.Log(ctx, slog.LevelInfo, "msg", "a", 1, "b", "two", "c", true)
		}
	})
	b.Run("WithContext", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(ctx, slog.LevelInfo, "msg2", slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
		}
	})
	b.Run("WithContext-parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
			}
		})
	})
}
