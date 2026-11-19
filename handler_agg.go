package main

import (
	"context"
	"fmt"
	"log"
	"time"
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
	// for init; condition; post
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
		fmt.Printf("Found post: %v\n", item.Title)
	}

	log.Printf("Feed %v collected, %v posts found", nextFeed.Name, len(feedData.Channel.Item))
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
