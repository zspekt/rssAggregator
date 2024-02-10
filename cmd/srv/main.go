package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	dbURL = os.Getenv("dbConn")
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
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
		MaxAge:         300,
	}))
	router.Use(secMiddleware)

	// endpoints that don't require auth
	v1 := chi.NewRouter()
	router.Mount("/api/v1", v1)

	// endpoinds that do
	v1Authd := chi.NewRouter()
	v1Authd.Use(authMiddleware)
	v1.Mount("/", v1Authd)

	v1.HandleFunc("/readiness", readinessHandler)
	v1.HandleFunc("/err", errorHandler)

	// users
	v1.Post("/users", usersCreateHandler)
	v1Authd.Get("/users", usersGetByApiKey)

	// feeds
	v1Authd.Post("/rssadd", feedsCreateHandler)
	v1.Get("/rssfeeds", feedsGetAllHandler)

	// feed follows
	v1Authd.Post("/feed_follows", feedFlPostHandler)
	v1Authd.Delete("/feed_follows/*", feedFlDeleteHandler)
	v1Authd.Get("/feed_follows", feedsGetByUserHandler)

	// posts
	v1Authd.Get("/posts", getPostsByUser)

	ctx := context.Background()
	go endlessFetching(ctx)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving on port: %v...\n\n\n", port)

	c := make(chan os.Signal, 1)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		<-c
		log.Println("Received interrupt signal. Shutting down gracefully...")
		srvErr := srv.Shutdown(shutdownCtx)
		if srvErr != nil {
			log.Printf("Error shutting down server -> %v\n", srvErr)
		}
	}()
	log.Fatal(srv.ListenAndServe())
}
