package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/zspekt/rssAggregator/internal/database"
	"github.com/zspekt/rssAggregator/internal/xmldecoding"
)

func GetApiKeyFromHeader(r *http.Request) (string, error) {
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("Authorization header is missing")
	}

	log.Printf("Retrieved token from header -> %v\n", apiKey)
	return apiKey, nil
}

func fetchFeed(url string) (io.Reader, error) {
	fmt.Print("\n\n\n")
	log.Printf("Fetching feed -> %v\n", url)

	httpResp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making Get request -> %v\n", err)
		return nil, err
	}

	return httpResp.Body, nil
}

func getRss(url string, rss *xmldecoding.Rss) error {
	fmt.Print("\n\n\n")
	log.Printf("Getting rss...")

	r, err := fetchFeed(url)
	if err != nil {
		log.Printf("Error fetching feed -> %v\n", err)
		return err
	}

	err = xmldecoding.DecodeXml(r, rss)
	if err != nil {
		log.Printf("Error decoding xml -> %v\n", err)
		return err
	}

	return nil
}

func processRss(feed database.Feed, wg *sync.WaitGroup) error {
	var (
		rssStruct xmldecoding.Rss   = xmldecoding.Rss{}
		db        *database.Queries = apiCfg.DB
	)

	defer wg.Done()

	r, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed -> %v\n", err)
		return err
	}

	err = xmldecoding.DecodeXml(r, &rssStruct)
	if err != nil {
		log.Printf("Error decoding xml -> %v\n", err)
		return err
	}

	time := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	arg := database.MarkFeedFetchedParams{
		LastFetchedAt: time,
		ID:            feed.ID,
	}

	db.MarkFeedFetched(context.Background(), arg)
	log.Printf("RSS <%v>\n", rssStruct.Channel.Title)
	log.Println("Go routine finished...")
	return nil
}

func endlessFetching(context context.Context) {
	fmt.Print("\n\n\n")
	log.Printf("Running endless fetching...")

	const interval = 20 * time.Second

	var (
		db                      = apiCfg.DB
		feedAmount int32        = 10
		ticker     *time.Ticker = time.NewTicker(interval)
		// done       chan bool          = make(chan bool)
		sliceUrls []database.Feed = make([]database.Feed, 10)
		// sliceRss   []*xmldecoding.Rss = make([]*xmldecoding.Rss, 10)
		err error          = nil
		wg  sync.WaitGroup = sync.WaitGroup{}
	)

	for {
		select {
		case <-ticker.C:
			log.Println("Endless fetcher ticker has ticked...")
			sliceUrls, err = db.GetNextFeedsToFetch(context, feedAmount)

			if err != nil {
				log.Printf("Error getting next urls to fetch -> %v\n", err)
				return
			}
			for _, url := range sliceUrls {
				wg.Add(1)
				log.Println("Go routine started...")
				go processRss(url, &wg)
			}
			wg.Wait()
		}
	}
}

/*




	// one function fetches, hits another functions channel telling it to look
	// into a certain var to get the urls
	// then the third function's channel is hit to process the rss feeds













*/

// interval is passed to the NewCache() fn – which i haven't written yet
// this func is the called by NewCache()

//   // one function fetches, hits another functions channel telling it to look
//   // into a certain var to get the urls
//   // then the third function's channel is hit to process the rss feeds

// first function is the ticker. every 60 seconds, will talk to th

func processFeed() {
}

/*














 */
