package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jnsanders1983/blog-aggregator/internal/database"
)

func addfeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addfeed requires name and url")
	}
	ctx := context.Background()
	name := s.c.CurrentUserName

	user, err := s.db.GetUser(ctx, name)
	if err != nil {
		return err
	}

	now := time.Now()

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}
	fmt.Printf("Name: %s\n", feed.Name)
	return nil
}
