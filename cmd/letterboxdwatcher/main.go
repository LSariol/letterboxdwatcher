package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lsariol/letterboxdwatcher/internal/app"
	"github.com/lsariol/letterboxdwatcher/internal/common"
	"github.com/lsariol/letterboxdwatcher/internal/config"
	"github.com/lsariol/letterboxdwatcher/internal/parser"
	"github.com/lsariol/letterboxdwatcher/internal/store"
)

func main() {
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("FATAL: Failed to load EST timezone: %v", err)
	}
	time.Local = est

	cfg := config.Load()

	client, err := initialize(cfg)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize service: %v", err)
	}

	ctx := context.Background()
	notificationDelay := time.Duration(cfg.NotificationDelaySeconds) * time.Second
	pollWindow := time.Duration(cfg.PollIntervalMinutes) * time.Minute
	minFeedDelay := time.Duration(cfg.MinFeedDelaySeconds) * time.Second

	runLoop(ctx, client, cfg.NotificationEndpoint, notificationDelay, pollWindow, minFeedDelay)
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

func runLoop(ctx context.Context, client *app.Client, notificationEndpoint string, notificationDelay time.Duration, pollWindow time.Duration, minFeedDelay time.Duration) {
	for {
		subscriptions, err := client.Store.GetSubscriptions(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to load subscriptions: %v. Retrying after poll window.", err)
			time.Sleep(pollWindow)
			continue
		}

		n := len(subscriptions)
		if n == 0 {
			log.Println("INFO: No subscriptions found. Sleeping for poll window.")
			time.Sleep(pollWindow)
			continue
		}

		feedDelay := pollWindow / time.Duration(n)
		if feedDelay < minFeedDelay {
			feedDelay = minFeedDelay
		}

		log.Printf("INFO: Starting staggered poll cycle: %d subscriptions, %s between feeds.", n, feedDelay)

		for _, subscription := range subscriptions {
			processSubscription(ctx, client, subscription, notificationEndpoint, notificationDelay)
			time.Sleep(feedDelay)
		}

		log.Println("INFO: Staggered poll cycle complete.")
	}
}

func processSubscription(ctx context.Context, client *app.Client, subscription store.Subscription, notificationEndpoint string, notificationDelay time.Duration) {
	rawFeed, err := client.GetRawRSSFeed(ctx, subscription)
	if err != nil {
		log.Printf("WARN: Skipping subscription for %s: %v", subscription.Username, err)
		return
	}

	account := common.FeedData{
		Subscription: subscription,
		Movies:       parser.ParseFeed(rawFeed),
	}

	feedUpdates := client.GetNewFeedActivity([]common.FeedData{account})
	notifications := client.BuildNotifications(feedUpdates)

	if len(notifications) == 0 {
		log.Printf("INFO: No new activity for %s.", subscription.Username)
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
}
