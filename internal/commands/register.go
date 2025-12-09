package commands

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/luism2302/gator/internal/database"
	"time"
)

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) < 1 {
		return fmt.Errorf("error: No arguments given to register command. Usage: gator register <username>")
	}
	username := cmd.Args[0]

	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	newUser, err := s.Db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		return fmt.Errorf("error: couldnt create user: %w", err)
	}
	err = s.Cfg.SetUser(newUser.Name)
	if err != nil {
		return fmt.Errorf("error: couldnt set created user as current user: %w", err)
	}
	fmt.Printf("Created User: %s, at %v. With ID: %v\n", newUser.Name, newUser.CreatedAt, newUser.ID)
	return nil
}
