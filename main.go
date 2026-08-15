package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/manuelbanchero/blog-aggregator/internal/commands"
	"github.com/manuelbanchero/blog-aggregator/internal/config"
	"github.com/manuelbanchero/blog-aggregator/internal/database"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := commands.State{
		Db:  dbQueries,
		Cfg: &c,
	}

	cmds := commands.Commands{
		Methods: map[string]func(*commands.State, commands.Command) error{},
	}

	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerResetDB)
	cmds.Register("users", commands.ListUsers)
	cmds.Register("agg", commands.Agg)
	cmds.Register("addfeed", commands.AddFeedHandler)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Must include a command")
		os.Exit(1)
	}

	commandName := args[1]
	arguments := args[2:]

	command := commands.Command{
		Name: commandName,
		Args: arguments,
	}

	if err := cmds.Run(&s, command); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
