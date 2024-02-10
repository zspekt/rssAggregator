package main

import (
	"fmt"
	"log"
	"net/http"

	jsongenerics "github.com/zspekt/rssAggregator/internal/jsonGenerics"
)

func getPostsByUser(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING getPostsByUser...")

	db := apiCfg.DB
	apiKey, err := GetApiKeyFromHeader(r)
	if err != nil {
		log.Printf(
			"Error getting token from header in getPostsByUser func --> %v\n",
			err,
		)
		jsongenerics.RespondWithError(w, 400, err.Error())
		return
	}

	userID, err := db.GetIdByApiKey(r.Context(), apiKey)
	if err != nil {
		log.Printf(
			"Error getting user id from api key in getPostsByUser func --> %v\n",
			err,
		)
		jsongenerics.RespondWithError(w, 400, err.Error())
		return
	}

	postsSlice, err := db.GetPostsByUser(r.Context(), userID)
	if err != nil {
		log.Printf(
			"Error getting posts by user id in getPostsByUser func --> %v\n",
			err,
		)
		jsongenerics.RespondWithError(w, 400, err.Error())
		return
	}

	// doing this bc the db method returns its own type of posts, and i cannot
	// add any tags to em.

	var jsonPostSlice []getPostsByUserResp
	for _, post := range postsSlice {
		jsonPostSlice = append(jsonPostSlice, getPostsByUserResp{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description.String,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		})
	}
	jsongenerics.RespondWithJSON(w, 200, jsonPostSlice)
}
