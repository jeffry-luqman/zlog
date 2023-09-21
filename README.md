# zlog

## Example
```go
package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jeffry-luqman/zlog"
)

func main() {
	zlog.HandlerOptions = &slog.HandlerOptions{Level: slog.LevelDebug}
	zlog.FmtDuration = []int{zlog.FgMagenta, zlog.FmtItalic}
	zlog.FmtPath = []int{zlog.FgHiCyan}

	start := time.Now()
	time.Sleep(200 * time.Millisecond)

	zlog.New().Debug("heloo 1",
		slog.Duration(zlog.KeyDuration, time.Now().Sub(start)),
		slog.Int(zlog.KeyStatus, http.StatusOK),
		slog.String(zlog.KeyMethod, http.MethodGet),
		slog.String(zlog.KeyPath, "/api/products"),
		slog.String("foo", "bar"),
		slog.Int("baz", 123),
	)
	time.Sleep(time.Millisecond)

	zlog.New().Info("heloo 2",
		slog.Duration(zlog.KeyDuration, time.Now().Sub(start)),
		slog.Int(zlog.KeyStatus, http.StatusCreated),
		slog.String(zlog.KeyMethod, http.MethodPost),
		slog.String(zlog.KeyPath, "/api/products"),
		slog.String("foo", "bar"),
		slog.Int("baz", 123),
	)
	time.Sleep(time.Millisecond)

	zlog.New().Warn("heloo 3",
		slog.Duration(zlog.KeyDuration, time.Now().Sub(start)),
		slog.Int(zlog.KeyStatus, http.StatusBadRequest),
		slog.String(zlog.KeyMethod, http.MethodPut),
		slog.String(zlog.KeyPath, "/api/products/1"),
		slog.String("foo", "bar"),
		slog.Int("baz", 123),
	)
	time.Sleep(time.Millisecond)

	zlog.New().Error("heloo 4",
		slog.Duration(zlog.KeyDuration, time.Now().Sub(start)),
		slog.Int(zlog.KeyStatus, http.StatusInternalServerError),
		slog.String(zlog.KeyMethod, http.MethodPatch),
		slog.String(zlog.KeyPath, "/api/products/1"),
		slog.String("foo", "bar"),
		slog.Int("baz", 123),
	)
	time.Sleep(time.Millisecond)

	zlog.New().Info("heloo 5",
		slog.Duration(zlog.KeyDuration, time.Now().Sub(start)),
		slog.Int(zlog.KeyStatus, http.StatusNoContent),
		slog.String(zlog.KeyMethod, http.MethodDelete),
		slog.String(zlog.KeyPath, "/api/products/1"),
		slog.String("foo", "bar"),
		slog.Int("baz", 123),
	)
}
```

## Output
![image](https://github.com/jeffry-luqman/zlog/assets/11884257/591d196b-cf72-48d5-aee0-34b3dac019de)
