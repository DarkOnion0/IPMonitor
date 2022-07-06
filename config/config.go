package config

import (
	"flag"
)

func init() {
	flag.Parse()
}

var (
	Debug      = flag.String("debug", "false", "Sets log level to debug")
	Cron       = flag.String("cron", "*/15 * * * *", "Set a custom cron scheduled to run the IP check, run every 15 minutes by default")
	EnableCron = flag.String("cron-enable", "true", "This flag enable the cron mode, it can be disable to it in API mode only")
	EnableAPI  = flag.String("api-enable", "true", "This flag enable the API mode, it can be disable to run it in cron mode only")
	APIPort    = flag.String("api-port", "8080", "Set a custom api's listen port")
)
