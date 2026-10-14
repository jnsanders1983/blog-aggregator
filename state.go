package main

import (
	"fmt"

	"github.com/jnsanders1983/blog-aggregator/internal/config"
)

type state struct {
	c *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandHandler map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f := c.commandHandler[cmd.name]
	if f == nil {
		return fmt.Errorf("invalid command: %s", cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandHandler[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username.")
	}
	err := s.c.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("The user, %s, has been set!\n", s.c.CurrentUserName)
	return nil
}
