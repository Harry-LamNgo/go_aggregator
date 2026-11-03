package main

import "errors"

type CLICommand struct {
	Name string
	Args []string
}

type Commands struct {
	// Create a map for registered commands for system to hold/remember which CLICommand is used
	registeredCommands map[string]func(*state, CLICommand) error
}

// function register command to register new command as f (function)
func (c *Commands) register(name string, f func(*state, CLICommand) error) {
	c.registeredCommands[name] = f
}

func (c *Commands) run(s *state, cmd CLICommand) error {
	// Check cmd (CLIcommand) is exist in map of system commands
	clicommand, exist := c.registeredCommands[cmd.Name]
	if !exist {
		return errors.New("command not found")
	}
	return clicommand(s, cmd)
}
