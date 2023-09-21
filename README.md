# zlog
Zlog is formatted slog.Logger. [slog](https://pkg.go.dev/log/slog) provides structured logging, in which log records include a message, a severity level, and various other attributes expressed as key-value pairs, zlog make it more colorfull, beautifull and readable.

## Getting Started
1. Import the zlog Package
```go
import "github.com/jeffry-luqman/zlog"
```

2. Override some config if needed
```go
zlog.TimeFormat := time.RFC3339Nano
```

3. Create a Logger Instance
```go
logger := zlog.New()
```

4. Call it as slog logger from everywhere
```go
logger.Debug("Hello, World!")
logger.Info("Hello, World!")
logger.Warn("Hello, World!", slog.String("foo", "bar"))
logger.Error("Hello, World!", slog.String("foo", "bar"))
```

## Example Usage
Here's an example of how to use zlog to log different types of messages:

### Example 1
```go
package main

import (
	"github.com/jeffry-luqman/zlog"
)

func main() {
	logger := zlog.New()
	logger.Debug("Hello, World!")
	logger.Info("Hello, World!")
	logger.Warn("Hello, World!", slog.String("foo", "bar"))
	logger.Error("Hello, World!", slog.String("foo", "bar"))
}
```

### Example 2
#### Code
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
#### Output
![image](https://github.com/jeffry-luqman/zlog/assets/11884257/591d196b-cf72-48d5-aee0-34b3dac019de)

## Customizing Log Output
You can customize the log output by modifying the logger instance and the log handler options. Refer to the zlog.HandlerOptions and other variables defined in the zlog package for customization options.

## Logging with Enhanced Formatting
You can use the Fmt function from zlog to format your log messages with various attributes and colors. Here's an example:
```go
// Log a message with bold red text
logger.Info(zlog.Fmt("This is an important message!", zlog.FmtBold, zlog.FgRed))
```
In the example above, we use zlog.Fmt to format the log message with bold and red text.

## Contribute
If you find issues or have suggestions for improvements, please open an issue or create a pull request on the GitHub repository.

Happy logging with zlog!
