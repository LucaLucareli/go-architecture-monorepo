package logger

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	ColorReset  = "\x1b[0m"
	ColorRed    = "\x1b[31m"
	ColorGreen  = "\x1b[32m"
	ColorYellow = "\x1b[33m"
	ColorCyan   = "\x1b[36m"
	ColorGray   = "\x1b[90m"
	ColorWhite  = "\x1b[97m"
	ColorBold   = "\x1b[1m"
	ColorPurple = "\x1b[35m"
	ColorBlue   = "\x1b[34m"
)

var levelEmojis = map[string]string{
	"INFO":  "✨",
	"WARN":  "⚠️",
	"ERROR": "💥",
	"FATAL": "💀",
	"DEBUG": "🔍",
	"TRACE": "👣",
}

func Init(appName, appColor, env string) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339

	setLogLevel(env)

	var output io.Writer
	if strings.ToUpper(env) == "DEV" {
		output = setupConsoleWriter(appName, appColor)
	} else {
		output = os.Stdout
	}

	l := zerolog.New(output).
		With().
		Timestamp().
		Logger()

	log.Logger = l

	stdlog.SetFlags(0)
	stdlog.SetOutput(logWriter{})
}

func setLogLevel(env string) {
	switch strings.ToUpper(env) {
	case "PROD":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func setupConsoleWriter(appName, appColor string) zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05.000",
		NoColor:    false,
		FormatLevel: func(i interface{}) string {
			level := strings.ToUpper(fmt.Sprintf("%s", i))
			emoji := levelEmojis[level]
			if emoji == "" {
				emoji = "•"
			}

			levelColor := getLevelColor(level)

			return fmt.Sprintf(
				"%s[%s - %d]%s  %s%s%-5s%s %s ",
				appColor, appName, os.Getpid(), ColorReset,
				levelColor, ColorBold, level, ColorReset,
				emoji,
			)
		},
		FormatMessage: func(i interface{}) string {
			if i == nil {
				return ""
			}

			var msg string
			switch v := i.(type) {
			case string:
				msg = v
			case fmt.Stringer:
				msg = v.String()
			default:
				msg = fmt.Sprintf("%v", v)
			}

			return fmt.Sprintf("%s%s%s%s", ColorBold, ColorWhite, msg, ColorReset)
		},

		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s%v:%s", ColorGray, i, ColorReset)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s%v%s", ColorCyan, i, ColorReset)
		},
	}
}

func getLevelColor(level string) string {
	switch level {
	case "INFO":
		return ColorGreen
	case "WARN":
		return ColorYellow
	case "ERROR", "FATAL":
		return ColorRed
	case "DEBUG", "TRACE":
		return ColorCyan
	default:
		return ColorReset
	}
}

type logWriter struct{}

func (l logWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSuffix(string(p), "\n")
	if msg == "<nil>" {
		msg = ""
	}
	log.Info().Msg(msg)
	return len(p), nil
}

func PrintRoutes(e *echo.Echo) {
	fmt.Println()
	log.Info().Msg("╔═════════════════════════════════════════════════════╗")
	log.Info().Msg("║             API ROUTES INITIALIZED                  ║")
	log.Info().Msg("╚═════════════════════════════════════════════════════╝")

	for _, r := range e.Routes() {
		if r.Path == "/*" || r.Path == "" || r.Method == "ECHO" {
			continue
		}

		handlerName := formatHandlerName(r.Name)

		log.Debug().
			Str("method", fmt.Sprintf("%-s", r.Method)).
			Str("path", fmt.Sprintf("%-s", r.Path)).
			Str("handler", handlerName).
			Msg("🔍 Mapped")
	}
	fmt.Println()
}

func formatHandlerName(name string) string {
	parts := strings.Split(name, "/")
	name = parts[len(parts)-1]

	name = strings.TrimSuffix(name, "-fm")
	name = strings.ReplaceAll(name, "(*", "")
	name = strings.ReplaceAll(name, ")", "")

	subParts := strings.Split(name, ".")
	if len(subParts) >= 2 {
		return subParts[len(subParts)-2]
	}

	return name
}
