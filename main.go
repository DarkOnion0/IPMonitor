package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/robfig/cron/v3"

	"github.com/DarkOnion0/IpMonitor/config"
)

var (
	currentIP  string
	previousIP string
)

// This is the response type for every request sent to the / endpoint when the API is enabled
type JSONRespT struct {
	CurrentIP  string
	PreviousIP string
	IpChanged  bool
}

func init() {
	// enable or not the debug level (default is Info)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *config.Debug == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()
	// activate the pretty logger for dev purpose only if the debug mode is enabled
	if *config.Debug == "true" {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().
		Str("type", "main").
		Str("function", "init").
		Msg("Logger is configured!")

	log.Debug().
		Str("type", "main").
		Str("function", "init").
		Msg("Debug mode is enabled!")
}

func init() {
	log.Debug().
		Str("Debug", *config.Debug).
		Str("Cron", *config.Cron).
		Str("EnableCron", *config.EnableCron).
		Str("APIPort", *config.APIPort).
		Str("EnableAPI", *config.EnableAPI).
		Msg("Printing the default settings")
}

// This function check if the public server IP has changed since the last lookup and return a warn log if it's the case
func ipChecker() (err error, ipChanged bool) {
	funcLog := log.With().
		Str("function", "IpChecker").
		Logger()

	funcLog.Debug().
		Msg("Start the function")

	funcLog.Debug().
		Msg("Send a GET request to the Cloudflare API")
	resp, err1 := http.Get("https://cloudflare.com/cdn-cgi/trace")

	if err1 != nil {
		funcLog.Error().
			Err(err1).
			Msg("Something bad append while getting the server if from Cloudlfare server, is your internet down ?")
		return errors.New("Something bad append while getting the server if from Cloudlfare server"), ipChanged
	}

	defer resp.Body.Close()

	funcLog.Debug().
		Msg("Start reading the response body")
	requestBody, err2 := io.ReadAll(resp.Body)

	if err2 != nil {
		funcLog.Error().
			Err(err2).
			Msg("Something bad append while reading the response body")

		return errors.New("Something bad append while reading the response body"), ipChanged
	}

	currentIP = strings.Split(strings.Split(string(requestBody), "\n")[2], "=")[1]

	switch previousIP {
	case currentIP:
		funcLog.Info().
			Str("currentIP", currentIP).
			Str("previousIP", previousIP).
			Bool("ipChanged", false).
			Msg("The ip has not changed since the last check: function finished successfully")

		return err, false
	case "":
		funcLog.Info().
			Str("currentIP", currentIP).
			Str("previousIP", previousIP).
			Bool("ipChanged", false).
			Msg("This is the first run of the function, previousIP is not set: function finished successfully")
		previousIP = currentIP
		currentIP = ""

		return err, false

	default:
		funcLog.Warn().
			Str("currentIP", currentIP).
			Str("previousIP", previousIP).
			Bool("ipChanged", true).
			Msg("The ip has changed since last check: function finished successfully")

		previousIP = currentIP
		currentIP = ""

		return err, true
	}
}

func main() {
	funcLog := log.With().
		Str("function", "Main").
		Logger()

	funcLog.Info().
		Msg("Start the app")

	if *config.EnableCron == "true" {
		c := cron.New()

		// set a cron job to check the public IP every 10 minutes
		// nolint
		c.AddFunc(*config.Cron, func() {
			err, _ := ipChecker()

			if err != nil {
				funcLog.Error().
					Str("type", "cron").
					Err(err).
					Msg("Something bad append while running the ipChecker function")
			}
		})

		// start all the cron jobs
		funcLog.Debug().
			Str("type", "cron").
			Msg("Start the cron jobs")
		c.Start()

		defer func(c *cron.Cron) {
			funcLog.Debug().
				Str("type", "cron").
				Msg("Closing cron jobs")
			c.Stop()
		}(c)
	}

	if *config.EnableAPI == "true" {
		app := fiber.New()

		app.Get("/", func(c *fiber.Ctx) error {
			funcLog.Debug().
				Str("type", "API").
				Str("endpoint", "/").
				Msg("Endpoint has been triggered")
			err, ipChanged := ipChecker()

			if err != nil {
				funcLog.Error().
					Str("type", "API").
					Err(err).
					Msg("Something bad append while running the ipChecker function")

				return err
			}

			resp := JSONRespT{
				CurrentIP:  currentIP,
				PreviousIP: previousIP,
				IpChanged:  ipChanged,
			}

			return c.JSON(resp)
		})

		funcLog.Debug().
			Str("type", "API").
			Msg("Starting fiber API server")

		err := app.Listen(":" + *config.APIPort)

		if err != nil {
			funcLog.Error().
				Str("type", "API").
				Err(err).
				Msg("Something bad append while starting the fiber API")

			return
		}
	}

	err, _ := ipChecker()

	if err != nil {
		funcLog.Error().
			Err(err).
			Msg("Something bad append while running the ipChecker function")

		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	funcLog.Info().
		Msg("Closing the app")
}
