package common

import (
	"github.com/lsariol/letterboxdwatcher/internal/store"
)

type RawRSSFeed struct {
	Channel struct {
		Items []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	GUID    GUID   `xml:"guid"`
	PubDate string `xml:"pubDate"`

	WatchedDate  string  `xml:"https://letterboxd.com watchedDate"`
	Rewatch      string  `xml:"https://letterboxd.com rewatch"`
	FilmTitle    string  `xml:"https://letterboxd.com filmTitle"`
	FilmYear     int     `xml:"https://letterboxd.com filmYear"`
	MemberRating float64 `xml:"https://letterboxd.com memberRating"`
	MemberLike   string  `xml:"https://letterboxd.com memberLike"`

	MovieId int `xml:"https://themoviedb.org movieId"`

	Description string `xml:"description"`
	Creator     string `xml:"http://purl.org/dc/elements/1.1/ creator"`
}

type GUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

// Account information gathered from internal DB and RSSFeed
type FeedData struct {
	Subscription store.Subscription
	Movies       []ParsedMovie
}

type ParsedMovie struct {
	Title       string
	Link        string
	Guid        string
	PubDate     string
	WatchedDate string
	Rewatch     bool
	FilmTitle   string
	FilmYear    int
	Rating      float64
	Liked       bool
	MovieId     int
	Review      string
}

// List of new entries not seen by this account before
type FeedUpdate struct {
	Subscription store.Subscription
	NewMovies    []ParsedMovie
}
