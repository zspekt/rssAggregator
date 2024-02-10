package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/zspekt/rssAggregator/internal/database"
	jsongenerics "github.com/zspekt/rssAggregator/internal/jsonGenerics"
)

func usersCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING usersCreateHandler...")

	var (
		// hold the database connection
		db = apiCfg.DB
		// struct to decode the http request's body
		userPostReq = decodeUserPost{}
	)

	jsongenerics.DecodeJson(r.Body, &userPostReq)

	if userPostReq.Name == "" {
		jsongenerics.RespondWithError(w, 400, "missing name field")
	}

	// params that will be used to run the CreateUser query
	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userPostReq.Name,
	}

	// creates users and returns so we can respond to the http.req with it
	newUser, err := db.CreateUser(r.Context(), newUserParams)
	if err != nil {
		log.Fatalf("DB error on usersCreateHandler while trying to create user -> %v\n", err)
	}

	jsongenerics.RespondWithJSON(w, 200, newUser)
	log.Println("usersCreateHandler exited without any errors...")
}

func usersGetByApiKey(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING usersGetByApiKey...")

	// will be sent back to the client
	var userResp database.User

	userPtr, ok := r.Context().Value("user").(*database.User)
	if !ok {
		jsongenerics.RespondWithError(w, 400, "Unauthorized access")
		return
	}
	userResp = *userPtr

	jsongenerics.RespondWithJSON(w, 200, userResp)
}
