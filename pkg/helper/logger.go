package helper

import "github.com/rs/zerolog"

func LoggerWithNodeClassContext(logger zerolog.Logger, class string) zerolog.Logger {
	return logger.With().Str("node-class", class).Logger()
}
