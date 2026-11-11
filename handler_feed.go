package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Harry-LamNgo/go_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd CLICommand) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	userFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to current user - %w", err)
	}

	fmt.Println("Successfully add new feed to current user!")
	fmt.Println("==================================================")
	fmt.Printf("- Name:	%v\n", currentUser.Name)
	fmt.Printf(" * Created:		%v\n", userFeed.CreatedAt)
	fmt.Printf(" * Feed name:	%v\n", userFeed.Name)
	fmt.Printf(" * Feed URL: 	%v\n", userFeed.Url)

	return nil
}
