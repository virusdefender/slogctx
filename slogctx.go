package slogctx

import (
	"context"
	"log/slog"
)

type Handler struct {
	parent  slog.Handler
	options *HandlerOptions
}

type HandlerOptions struct {
	// get these attrs from context and add to record
	AttrsFromCtx []string
	// add prefix to attrs to avoid conflict
	// for example, if AttrsFromCtx contains "trace_id", and AttrPrefix is "ctx_",
	// then the record will contain "ctx_trace_id" instead of "trace_id"
	AttrPrefix string
	// still show nil value in the log if some attr key is not found in the context
	ShowNilValue bool

	cachedAttrKey []string
}

func NewHandler(parent slog.Handler, options *HandlerOptions) slog.Handler {
	options.cachedAttrKey = make([]string, len(options.AttrsFromCtx))
	for idx, key := range options.AttrsFromCtx {
		options.cachedAttrKey[idx] = options.AttrPrefix + key
	}
	return &Handler{parent: parent, options: options}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.parent.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.parent.WithAttrs(attrs)
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h.parent.WithGroup(name)
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
