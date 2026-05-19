package parser

import (
	"github.com/lsariol/letterboxdwatcher/internal/common"
)

// Takes in a single users raw RSSFeed, and returns a clean parsedFeed
func ParseFeed(RSSFeed common.RawRSSFeed) []common.ParsedMovie {
	parsedFeed := []common.ParsedMovie{}

	for _, item := range RSSFeed.Channel.Items {

		var rewatch bool
		if item.Rewatch == "Yes" {
			rewatch = true
		} else {
			rewatch = false
		}

		var liked bool
		if item.MemberLike == "Yes" {
			liked = true
		} else {
			liked = false
		}

		parsedItem := common.ParsedMovie{
			Title:       item.Title,
			Link:        item.Link,
			Guid:        item.GUID.Value,
			PubDate:     item.PubDate,
			WatchedDate: item.WatchedDate,
			Rewatch:     rewatch,
			FilmTitle:   item.FilmTitle,
			FilmYear:    item.FilmYear,
			Rating:      item.MemberRating,
			Liked:       liked,
			MovieId:     item.MovieId,
			Review:      item.Description,
		}

		parsedFeed = append(parsedFeed, parsedItem)
	}

	return parsedFeed
}
