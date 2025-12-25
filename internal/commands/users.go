package commands

import (
	"context"
	"fmt"
)

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error: users command doesnt accept arguments")
	}
	allUsers, err := s.Db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: couldnt get users from database: %w", err)
	}

	for _, user := range allUsers {
		if user == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}
