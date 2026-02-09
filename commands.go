package main

import (
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
)

type state struct {
	db     *database.Queries
	config config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	newCommand, exists := c.commands[cmd.name]
	if !exists {
		return fmt.Errorf("command '%s' not registered\n", cmd.name)
	}
	fmt.Printf("running command: %s\n", cmd.name)
	err := newCommand(s, cmd)
	if err != nil {
		return fmt.Errorf("execution error: %v", err)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f

}
