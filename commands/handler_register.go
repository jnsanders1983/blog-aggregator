package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jnsanders1983/blog-aggregator/internal/database"
)

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("register requires exactly one username")
	}

	name := cmd.Args[0]
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
