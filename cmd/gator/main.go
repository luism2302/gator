package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/luism2302/gator/internal/cli"
	"github.com/luism2302/gator/internal/config"
	"github.com/luism2302/gator/internal/database"
	"github.com/luism2302/gator/internal/rss"
)

func main() {
	//read config
	config, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//check if num of arguments if valid
	if len(os.Args) < 2 {
		fmt.Println("too few arguments provided")
		os.Exit(1)
	}
	//db stuff
	dbURL := "postgres://postgres@localhost:5432/gator?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	//initialize state and commands
	state := cli.State{
		Db:  dbQueries,
		Cfg: config,
	}
	commands := cli.Commands{
		Name_to_function: map[string]func(s *cli.State, cmd cli.Command) error{},
	}

	//commands
	command_name := os.Args[1]
	args := os.Args[2:]
	command := cli.Command{
		Name:      command_name,
		Arguments: args,
	}
	switch command.Name {
	case "login":
		commands.Register(command.Name, cli.HandlerLogin)
	case "register":
		commands.Register(command.Name, cli.HandlerRegister)
	case "reset":
		commands.Register(command.Name, cli.HandlerReset)
	case "users":
		commands.Register(command.Name, cli.HandlerUsers)
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
	err = commands.Run(&state, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	test, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(test)
}
