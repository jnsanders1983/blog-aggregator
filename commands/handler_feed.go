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

func feeds(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("feeds has no arguements")
	}
	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		name := feed.Name
		url := feed.Url
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s, Url: %s, Creator:%s\n", name, url, user.Name)
	}

	return nil
}
