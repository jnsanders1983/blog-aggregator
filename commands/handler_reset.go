package commands

import (
	"context"
	"fmt"
)

func handlerReset(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("reset does not accept arguments")
	}

	if err := s.db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("reset users: %w", err)
	}

	fmt.Println("All users were deleted successfully.")
	return nil
}
