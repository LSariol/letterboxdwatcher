package store

import (
	"context"
	"fmt"
)

func (s *Store) CreateFeed(ctx context.Context, feed Subscription) (Subscription, error) {
	query := `
		INSERT INTO botsuite.letterboxd_feed_subscriptions (
			username,
			user_id,
			letterboxd_username,
			feed_url,
			alert_channels
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			id,
			username,
			user_id,
			letterboxd_username,
			feed_url,
			last_seen_guid,
			last_alerted_at,
			created_at;
	`
	var newSubscription Subscription

	err := s.Pool.QueryRow(
		ctx,
		query,
		feed.Username,
		feed.UserId,
		feed.LetterboxdUsername,
		feed.FeedURL,
		feed.AlertChannels,
	).Scan(
		&newSubscription.ID,
		&newSubscription.Username,
		&newSubscription.UserId,
		&newSubscription.LetterboxdUsername,
		&newSubscription.FeedURL,
		&newSubscription.LastSeenGUID,
		&newSubscription.LastFetchedAt,
		&newSubscription.CreatedAt,
		&newSubscription.AlertChannels,
	)

	if err != nil {
		return Subscription{}, err
	}

	return newSubscription, nil
}

func (s *Store) GetSubscriptions(ctx context.Context) ([]Subscription, error) {
	query := `
		SELECT
			id,
			username,
			user_id,
			letterboxd_username,
			feed_url,
			last_seen_guid,
			last_alerted_at,
			created_at,
			alert_channels
		FROM botsuite.letterboxd_feed_subscriptions
		ORDER BY created_at DESC;
	`
	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feeds []Subscription

	for rows.Next() {
		var feed Subscription

		err := rows.Scan(
			&feed.ID,
			&feed.Username,
			&feed.UserId,
			&feed.LetterboxdUsername,
			&feed.FeedURL,
			&feed.LastSeenGUID,
			&feed.LastFetchedAt,
			&feed.CreatedAt,
			&feed.AlertChannels,
		)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (s *Store) UpdateLastSeenGUID(ctx context.Context, UserId string, GUID string) error {
	query := `
	UPDATE botsuite.letterboxd_feed_subscriptions
	SET
		last_seen_guid = $1,
		last_alerted_at = NOW()
	WHERE user_id = $2;
	`

	tag, err := s.Pool.Exec(ctx, query, GUID, UserId)
	if err != nil {
		return fmt.Errorf("update last seen guid for user_id %s: %w", UserId, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no letterboxd feed subscription found for user_id %s", UserId)
	}

	return nil
}
