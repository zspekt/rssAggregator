package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/zspekt/rssAggregator/internal/database"
)

func feedFlPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING feedFlPostHandler...")

	var (
		// hold the database connection
		db = apiCfg.DB
		// struct to decode the http request's body
		feedFlPostReq = decodeFeedFlPost{}
		apiKey        string
	)

	apiKey, err := GetApiKeyFromHeader(r)
	if err != nil {
		log.Println("Error getting token from header on feedFlPostHandler func -> %v\n", err)
		respondWithError(w, 400, err.Error())
		return
	}

	err = decodeJson(r.Body, &feedFlPostReq)
	if err != nil {
		log.Printf(
			"Error decoding json in feedFlPostHandler func --> %v\n",
			err,
		)
	}

	userID, err := db.GetIdByApiKey(r.Context(), apiKey)
	if err != nil {
		log.Printf(
			"Error decoding json in feedFlPostHandler func --> %v\n",
			err,
		)
	}

	arg := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feedFlPostReq.FeedId,
		UserID:    userID,
	}

	feedFollow, err := db.CreateFeedFollow(r.Context(), arg)
	if err != nil {
		log.Printf(
			"Error creating feed follow in feedFlPostHandler func --> %v\n",
			err,
		)
	}

	respondWithJSON(w, 200, feedFollow)
	log.Println("feedFlPostHandler exited without any errors...")
}

func feedFlDeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING feedFlDeleteHandler...")

	// hold the database connection
	db := apiCfg.DB

	feedFlId := chi.URLParam(r, "*")
	if feedFlId == "" {
		log.Printf(
			"Error getting feedFlId from URL in feedFlDeleteHandler func")
		respondWithError(w, 400, "missing url param")
		return
	}

	uuid := uuid.MustParse(feedFlId)

	log.Printf("feedFlId <%v>\n", uuid)

	err := db.DeleteFeedFollow(r.Context(), uuid)
	if err != nil {
		log.Printf(
			"Error deleting feed follow in feedFlDeleteHandler func --> %v\n",
			err,
		)
	}

	respondWithJSON(w, 200, "")
}

// write something to get param id from the url for chi
// id := chi.URLParam(r, "id")
// log.Println("id: ", id)
