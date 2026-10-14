package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jnsanders1983/blog-aggregator/internal/config"
	"github.com/jnsanders1983/blog-aggregator/internal/database"
)

type state struct {
	c  *config.Config
	db *database.Queries
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
		return fmt.Errorf("login requires exactly one username")
	}

	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("get user %q: %w", name, err)
	}

	if err := s.c.SetUser(name); err != nil {
		return fmt.Errorf("set current user: %w", err)
	}

	fmt.Printf("The user, %s, has been set!\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register requires exactly one username")
	}

	name := cmd.args[0]
	ctx := context.Background()

	_, err := s.db.GetUser(ctx, name)
	if err == nil {
		return fmt.Errorf("user %q already exists", name)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("check existing user: %w", err)
	}

	now := time.Now()
	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	if err := s.c.SetUser(name); err != nil {
		return fmt.Errorf("set current user: %w", err)
	}

	fmt.Printf("User %s was created!\n", user.Name)
	fmt.Printf("%+v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset does not accept arguments")
	}

	if err := s.db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("reset users: %w", err)
	}

	fmt.Println("All users were deleted successfully.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("users does not accept arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.c.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
