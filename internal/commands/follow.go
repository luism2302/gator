package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luism2302/gator/internal/database"
)

func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: follow only accepts a url command")
	}
	url := cmd.Args[0]
	fetchedFeed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error: couldnt get feed: %w", err)
	}
	currentUser, err := s.Db.GetUserByName(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error: couldnt get current user: %w", err)
	}

	newfeedParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    fetchedFeed.ID,
	}
	newFeed, err := s.Db.CreateFeedFollow(context.Background(), newfeedParams)
	if err != nil {
		return fmt.Errorf("error: couldnt create new feed follow: %w", err)
	}
	fmt.Printf("User: %s now follows feed: %s\n", newFeed.UserName, newFeed.FeedName)
	return nil
}
