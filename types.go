package main

import "github.com/zspekt/rssAggregator/internal/database"

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
