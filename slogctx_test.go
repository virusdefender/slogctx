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
	Userid    *int    `json:"uid"`
	RequestId *string `json:"request_id"`
	NotFound  *string `json:"not_found"`
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

type RequestId string
type UserId struct{}

func testContextValues(t *testing.T, showNilValue bool) {
	buf := &bytes.Buffer{}
	h := NewHandler(
		slog.NewJSONHandler(buf, &slog.HandlerOptions{AddSource: true}),
		&HandlerOptions{
			AttrsFromCtx: []any{"url", RequestId("request_id"), UserId{}, "not_found"},
			AttrsNameMap: map[any]string{UserId{}: "uid"},
			ShowNilValue: showNilValue,
		})
	logger := slog.New(h)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "url", "https://example.com/api")
	ctx = context.WithValue(ctx, RequestId("request_id"), "abcdef")
	ctx = context.WithValue(ctx, UserId{}, 123)

	logger.ErrorContext(ctx, "error msg with context")

	logString := buf.String()
	fmt.Println(logString)

	logItem := LogItem{}
	err := json.Unmarshal([]byte(logString), &logItem)
	if err != nil {
		t.Fatalf("unmarshal line %s failed: %s", logString, err)
	}
	if *logItem.Userid != 123 || *logItem.RequestId != "abcdef" || logItem.NotFound != nil {
		t.Fatalf("unexpected log line")
	}

	if showNilValue && !strings.Contains(logString, `"not_found":null`) {
		t.Fatalf("unexpected log line")
	} else if !showNilValue && strings.Contains(logString, `"not_found":null`) {
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
		&HandlerOptions{AttrsFromCtx: []any{"uid"}})

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
		&HandlerOptions{AttrsFromCtx: []any{"uid"}})

	logger := slog.New(h).WithGroup("group_name")

	logger.ErrorContext(context.WithValue(context.Background(), "uid", 123), "error msg", "a", "b")

	logString := buf.String()
	fmt.Println(logString)
	if !strings.Contains(logString, `group_name.a=b group_name.uid=123`) {
		t.Fatalf("unexpected log line")
	}
}
