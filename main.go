package main

import (
	"fmt"
	"os"

	"github.com/manuelbanchero/blog-aggregator/internal/commands"
	"github.com/manuelbanchero/blog-aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	s := commands.State{
		Cfg: &c,
	}

	cmds := commands.Commands{
		Methods: map[string]func(*commands.State, commands.Command) error{},
	}

	cmds.Register("login", commands.HandlerLogin)

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
