package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jnsanders1983/blog-aggregator/internal/database"
)

func addfeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addfeed requires name and url")
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
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("Name: %s\n", feed.Name)

	feed_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feed_follow_row, err := s.db.CreateFeedFollow(context.Background(), feed_follow)
	if err != nil {
		return err
	}
	fmt.Printf("Name: %s\n", feed_follow_row.FeedName)
	return nil
}

func feeds(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("feeds has no arguments")
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

func follow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("follow takes exactly one argument, a url")
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Printf("Name: %s\n", feed.Name)

	now := time.Now()
	feed_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feed_follow_row, err := s.db.CreateFeedFollow(context.Background(), feed_follow)
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %s, User: %s\n", feed_follow_row.FeedName, feed_follow_row.UserName)

	return nil
}

func unfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("unfollow takes exactly one argument, a url")
	}
	url := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	remove_feed := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	return s.db.DeleteFeedFollow(context.Background(), remove_feed)
}

func following(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("following takes zeor arguments")
	}

	user_feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed := range user_feed_follows {
		fmt.Printf("Name: %s\n", feed.FeedName)
	}
	return nil
}
