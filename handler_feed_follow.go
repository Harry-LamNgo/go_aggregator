package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Harry-LamNgo/go_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd CLICommand, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	feedUrl, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedUrl.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating new feed follow for the current user - %w", err)
	}

	fmt.Println("Successfully create new feed follow record for the current user")
	fmt.Printf(" - Name feed: %v\n", feedFollow.FeedName)
	fmt.Printf(" - User: %v\n", feedFollow.UserName)

	return nil
}

func handlerListFeedFollows(s *state, cmd CLICommand, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	userFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error finding current user feed follows - %w", err)
	}

	if len(userFeeds) == 0 {
		fmt.Printf("No feed follows record found for current user.")
		return nil
	}

	fmt.Printf("User %v following feeds:\n", user.Name)
	for _, userFeed := range userFeeds {
		fmt.Printf(" - %v\n", userFeed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd CLICommand, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	feedURL, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	err = s.db.RemoveFeedFollow(context.Background(), database.RemoveFeedFollowParams{
		UserID: user.ID,
		FeedID: feedURL.ID,
	})
	if err != nil {
		return fmt.Errorf("error removing feed follow record - %w", err)
	}

	fmt.Println("Successfully unfollow feed record.")

	return nil
}
