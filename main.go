package main

import (
	"log"
	"os"

	"github.com/Harry-LamNgo/go_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: ", err)
	}

	// Store cfg (config) to state
	programState := &state{
		cfg: &cfg,
	}

	cmds := Commands{
		registeredCommands: make(map[string]func(*state, CLICommand) error),
	}

	cmds.register("login", handlerLogin)

	// Because the 1st argument of os is always program name/path of cli -> ignore this program name/path
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, CLICommand{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
