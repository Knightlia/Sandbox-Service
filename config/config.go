package config

import (
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct{}

func NewConfig() Config {
	return Config{}
}

// InitFlags initialise CLI flags.
func (_ Config) InitFlags() {
	pflag.Uint("port", 8080, "The port to run the server on. Default is 8080.")
	pflag.String("config", ".", "The directory where the config.yaml exists. Defaults to current working directory.")
	pflag.Bool("debug", false, "Enable debug for development. False by default.")
	pflag.String("logs", "logs", "Directory for the logs. In 'logs' folder at the current working directory by default.")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}

// InitLogger initialises the rolling file logger and the logging format.
func (_ Config) InitLogger() {
	writers := io.MultiWriter(
		zerolog.ConsoleWriter{Out: os.Stdout},
		&lumberjack.Logger{
			Filename:   viper.GetString("logs") + "/sandbox.log",
			MaxSize:    100,
			MaxAge:     10,
			MaxBackups: 3,
			Compress:   true,
		},
	)

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

	log.Logger = zerolog.New(writers).With().Timestamp().Caller().Logger()
	log.Info().Msg("Logger initialised.")

	if viper.GetBool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Logger debug level enabled.")
	}
}

// InitConfigFile loads an external configuration file and registers it into [viper].
func (_ Config) InitConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(viper.GetString("config"))
	viper.AddConfigPath("$HOME/.config/sandbox-service") // call multiple times to add many search paths
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic().
			Err(err).
			Msg("Error reading configuration file.")
	}

	log.Info().Msg("Configuration file loaded.")
}
