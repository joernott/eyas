package server

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
)

var Static embed.FS
var KeyPath string

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	i, err := Static.ReadFile("index.html")
	if err != nil {
		errorOut(w, 500, "ERR00004", "Could read index.html", err)
		return
	}
	fmt.Fprint(w, string(i[:]))
}

func Server(SSL bool, Port int, LogFile string, LogLevel string, Keydir string) {
	initLogging(LogFile, LogLevel)
	if _, err := os.Stat(Keydir); errors.Is(err, os.ErrNotExist) {
		log.Fatal().Err(err).Str("id", "ERR00003").Str("parameter", "keydir").Str("got", Keydir)
	}
	KeyPath = Keydir

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/index.html", Index)
	router.GET("/api/keys", PKCS7Keys)
	router.POST("/api/encrypt/single", EncryptSingle)
	router.POST("/api/encrypt/yaml", EncryptYaml)
	router.POST("/api/encrypt/csv", EncryptCSV)
	router.ServeFiles("/static/*filepath", http.FS(Static))
	log.Info().Bool("ssl", SSL).Int("port", Port).Str("logfile", LogFile).Str("loglevel", LogLevel).Msg("Starting server")
	if SSL {
		if err := http.ListenAndServeTLS(":"+strconv.Itoa(Port), "ssl/server.crt", "ssl/server.key", router); err != nil {
			log.Fatal().Err(err).Msg("Startup failed")
		}
	} else {
		if err := http.ListenAndServe(":"+strconv.Itoa(Port), router); err != nil {
			log.Fatal().Err(err).Msg("Startup failed")
		}
	}
}

func initLogging(LogFile string, LogLevel string) {
	switch strings.ToLower(LogLevel) {
	case "trace":
		{
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		}
	case "debug":
		{
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	case "info":
		{
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	case "warn":
		{
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		}
	case "error":
		{
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}
	case "fatal":
		{
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		}
	case "panic":
		{
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		}
	default:
		{
			expected := zerolog.Arr().
				Str("trace").
				Str("debug").
				Str("info").
				Str("warn").
				Str("error").
				Str("fatal").
				Str("panic")
			log.Error().
				Str("error", "Wrong Parameter").
				Str("id", "ERR00001").
				Str("parameter", "loglevel").
				Array("expected", expected).
				Str("got", LogLevel).
				Msg("Illegal log level " + LogLevel)
			fmt.Println("Illegal log level " + LogLevel)
			os.Exit(1)
		}
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if LogFile == "" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		f, err := os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND, os.FileMode(0640))
		if err != nil {
			log.Error().
				Err(err).
				Str("id", "ERR00002").
				Str("parameter", "logfile").
				Str("got", LogFile).
				Msg("Can't open log file  " + LogFile)
			fmt.Println("Can't open log file " + LogFile)
			os.Exit(1)
		}
		log.Logger = zerolog.New(f)
	}

}
