package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Harry-LamNgo/go_aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd CLICommand) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time between requests>", cmd.Name)
	}

	timeDur := cmd.Args[0]

	timeBetweenReqs, err := time.ParseDuration(timeDur)
	if err != nil {
		return fmt.Errorf("invalid duration - %w", err)
	}

	log.Printf("Collecting feeds every %s...\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	// for init; condition; post ---> Why use this ? --> At the end of this file
	// init omit - condition omit --> does not need initialize and runs infinite loops
	// post: "<-ticker.C" means receive "tick" from ticker channel
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) {
	// Get the next feed to fetch from the Database
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")

	// Mark that feed --> update it last time fetched to pushing it to the back of the queue
	// so it won't be picked up again right away.
	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %v fetched: %v", nextFeed.Name, err)
		return
	}

	// Fetch latest post (data) from Mark feed URL
	feedData, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %v: %v", nextFeed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: nullstr(item.Description),
			PublishedAt: parseRSSDate(item.PubDate),
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
					continue
				}
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	log.Printf("Feed %v collected, %v posts found", nextFeed.Name, len(feedData.Channel.Item))

}

// To use string nullable in postgresql --> using sql.NullString to convert feed string
func nullstr(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

// To make String Datetime in feed into TIMESTAMP (NULLABLE) in postgresql
// Using time.Parse with constant time.<layout>
func parseRSSDate(dateStr string) sql.NullTime {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return sql.NullTime{Valid: false}
	}

	timeFormat := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC822,
		time.RFC822Z,
	}

	for _, layout := range timeFormat {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return sql.NullTime{
				Time:  t.UTC(),
				Valid: true,
			}
		}
	}
	return sql.NullTime{Valid: false}
}

/*
Why write like this ?
for ; ; <-ticker.C {
	body code
}

Not:
for range ticker.C {
	body code
}
Because this not immediately run body code, it waits for time duration
ticker.C to run the 1st loop of body code

*/
