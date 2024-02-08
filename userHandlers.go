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

	ctx := r.Context()

	db := apiCfg.DB

	userPostReq := decodeUserPost{}

	decodeJson(r.Body, &userPostReq)

	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userPostReq.Name,
	}

	newUser, err := db.CreateUser(ctx, newUserParams)
	if err != nil {
		log.Fatalf("DB error on usersCreateHandler while trying to create user -> %v\n", err)
	}

	respondWithJSON(w, 200, newUser)
	log.Println("usersCreateHandler exited without any errors...")
}
