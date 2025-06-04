package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luism2302/gator/internal/config"
	"github.com/luism2302/gator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	Name_to_function map[string]func(s *State, cmd Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if s == nil {
		return fmt.Errorf("nil pointer to state")
	}
	err := c.Name_to_function[cmd.Name](s, cmd)
	return err
}
func (c *Commands) Register(name string, f func(s *State, cmd Command) error) {
	c.Name_to_function[name] = f
}
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("you must provide a username as an argument")
	}
	if s == nil {
		return fmt.Errorf("nil pointer to state")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("user %s doesnt exist", cmd.Arguments[0])
	}
	s.Cfg.SetUser(cmd.Arguments[0])
	fmt.Printf("user: %s has been set", cmd.Arguments[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("you must pass a name as an argument")
	}
	name := cmd.Arguments[0]
	_, err := s.Db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("user %s already exists", name)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	newUser, err2 := s.Db.CreateUser(context.Background(), params)
	if err2 != nil {
		return err2
	}
	s.Cfg.SetUser(name)
	fmt.Printf("created user %s: %v", name, newUser)
	return nil
}
