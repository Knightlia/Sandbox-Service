package config

import (
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitLogging() {
	writers := io.MultiWriter(
		zerolog.ConsoleWriter{Out: os.Stdout},
	)

	// Display on filename + line number for caller field
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = zerolog.New(writers).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Info().Msg("Logger initialised.")

	// Debug level if debug flag set
	if viper.GetBool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Logger debug level enabled.")
	}
}
