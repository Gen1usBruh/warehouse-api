package sl

import (
	"log/slog"
	"os"
	"github.com/Gen1usBruh/warehouse-api/internal/config"
)

func SetupLogger(conf *config.Logger) *slog.Logger {
	var log *slog.Logger
	if conf.Level == string("info") {
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	// not the right way to parse from Config, but now is ok

	return log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}