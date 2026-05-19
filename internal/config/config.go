package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL                string
	NotificationEndpoint       string
	PollIntervalMinutes        int
	NotificationDelaySeconds   int
}

func Load() Config {
	interval, _ := strconv.Atoi(os.Getenv("POLL_INTERVAL_MINUTES"))
	if interval == 0 {
		interval = 25
	}
	delay, _ := strconv.Atoi(os.Getenv("NOTIFICATION_DELAY_SECONDS"))
	if delay == 0 {
		delay = 3
	}
	return Config{
		DatabaseURL:              os.Getenv("DATABASE_URL"),
		NotificationEndpoint:     os.Getenv("NOTIFICATION_ENDPOINT"),
		PollIntervalMinutes:      interval,
		NotificationDelaySeconds: delay,
	}
}
