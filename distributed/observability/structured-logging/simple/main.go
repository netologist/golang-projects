package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info(
		"user login",
		slog.String("user", "alice"),
		slog.Int("attempt", 1),
		slog.Bool("success", true),
	)
}
