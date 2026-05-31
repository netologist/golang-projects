package main

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestContextHandler_enrichesRecord(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(slog.NewJSONHandler(&buf, nil))

	ctx := WithAttrs(
		context.Background(),
		slog.String("request_id", "req-7"),
		slog.String("tenant", "acme"),
	)
	logger.InfoContext(ctx, "handling request")

	out := buf.String()
	if !strings.Contains(out, "req-7") {
		t.Errorf("expected request_id in log, got: %s", out)
	}
	if !strings.Contains(out, "acme") {
		t.Errorf("expected tenant in log, got: %s", out)
	}
}

func TestWithAttrs_accumulates(t *testing.T) {
	ctx := WithAttrs(context.Background(), slog.String("a", "1"))
	ctx = WithAttrs(ctx, slog.String("b", "2"))

	var buf bytes.Buffer
	logger := NewLogger(slog.NewJSONHandler(&buf, nil))
	logger.InfoContext(ctx, "msg")

	out := buf.String()
	if !strings.Contains(out, `"a":"1"`) || !strings.Contains(out, `"b":"2"`) {
		t.Errorf("expected both attrs, got: %s", out)
	}
}
