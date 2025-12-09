package commands

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("error: No arguments given to login command. Usage: gator login <username>")
	}
	username := strings.TrimSpace(cmd.Args[0])
	registeredUsers, err := s.Db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: couldnt get registered users from db: %w", err)
	}
	if !slices.Contains(registeredUsers, username) {
		return fmt.Errorf("error: username %s not registered in db", username)
	}
	err = s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error: couldnt set user in config: %w", err)
	}
	fmt.Printf("Logged in as: %s\n", username)
	return nil
}
