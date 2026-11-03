package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd CLICommand) error {
	// len(cmd.Args) must specify 1 arg to prevent username like John Sand (2 args)
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>}", cmd.Name)
	}

	name := cmd.Args[0]

	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("could not set current user %w", err)
	}

	fmt.Printf("User switched successfully to %q\n", name)
	return nil
}
