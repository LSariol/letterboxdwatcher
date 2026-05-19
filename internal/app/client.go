package app

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lsariol/letterboxdwatcher/internal/common"
	"github.com/lsariol/letterboxdwatcher/internal/parser"
	"github.com/lsariol/letterboxdwatcher/internal/store"
)

type Client struct {
	HttpClient *http.Client
	Store      *store.Store
}

func NewClient(store *store.Store, http *http.Client) *Client {
	c := Client{
		HttpClient: http,
		Store:      store,
	}

	return &c
}

func (c *Client) Notify() {

}

// Loads subscriptions from database, pulls RSS feeds and returns a slice of common.FeedData.
// A fatal DB failure returns nil so the caller can skip the poll cycle entirely.
// Individual feed fetch failures are logged and skipped — one bad subscription won't halt the rest.
func (c *Client) GetAccounts(ctx context.Context) []common.FeedData {

	subscriptions, err := c.Store.GetSubscriptions(ctx)
	if err != nil {
		log.Printf("ERROR: failed to get subscriptions: %v", err)
		return nil
	}

	if len(subscriptions) == 0 {
		log.Println("INFO: no subscriptions found, nothing to do.")
		return nil
	}

	accounts := make([]common.FeedData, 0, len(subscriptions))

	for _, subscription := range subscriptions {
		rawRSSFeed, err := c.GetRawRSSFeed(ctx, subscription)
		if err != nil {
			log.Printf("WARN: skipping subscription for %s: %v", subscription.Username, err)
			continue
		}

		accounts = append(accounts, common.FeedData{
			Subscription: subscription,
			Movies:       parser.ParseFeed(rawRSSFeed),
		})
	}

	return accounts
}

// Takes in a UserRecord, and returns the entire RSS Feed as a single object
func (c *Client) GetRawRSSFeed(ctx context.Context, subscription store.Subscription) (common.RawRSSFeed, error) {
	userFeed := common.RawRSSFeed{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, subscription.FeedURL, nil)
	if err != nil {
		return userFeed, fmt.Errorf("create RSS request: %w", err)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return userFeed, fmt.Errorf("pull RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return userFeed, fmt.Errorf("unexpected RSS response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return userFeed, fmt.Errorf("read RSS body: %w", err)
	}

	if err := xml.Unmarshal(body, &userFeed); err != nil {
		return userFeed, err
	}

	return userFeed, nil
}

func (c *Client) GetNewFeedActivity(accounts []common.FeedData) []common.FeedUpdate {
	feedUpdates := []common.FeedUpdate{}

	for _, account := range accounts {

		feedUpdate := common.FeedUpdate{}
		feedUpdate.Subscription = account.Subscription

		if account.Subscription.LastSeenGUID == nil {
			if len(account.Movies) == 0 {
				continue
			} else {

				feedUpdate.NewMovies = append(feedUpdate.NewMovies, account.Movies[0])
				feedUpdate.Subscription.LastSeenGUID = &feedUpdate.NewMovies[0].Guid
				feedUpdates = append(feedUpdates, feedUpdate)

				continue
			}
		} else {

			lastSeenGUID := account.Subscription.LastSeenGUID

			feedUpdate.Subscription.LastSeenGUID = &account.Movies[0].Guid

			for _, movie := range account.Movies {
				if *lastSeenGUID == movie.Guid {
					break
				} else {

					feedUpdate.NewMovies = append(feedUpdate.NewMovies, movie)

				}
			}

			if len(feedUpdate.NewMovies) > 0 {
				feedUpdates = append(feedUpdates, feedUpdate)
			}

			continue
		}
	}

	return feedUpdates
}
