package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zspekt/rssAggregator/internal/database"
	jsongenerics "github.com/zspekt/rssAggregator/internal/jsonGenerics"
)

func getPostsByUser(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING getPostsByUser...")

	var (
		db   database.Queries = *apiCfg.DB
		user database.User
	)

	userPtr, ok := r.Context().Value("user").(*database.User)
	if !ok {
		jsongenerics.RespondWithError(w, 400, "Unauthorized access")
		return
	}
	user = *userPtr

	postsSlice, err := db.GetPostsByUser(r.Context(), user.ID)
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
