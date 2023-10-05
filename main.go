package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	tgclient "read-adviser-bot/clients/telegram"
	event_consumer "read-adviser-bot/consumer/event-consumer"
	"read-adviser-bot/events/telegram"
	"read-adviser-bot/lib/logger/handlers/slogpretty"
	"read-adviser-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

var logLevelStr, token *string

func init() {
	token = flag.String(
		"tg-bot-token",
		"",
		"token for access to tg bot",
	)

	logLevelStr = flag.String(
		"log-level",
		"info",
		"logging level (debug, info, warn, error)",
	)

	flag.Parse()
}

func main() {
	logger := setupLogger()

	logger.Debug("creating event processor")

	eventsProcessor := telegram.New(
		tgclient.New(tgBotHost, mustToken(logger), logger),
		files.New(storagePath),
		logger,
	)

	logger.Debug("creating consumer")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize, logger)

	logger.Debug("starting consumer")

	if err := consumer.Start(); err != nil {
		logger.Error("Service is stopped", err)
		os.Exit(1)
	}
}

func mustToken(logger *slog.Logger) string {
	if *token == "" {
		logger.Error("token is not specified")
		os.Exit(1)
	}

	return *token
}

func setupLogger() *slog.Logger {

	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: logLevel(),
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func logLevel() slog.Level {
	var logLevel slog.Level

	switch *logLevelStr {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		log.Fatalf("Unsupported log level: %s", *logLevelStr)
	}

	return logLevel
}
