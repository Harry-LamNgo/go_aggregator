package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Harry-LamNgo/go_aggregator/internal/config"
	"github.com/Harry-LamNgo/go_aggregator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Open database from url - using 'postgres'
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	// Store database and cfg (config) to state
	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := Commands{
		registeredCommands: make(map[string]func(*state, CLICommand) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

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
