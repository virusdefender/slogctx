# slogctx

[slog](https://pkg.go.dev/log/slog) is a new logging library for Go. 
It supports passing context with log entries, you can use it to pass request id, user id, etc.

However, the default handler in slog does not record the context values, this library provides a handler that does.

## Demo

```go
package main

import (
	"context"
	"github.com/virusdefender/slogctx"
	"log/slog"
	"os"
)

type UserId struct{}

func main() {

	h := slogctx.NewHandler(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		&slogctx.HandlerOptions{
			AttrsFromCtx: []any{"request_id", UserId{}},
			AttrsNameMap: map[any]string{UserId{}: "uid"},
		})

	logger := slog.New(h)
	ctx := context.WithValue(
		context.WithValue(context.Background(), UserId{}, 123),
		"request_id", "abcdef")
	logger.ErrorContext(ctx, "error msg with context")
}

```

the output is

```
time=.. level=ERROR msg="error msg with context" request_id=abcdef uid=123
```