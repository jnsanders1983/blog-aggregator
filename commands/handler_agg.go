package commands

import (
	"context"
	"fmt"
	"time"
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

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("* %s\n", item.Title)
	}
	fmt.Printf("Feed %s: %d posts\n", feed.Name, len(rssFeed.Channel.Item))
}
