package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lsariol/letterboxdwatcher/internal/app"
	"github.com/lsariol/letterboxdwatcher/internal/config"
	"github.com/lsariol/letterboxdwatcher/internal/store"
)

func main() {
	cfg := config.Load()

	client, err := initialize(cfg)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize service: %v", err)
	}

	ctx := context.Background()

	notificationDelay := time.Duration(cfg.NotificationDelaySeconds) * time.Second
	runOnce(ctx, client, cfg.NotificationEndpoint, notificationDelay)
	ticker := time.NewTicker(time.Duration(cfg.PollIntervalMinutes) * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		runOnce(ctx, client, cfg.NotificationEndpoint, notificationDelay)
	}
}

func initialize(cfg config.Config) (*app.Client, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	s := store.NewStore(pool)

	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	return app.NewClient(s, httpClient), nil
}

func runOnce(ctx context.Context, client *app.Client, notificationEndpoint string, notificationDelay time.Duration) {
	log.Println("INFO: Starting feed poll...")

	accounts := client.GetAccounts(ctx)
	if accounts == nil {
		log.Println("WARN: Skipping poll cycle, failed to load accounts.")
		return
	}

	newAccountActivity := client.GetNewFeedActivity(accounts)
	notifications := client.BuildNotifications(newAccountActivity)

	if len(notifications) == 0 {
		log.Println("INFO: No new activity found.")
		return
	}

	for _, notification := range notifications {
		res, err := client.SendNotifications(ctx, notification, notificationEndpoint)
		if err != nil {
			log.Printf("ERROR: Failed to send notification for user %s: %v", notification.UserId, err)
			continue
		}

		log.Printf("INFO: Notification sent for user %s: %s", notification.UserId, res)
		time.Sleep(notificationDelay)

		if err := client.Store.UpdateLastSeenGUID(ctx, notification.UserId, notification.GUID); err != nil {
			log.Printf("ERROR: Failed to update last seen GUID for user %s: %v", notification.UserId, err)
		}
	}

	log.Println("INFO: Feed poll complete.")
}
