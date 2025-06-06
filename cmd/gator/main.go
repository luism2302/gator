package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/luism2302/gator/internal/cli"
	"github.com/luism2302/gator/internal/config"
	"github.com/luism2302/gator/internal/database"
)

func main() {
	//read config
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)
	//initialize state and commands
	state := &cli.State{
		Db:  dbQueries,
		Cfg: &cfg,
	}
	commands := cli.Commands{
		Name_to_function: map[string]func(s *cli.State, cmd cli.Command) error{},
	}
	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)
	commands.Register("reset", cli.HandlerReset)
	commands.Register("users", cli.HandlerUsers)
	commands.Register("agg", cli.HandlerAgg)
	commands.Register("addfeed", cli.MiddlewareLoggedIn(cli.HandlerAddFeed))
	commands.Register("feeds", cli.HandlerFeeds)
	commands.Register("follow", cli.MiddlewareLoggedIn(cli.HandlerFollow))
	commands.Register("following", cli.MiddlewareLoggedIn(cli.HandlerFollowing))
	commands.Register("unfollow", cli.MiddlewareLoggedIn(cli.HandlerUnfollow))

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
	}
	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	err = commands.Run(state, cli.Command{Name: commandName, Arguments: commandArgs})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
