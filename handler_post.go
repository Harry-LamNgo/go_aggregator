package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Harry-LamNgo/go_aggregator/internal/database"
)

func handlerBrowsePost(s *state, cmd CLICommand, user database.User) error {
	// Default limit
	var limit int32
	limit = 2
	if len(cmd.Args) >= 2 {
		return fmt.Errorf("usage: %v [<limit>]", cmd.Name)
	} else if len(cmd.Args) == 1 {
		parseLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit %q - %w", cmd.Args[0], err)
		}
		limit = int32(parseLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("error getting posts for user %v - %w", user.ID, err)
	}

	fmt.Printf("Successfully get posts for %v\n", user.Name)
	fmt.Println("==================================================")
	for _, post := range posts {
		fmt.Printf("%v from %v\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %v ---\n", post.Title)
		fmt.Printf("   %v\n", post.Description.String)
		fmt.Printf("** Link: %v\n", post.Url)
		fmt.Println("==================================================")
	}

	return nil
}
