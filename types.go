package main

import (
	"github.com/google/uuid"

	"github.com/zspekt/rssAggregator/internal/database"
)

type readinessResp struct {
	Status string `json:"status"`
}

type errorResp struct {
	Error string `json:"error"`
}

type apiConfig struct {
	DB *database.Queries
}

type decodeUserPost struct {
	Name string `json:"name"`
}

type decodeUserGet struct {
	ApiKey string `json:"api_key"`
}

type decodeFeedPost struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// we get the user id from the apikey
type decodeFeedFlPost struct {
	FeedId uuid.UUID `json:"feed_id"`
}

type createFeedResp struct {
	Feed       database.Feed       `json:"feed"`
	FeedFollow database.Feedfollow `json:"feed_follow"`
}
