package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/zspekt/rssAggregator/internal/database"
)

func usersCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING usersCreateHandler...")

	userPostReq := decodeUserPost{}

	decodeJson(r.Body, &userPostReq)

	newUser := database.User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userPostReq.Name,
	}

	respondWithJSON(w, 200, newUser)
}
