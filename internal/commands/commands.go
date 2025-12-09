package commands

import (
	"fmt"
	"github.com/luism2302/gator/internal/config"
	"github.com/luism2302/gator/internal/database"
)

type State struct {
	Cfg *config.Config
	Db  *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	MapHandler map[string]func(*State, Command) error
}

func (cmds *Commands) Run(s *State, cmd Command) error {
	handler, ok := cmds.MapHandler[cmd.Name]
	if !ok {
		return fmt.Errorf("error: %s command not found", cmd.Name)
	}
	err := handler(s, cmd)
	if err != nil {
		return fmt.Errorf("error: couldnt run command %s: %w", cmd.Name, err)
	}
	return nil
}

func (cmds *Commands) Register(name string, handler func(*State, Command) error) {
	cmds.MapHandler[name] = handler
}
