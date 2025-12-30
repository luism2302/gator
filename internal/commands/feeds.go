package commands

import (
	"context"
	"fmt"
)

func HandlerFeeds(s *State, cmd Command) error {
	allFeeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error: couldnt get feeds from db: %w", err)
	}
	for _, feed := range allFeeds {
		fmt.Printf("Name: %s, URL: %s, Added By: %s\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}
