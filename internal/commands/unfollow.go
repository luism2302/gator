package commands

import (
	"context"
	"fmt"
	"github.com/luism2302/gator/internal/database"
)

func HandlerUnfollow(s *State, cmd Command, u database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("error: url not provided")
	}
	deleteFeed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error: couldnt fetch feed: %w", err)
	}
	deleteParams := database.DeleteFeedFollowParams{
		UserID: u.ID,
		FeedID: deleteFeed.ID,
	}
	err = s.Db.DeleteFeedFollow(context.Background(), deleteParams)
	if err != nil {
		return fmt.Errorf("error: couldnt delete feed: %w", err)
	}
	return nil
}
