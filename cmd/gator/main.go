package main

import _ "github.com/lib/pq"
import (
	"database/sql"
	"github.com/luism2302/gator/internal/commands"
	"github.com/luism2302/gator/internal/config"
	"github.com/luism2302/gator/internal/database"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	state := commands.State{
		Cfg: &cfg,
		Db:  dbQueries,
	}

	var cmds = commands.Commands{
		MapHandler: make(map[string]func(*commands.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAgg)

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

}
