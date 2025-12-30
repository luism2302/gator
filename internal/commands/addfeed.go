package commands

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/luism2302/gator/internal/database"
	"time"
)

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("error: expected args <name> <url>")
	}
	loggedUser, err := s.Db.GetUserByName(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error: couldnt get user: %s from db", s.Cfg.CurrentUserName)
	}
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]
	newFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    loggedUser.ID,
	}

	newFeed, err := s.Db.CreateFeed(context.Background(), newFeedParams)
	if err != nil {
		return fmt.Errorf("error: couldnt create feed")
	}
	fmt.Printf("Succesfully Created Feed. ID: %v, URL: %s, Name: %s. Created by UserID: %v at time: %v\n", newFeed.ID, newFeed.Url, newFeed.Name, newFeed.UserID, newFeed.CreatedAt)
	return nil
}
