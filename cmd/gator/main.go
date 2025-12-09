package main

import (
	"fmt"
	"github.com/luism2302/gator/internal/commands"
	"github.com/luism2302/gator/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	state := commands.State{
		Cfg: &cfg,
	}
	cmds := commands.Commands{
		MapHandler: make(map[string]func(*commands.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin)
	args := os.Args
	if len(args) < 2 {
		log.Fatal("error: no command provided. Usage: gator <command> [args]")
	}
	cmdName := args[1]
	cmdArgs := []string{}
	if len(args) > 2 {
		cmdArgs = args[2:]
	}
	cmd := commands.Command{
		Name: cmdName,
		Args: cmdArgs,
	}
	err = cmds.Run(&state, cmd)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Println(cfg)
}
