package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/zspekt/rssAggregator/internal/database"
)

var (
	port   string
	dbURL  string
	db     *sql.DB
	apiCfg *apiConfig = &apiConfig{}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env --> %v\n", err)
		return
	}

	port = os.Getenv("PORT")
	log.Println("port var has been set...")

	dbURL := os.Getenv("dbConn")
	log.Println("dbURL var has been set...")

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database -> %v\n", err)
	}

	apiCfg.DB = database.New(db)

	log.Println(
		"The database connection has been stored in the apiCfg...",
	)

	log.Println("init func exited without errors...")
}

func main() {
	fmt.Print("\n\n\n")
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
		MaxAge:         300,
	}))

	v1 := chi.NewRouter()
	router.Mount("/api/v1", v1)

	v1.HandleFunc("/readiness", readinessHandler)
	v1.HandleFunc("/err", errorHandler)

	v1.Post("/users", usersCreateHandler)
	v1.Get("/users", usersGetByApiKey)

	v1.Post("/rssadd", feedsCreateHandler)
	v1.Get("/rssfeeds", feedsGetAllHandler)

	v1.Post("/feed_follows", feedFlPostHandler)

	v1.Delete("/feed_follows/*", feedFlDeleteHandler)

	v1.Get("/feed_follows", feedsGetByUserHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving on port: %v...\n\n\n\n", port)
	log.Fatal(srv.ListenAndServe())
}
