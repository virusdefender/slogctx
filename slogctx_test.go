package slogctx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
)

type LogItem struct {
	Source struct {
		File string `json:"file"`
		Line int    `json:"line"`
	}
	CtxUid       *int    `json:"ctx_uid"`
	CtxRequestId *string `json:"ctx_request_id"`
}

func TestSlogCtx(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewHandler(
		slog.NewJSONHandler(buf, &slog.HandlerOptions{AddSource: true}),
		&HandlerOptions{
			AttrsFromCtx: []string{"uid", "request_id"},
			AttrPrefix:   "ctx_",
		})
	logger := slog.New(h)

	ctx := context.WithValue(context.WithValue(context.Background(), "uid", 123), "request_id", "abcdef")
	logger.Error("error msg")
	logger.ErrorContext(ctx, "error msg with context")

	logString := buf.String()
	fmt.Println(logString)

	logs := make([]LogItem, 2)
	for idx, line := range strings.Split(strings.TrimSpace(logString), "\n") {
		err := json.Unmarshal([]byte(line), &logs[idx])
		if err != nil {
			t.Fatalf("unmarshal line %s failed: %s", line, err)
		}
		if filepath.Base(logs[idx].Source.File) != "slogctx_test.go" {
			t.Fatalf("file name %s is not expected", logs[idx].Source.File)
		}
	}
	if logs[0].CtxUid != nil || logs[0].CtxRequestId != nil {
		t.Fatalf("unexpected log line 1")
	}
	if *logs[1].CtxUid != 123 || *logs[1].CtxRequestId != "abcdef" {
		t.Fatalf("unexpected log line 2")
	}
}
