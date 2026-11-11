package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Harry-LamNgo/go_aggregator/internal/database"
)

func handlerLogin(s *state, cmd CLICommand) error {
	// len(cmd.Args) must specify 1 arg to prevent username like John Sand (2 args)
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user %q - %w", name, err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user %q - %w", name, err)
	}

	fmt.Printf("User switched successfully to %q\n", name)
	return nil
}

func handlerRegister(s *state, cmd CLICommand) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("error creating new user %q - %w", name, err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting current user %q - %w", user.Name, err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func handlerReset(s *state, cmd CLICommand) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	err := s.db.ResetData(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset data - %w", err)
	}
	fmt.Println("Successfully clear all data!")

	return nil
}

func handlerUsers(s *state, cmd CLICommand) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get all users data - %w", err)
	}

	if len(users) == 0 {
		fmt.Println("users data is empty")
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf(" * %v (current)\n", user.Name)
		} else {
			fmt.Printf(" * %v\n", user.Name)
		}
	}
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" - ID:		%v\n", user.ID)
	fmt.Printf(" - Name: 	%v\n", user.Name)
}
