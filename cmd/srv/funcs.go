/*
current possible bug:
  there is nothing preventing multiple bad urls from flooding the db.
  considering we only fetch a predetemined amount of feeds at once, and nulls
  come first, with enough of these bad rss feeds, it would really slow down the
  fetching, even to a halt, given enough feeds.

fix part 1: regexp to make sure at least the format of the url is valid

    part 2: new columns on database. 1 to track the amount of times we have been
            unable to fetch them,    1 to flag it after a certain amount of tries
*/

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

	"github.com/google/uuid"

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

	time1 := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	arg := database.MarkFeedFetchedParams{
		LastFetchedAt: time1,
		ID:            feed.ID,
	}
	db.MarkFeedFetched(context.Background(), arg)

	// initializing some args for the creation of the post
	var (
		arg2        database.CreatePostParams
		time2       time.Time
		description sql.NullString
	)

	for _, item := range rssStruct.Channel.Item {
		time2, err = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", item.PubDate)
		if err != nil {
			log.Printf("Error parsing time from rss item -> %v\n", err)
			return err
		}

		if item.Description == "" {
			description = sql.NullString{
				String: "",
				Valid:  false,
			}
		} else {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}

		arg2 = database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: time2,
			FeedID:      feed.ID,
		}

		db.CreatePost(context.Background(), arg2)
	}

	// arg = database.CreatePostParams{
	// 	ID:          uuid.New(),
	// 	CreatedAt:   time.Now(),
	// 	Title:
	// 	Url:         "",
	// 	Description: sql.NullString{},
	// 	PublishedAt: time.Time{},
	// 	FeedID:      [16]byte{},
	// }

	// Item          []struct {
	// 	Text        string `xml:",chardata"`
	// 	Title       string `xml:"title"`
	// 	Link        string `xml:"link"`
	// 	PubDate     string `xml:"pubDate"`
	// 	Guid        string `xml:"guid"`
	// 	Description string `xml:"description"`
	// } `xml:"item"`

	log.Printf("RSS <%v>\n", rssStruct.Channel.Title)
	log.Println("Go routine finished...")
	return nil
}

func endlessFetching(context context.Context) {
	fmt.Print("\n\n\n")
	log.Printf("Running endless fetching...")

	const interval = 60 * time.Second

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
