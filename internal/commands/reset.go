package commands

import (
	"context"
	"fmt"
)

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error: reset command doesnt accept arguments")
	}
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: couldnt delete users from database: %w", err)
	}
	return nil
}
