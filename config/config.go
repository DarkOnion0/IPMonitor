package config

import (
	"flag"
)

func init() {
	flag.Parse()
}

var (
	Debug = flag.String("debug", "false", "Sets log level to debug")
	Cron  = flag.String("cron", "*/15 * * * *", "Set a custom cron scheduled to run the IP check, run every 15 minutes by default")
)
