package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd CLICommand) error {
	feedData, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed - %w", err)
	}

	fmt.Printf("%+v\n", feedData)

	return nil
}
