package observability

import (
	"context"
	"log/slog"
)

type ContextHandler struct{ slog.Handler }
type ctxKey string

const requestIDKey ctxKey = "request_id"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h ContextHandler) WithAttrs(a []slog.Attr) slog.Handler {
	return ContextHandler{h.Handler.WithAttrs(a)}
}
func (h ContextHandler) WithGroup(n string) slog.Handler {
	return ContextHandler{h.Handler.WithGroup(n)}
}
