package slogkit

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	var l slog.Level
	var log *slog.Logger

	var err = l.UnmarshalText([]byte(level))
	if err != nil {
		l = slog.LevelInfo
	}
	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	}))
	// Print warning later, so it's formated correctly
	if err != nil {
		log.Warn("Incorrect log level provided. Setting to default level", "level", slog.LevelInfo)
	}
	return log
}

func Error(err any) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Any("error", err)
}
