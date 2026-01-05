package commands

import (
	"context"
	"fmt"
	"github.com/luism2302/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(s *State, cmd Command) error {
	return func(s *State, cmd Command) error {
		loggedUser, err := s.Db.GetUserByName(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error: couldnt get current user: %w", err)
		}
		return handler(s, cmd, loggedUser)
	}
}
