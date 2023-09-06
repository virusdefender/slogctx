# slogctx

[slog](https://pkg.go.dev/log/slog) is a new logging library for Go. 
It supports passing context with log entries, you can use it to pass request id, user id, etc.

However, the default handler in slog does not record the context values, this library provides a handler that does.

## Demo

```go
func main() {
	h := slogctx.NewHandler(
		slog.NewTextHandler(os.Stdout, nil),
		&slogctx.HandlerOptions{
			AttrsFromCtx: []string{"uid", "request_id"},
			AttrPrefix:   "ctx_",
		})
	logger := slog.New(h)
	ctx := context.WithValue(
		context.WithValue(context.Background(), "uid", 123), 
		"request_id", "abcdef")
	logger.Error("error msg")
	logger.ErrorContext(ctx, "error msg with context")
}
```

the output is

```
time=.. level=ERROR msg="error msg" ctx_uid=<nil> ctx_request_id=<nil>
time=.. level=ERROR msg="error msg with context" ctx_uid=123 ctx_request_id=abcdef
```

Playground: https://go.dev/play/p/MU0opRrNbIW