package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	logger := NewLogger(slog.NewJSONHandler(os.Stdout, nil))

	ctx := WithAttrs(
		context.Background(),
		slog.String("request_id", "req-42"),
		slog.String("user", "alice"),
	)
	logger.InfoContext(ctx, "order placed", slog.Int("order_id", 100))
	logger.WarnContext(ctx, "inventory low", slog.String("sku", "ABC"))
}
