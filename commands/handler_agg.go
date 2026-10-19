package commands

import (
	"context"
	"fmt"
)

func handlerAgg(s *State, cmd Command) error {
	ctx := context.Background()

	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}
