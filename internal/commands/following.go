package commands

import (
	"context"
	"fmt"

	"github.com/luism2302/gator/internal/database"
)

func HandlerFollowing(s *State, cmd Command, u database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error: following doesnt accept arguments")
	}

	following, err := s.Db.GetFeedFollowsForUser(context.Background(), u.ID)
	if err != nil {
		return fmt.Errorf("error: couldnt get followed feeds for user: %w", err)
	}
	fmt.Printf("User: %s follows:\n", u.Name)
	for _, followed := range following {
		fmt.Printf("-\t%s\n", followed)
	}
	return nil
}
