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

	fmt.Println(s.Cfg.CurrentUserName)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Must include a command")
		return
	}

}
