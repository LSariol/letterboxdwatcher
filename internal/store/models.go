package store

import "time"

type Subscription struct {
	ID                 int64      // Database ID
	Username           string     // Twitch display name
	UserId             string     // Twitch ID
	LetterboxdUsername string     //Letterboxd Username
	FeedURL            string     //RSS Feed URL
	LastSeenGUID       *string    //The last seen GUID from this profile. used to check for new movies
	LastFetchedAt      *time.Time //last time the feed was taken
	CreatedAt          time.Time  //The time the Subscription was created
	AlertChannels      []string
}
