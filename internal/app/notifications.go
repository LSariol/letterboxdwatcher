package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/lsariol/letterboxdwatcher/internal/common"
)

func (c *Client) BuildNotifications(accountActivity []common.FeedUpdate) []Notification {
	notifications := []Notification{}

	for _, account := range accountActivity {
		notification := Notification{}

		header := ""
		reviews := ""
		footer := ""

		if len(account.NewMovies) == 1 {
			movie := account.NewMovies[0]
			if movie.Rewatch {
				header = fmt.Sprintf("🎬 %s rewatched ", account.Subscription.Username)
			} else {
				header = fmt.Sprintf("🎬 %s watched ", account.Subscription.Username)
			}

			if movie.Liked {
				reviews = fmt.Sprintf("%s (%d) on %s | %s |❤️|", movie.FilmTitle, movie.FilmYear, movie.WatchedDate, getStars(movie.Rating))
			} else {
				reviews = fmt.Sprintf("%s (%d) on %s | %s |", movie.FilmTitle, movie.FilmYear, movie.WatchedDate, getStars(movie.Rating))
			}

			footer = fmt.Sprintf("Read their review here: %s", movie.Link)

		} else {
			header = fmt.Sprintf("🎬 %s logged %d new movies. ", account.Subscription.Username, len(account.NewMovies))
			footer = fmt.Sprintf("See other entries and reviews here: %sactivity", account.Subscription.FeedURL)

			reviewsLength := 480 - (len(header) + len(footer))

			var reviewList = []string{}
			currentLength := 2

			for _, movie := range account.NewMovies {
				review := ""

				if movie.Liked && movie.Rewatch {
					review = fmt.Sprintf("%s (%d)🔁 - %s ❤️", movie.FilmTitle, movie.FilmYear, getStars(movie.Rating))
				} else if movie.Liked && !movie.Rewatch {
					review = fmt.Sprintf("%s (%d) - %s ❤️", movie.FilmTitle, movie.FilmYear, getStars(movie.Rating))
				} else if !movie.Liked && movie.Rewatch {
					review = fmt.Sprintf("%s (%d)🔁 - %s", movie.FilmTitle, movie.FilmYear, getStars(movie.Rating))
				} else {
					review = fmt.Sprintf("%s (%d) - %s", movie.FilmTitle, movie.FilmYear, getStars(movie.Rating))
				}

				addedLength := len(review)
				if len(reviews) > 0 {
					addedLength += len(", ")
				}

				if currentLength+addedLength > reviewsLength {
					break
				}

				reviewList = append(reviewList, review)
				currentLength += addedLength
			}
			reviews = strings.Join(reviewList, ", ")
			reviews += ". "

		}

		notification.AlertChannels = account.Subscription.AlertChannels
		notification.GUID = *account.Subscription.LastSeenGUID
		notification.UserId = account.Subscription.UserId

		message := header + reviews + footer
		notification.Message = message

		notifications = append(notifications, notification)

	}

	return notifications
}

func getStars(stars float64) string {
	s := strings.Repeat("★", int(stars))

	if stars != math.Trunc(stars) {
		s += "½"
	}
	return s
}

func (c *Client) SendNotifications(ctx context.Context, notification Notification, url string) ([]byte, error) {

	body, err := json.Marshal(notification)
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("api returned status %d: %s", res.StatusCode, string(resBody))
	}

	return resBody, nil
}
