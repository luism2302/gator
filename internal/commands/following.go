package commands

import (
	"context"
	"fmt"
)

func HandlerFollowing(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error: following doesnt accept arguments")
	}
	currentUser, err := s.Db.GetUserByName(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error: couldnt get current user: %w", err)
	}
	following, err := s.Db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("error: couldnt get followed feeds for user: %w", err)
	}
	fmt.Printf("User: %s follows:\n", currentUser.Name)
	for _, followed := range following {
		fmt.Printf("-\t%s\n", followed)
	}
	return nil
}
