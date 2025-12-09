package commands

import (
	"fmt"
	"strings"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("error: No arguments given to login command. Usage: gator login <username>")
	}
	username := strings.TrimSpace(cmd.Args[0])
	err := s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error: couldnt set user in config: %w", err)
	}
	fmt.Println("succesfully set user")
	return nil
}
