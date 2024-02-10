package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/zspekt/rssAggregator/internal/database"
	jsongenerics "github.com/zspekt/rssAggregator/internal/jsonGenerics"
)

func feedFlPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING feedFlPostHandler...")

	var (
		// hold the database connection
		db = apiCfg.DB
		// struct to decode the http request's body
		feedFlPostReq = decodeFeedFlPost{}
		user          database.User
	)

	userPtr, ok := r.Context().Value("user").(*database.User)
	if !ok {
		jsongenerics.RespondWithError(w, 400, "Unauthorized access")
		return
	}
	user = *userPtr

	err := jsongenerics.DecodeJson(r.Body, &feedFlPostReq)
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
		UserID:    user.ID,
	}

	feedFollow, err := db.CreateFeedFollow(r.Context(), arg)
	if err != nil {
		log.Printf(
			"Error creating feed follow in feedFlPostHandler func --> %v\n",
			err,
		)
	}

	jsongenerics.RespondWithJSON(w, 200, feedFollow)
	log.Println("feedFlPostHandler exited without any errors...")
}

func feedFlDeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING feedFlDeleteHandler...")

	var (
		db   *database.Queries = apiCfg.DB
		user database.User
	)

	userPtr, ok := r.Context().Value("user").(*database.User)
	if !ok {
		jsongenerics.RespondWithError(w, 400, "Unauthorized access")
		return
	}
	user = *userPtr

	feedFlId := chi.URLParam(r, "*")
	if feedFlId == "" {
		log.Printf(
			"Error getting feedFlId from URL in feedFlDeleteHandler func")
		jsongenerics.RespondWithError(w, 400, "missing url param")
		return
	}

	uuid := uuid.MustParse(feedFlId)

	log.Printf("feedFlId <%v>\n", uuid)

	args := database.DeleteFeedFollowParams{
		ID:     uuid,
		UserID: user.ID,
	}

	err := db.DeleteFeedFollow(r.Context(), args)
	if err != nil {
		log.Printf(
			"Error deleting feed follow in feedFlDeleteHandler func --> %v\n",
			err,
		)
	}

	jsongenerics.RespondWithJSON(w, 200, "")
}

func feedsGetByUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING feedsGetAllByUserHandler...")

	var (
		db   *database.Queries = apiCfg.DB
		user database.User
	)

	userPtr, ok := r.Context().Value("user").(*database.User)
	if !ok {
		jsongenerics.RespondWithError(w, 400, "Unauthorized access")
		return
	}
	user = *userPtr

	feedSlice, err := db.GetFeedFollowsByUser(r.Context(), user.ID)
	if err != nil {
		log.Printf(
			"Error getting feeds by user id in feedsGetAllByUserHandler func --> %v\n",
			err,
		)
		jsongenerics.RespondWithError(w, 400, err.Error())
		return
	}

	jsongenerics.RespondWithJSON(w, 200, feedSlice)
	log.Println("feedsGetAllByUserHandler exited without any errors...")
}
