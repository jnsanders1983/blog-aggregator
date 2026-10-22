package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jnsanders1983/blog-aggregator/internal/database"
	"github.com/lib/pq"
)

func handlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs> (e.g. 1m)")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", cmd.Args[0], err)
	}
	if timeBetweenRequests <= 0 {
		return fmt.Errorf("time_between_reqs must be positive")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *State) {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("couldn't get next feed to fetch: %v\n", err)
		return
	}

	if err := s.db.MarkFeedFetched(ctx, feed.ID); err != nil {
		fmt.Printf("couldn't mark feed %s fetched: %v\n", feed.Name, err)
		return
	}

	fmt.Printf("Fetching %s (%s)\n", feed.Name, feed.Url)
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		fmt.Printf("couldn't fetch feed %s: %v\n", feed.Name, err)
		return
	}

	saved := 0
	for _, item := range rssFeed.Channel.Item {
		now := time.Now()
		_, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: parsePublishedAt(item.PubDate),
			FeedID:      feed.ID,
		})
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			}
			fmt.Printf("couldn't save post %q: %v\n", item.Link, err)
			continue
		}
		saved++
	}
	fmt.Printf("Feed %s: %d new posts (%d total)\n", feed.Name, saved, len(rssFeed.Channel.Item))
}

var pubDateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC3339,
	time.RFC3339Nano,
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

func parsePublishedAt(s string) sql.NullTime {
	s = strings.TrimSpace(s)
	if s == "" {
		return sql.NullTime{}
	}
	for _, layout := range pubDateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	return sql.NullTime{}
}
