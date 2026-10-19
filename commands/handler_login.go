package commands

import (
	"context"
	"fmt"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login requires exactly one username")
	}

	name := cmd.Args[0]
	if _, err := s.db.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("get user %q: %w", name, err)
	}

	if err := s.c.SetUser(name); err != nil {
		return fmt.Errorf("set current user: %w", err)
	}

	fmt.Printf("The user, %s, has been set!\n", name)
	return nil
}
