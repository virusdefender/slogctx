package slogctx

import (
	"context"
	"fmt"
	"log/slog"
)

type Handler struct {
	parent  slog.Handler
	options *HandlerOptions
}

type HandlerOptions struct {
	// get these attrs from context and add to record
	AttrsFromCtx []any
	// to avoid conflict, context value key can be any type, so we need a map to convert it to string
	// it's not necessary to add it if the key can be converted to string directly
	AttrsNameMap map[any]string
	// still show nil value in the log if some attr key is not found in the context
	ShowNilValue bool

	cachedAttrKey []string
}

func NewHandler(parent slog.Handler, options *HandlerOptions) slog.Handler {
	options.cachedAttrKey = make([]string, len(options.AttrsFromCtx))
	for idx, key := range options.AttrsFromCtx {
		k, exists := options.AttrsNameMap[key]
		if !exists {
			k = fmt.Sprintf("%v", key)
		}
		options.cachedAttrKey[idx] = k
	}
	return &Handler{parent: parent, options: options}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.parent.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{parent: h.parent.WithAttrs(attrs), options: h.options}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{parent: h.parent.WithGroup(name), options: h.options}
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	attrs := make([]slog.Attr, 0, len(h.options.AttrsFromCtx))
	for idx, key := range h.options.AttrsFromCtx {
		ctxVal := ctx.Value(key)
		if h.options.ShowNilValue || ctxVal != nil {
			attrs = append(attrs, slog.Any(h.options.cachedAttrKey[idx], ctxVal))
		}
	}
	record.AddAttrs(attrs...)
	return h.parent.Handle(ctx, record)
}
