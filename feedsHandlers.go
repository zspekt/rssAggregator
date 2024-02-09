package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/zspekt/rssAggregator/internal/database"
)

func feedsGetAllHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING feedsGetAllHandler...")

	var (
		db        *database.Queries = apiCfg.DB
		feedSlice []database.Feed   = make([]database.Feed, 0)
	)

	feedSlice, err := db.GetAllFeeds(r.Context())
	if err != nil {
		log.Printf(
			"Error getting all feeds from database during feedsGetAllHandler func --> %v\n",
			err,
		)
	}

	respondWithJSON(w, 200, feedSlice)
	log.Println("feedsGetAllHandler exited without any errors...")
}

func feedsCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING feedsCreateHandler...")

	var (
		// hold the database connection
		db = apiCfg.DB
		// struct to decode the http request's body
		feedPostReq = decodeFeedPost{}
	)

	apiKey, err := GetApiKeyFromHeader(r)
	if err != nil {
		log.Println("Error getting token from header on feedsCreateHandler func -> %v\n", err)
		respondWithError(w, 400, err.Error())
		return
	}

	decodeJson(r.Body, &feedPostReq)

	log.Printf("name field  <%v>\turl field  <%v>\n", feedPostReq.Name, feedPostReq.Url)

	if feedPostReq.Name == "" || feedPostReq.Url == "" {
		respondWithError(w, 400, "missing name or url field")
		log.Println("Request body was mising name or url fields...")
		return
	}

	userUUID, err := db.GetIdByApiKey(r.Context(), apiKey)
	if err != nil {
		log.Println("Error getting userUUID from DB on feedsCreateHandler func -> %v\n", err)
		respondWithError(w, 500, err.Error())
	}

	// params that will be used to run the CreateUser query
	newFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedPostReq.Name,
		Url:       feedPostReq.Url,
		UserID:    userUUID,
	}

	// creates users and returns so we can respond to the http.req with it
	newFeed, err := db.CreateFeed(r.Context(), newFeedParams)
	if err != nil {
		log.Fatalf("DB error on feedsCreateHandler while trying to create user -> %v\n", err)
	}

	respondWithJSON(w, 200, newFeed)
	log.Println("feedsCreateHandler exited without any errors...")
}
