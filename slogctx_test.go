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
	CtxNotFound  *string `json:"ctx_not_found"`
}

func TestSourceLocation(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewHandler(
		slog.NewJSONHandler(buf, &slog.HandlerOptions{AddSource: true}),
		&HandlerOptions{})
	logger := slog.New(h)

	logger.Error("error msg")

	logString := buf.String()
	fmt.Println(logString)

	logItem := LogItem{}
	err := json.Unmarshal([]byte(logString), &logItem)
	if err != nil {
		t.Fatalf("unmarshal line %s failed: %s", logString, err)
	}
	if filepath.Base(logItem.Source.File) != "slogctx_test.go" {
		t.Fatalf("file name %s is not expected", logItem.Source.File)
	}
}

func testContextValues(t *testing.T, showNilValue bool) {
	buf := &bytes.Buffer{}
	h := NewHandler(
		slog.NewJSONHandler(buf, &slog.HandlerOptions{AddSource: true}),
		&HandlerOptions{
			AttrsFromCtx: []string{"uid", "request_id", "not_found"},
			AttrPrefix:   "ctx_",
			ShowNilValue: showNilValue,
		})
	logger := slog.New(h)

	ctx := context.WithValue(context.WithValue(context.Background(), "uid", 123), "request_id", "abcdef")
	logger.ErrorContext(ctx, "error msg with context")

	logString := buf.String()
	fmt.Println(logString)

	logItem := LogItem{}
	err := json.Unmarshal([]byte(logString), &logItem)
	if err != nil {
		t.Fatalf("unmarshal line %s failed: %s", logString, err)
	}
	if *logItem.CtxUid != 123 || *logItem.CtxRequestId != "abcdef" || logItem.CtxNotFound != nil {
		t.Fatalf("unexpected log line")
	}

	if showNilValue && !strings.Contains(logString, `"ctx_not_found":null`) {
		t.Fatalf("unexpected log line")
	} else if !showNilValue && strings.Contains(logString, `"ctx_not_found":null`) {
		t.Fatalf("unexpected log line")
	}
}

func TestContextValues(t *testing.T) {
	testContextValues(t, false)
	testContextValues(t, true)
}

func TestWithAttr(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewHandler(
		slog.NewTextHandler(buf, nil),
		&HandlerOptions{AttrsFromCtx: []string{"uid"}})

	logger := slog.New(h).With("foo", "bar")

	logger.ErrorContext(context.WithValue(context.Background(), "uid", 123), "error msg", "a", "b")

	logString := buf.String()
	fmt.Println(logString)
	if !strings.Contains(logString, `foo=bar a=b uid=123`) {
		t.Fatalf("unexpected log line")
	}
}

func TestWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewHandler(
		slog.NewTextHandler(buf, nil),
		&HandlerOptions{AttrsFromCtx: []string{"uid"}})

	logger := slog.New(h).WithGroup("group_name")

	logger.ErrorContext(context.WithValue(context.Background(), "uid", 123), "error msg", "a", "b")

	logString := buf.String()
	fmt.Println(logString)
	if !strings.Contains(logString, `group_name.a=b group_name.uid=123`) {
		t.Fatalf("unexpected log line")
	}
}
