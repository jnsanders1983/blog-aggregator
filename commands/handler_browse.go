package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jnsanders1983/blog-aggregator/internal/database"
)

func browse(s *State, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: browse [limit]")
	}
	if len(cmd.Args) == 1 {
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil || n <= 0 {
			return fmt.Errorf("limit must be a positive integer")
		}
		limit = n
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Printf("%s (%s)\n", p.Title, p.FeedName)
		if p.PublishedAt.Valid {
			fmt.Printf("  Published: %s\n", p.PublishedAt.Time.Format("Mon Jan 2 2006"))
		}
		fmt.Printf("  %s\n", p.Url)
		if p.Description.Valid {
			fmt.Printf("  %s\n", p.Description.String)
		}
		fmt.Println()
	}
	return nil
}
