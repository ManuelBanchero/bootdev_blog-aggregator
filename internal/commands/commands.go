package commands

import (
	"fmt"

	"github.com/manuelbanchero/blog-aggregator/internal/config"
)

type State struct {
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Methods map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	method, ok := c.Methods[cmd.Name]
	if !ok {
		return fmt.Errorf("The command does not exists")
	}

	if err := method(s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Methods[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	username := cmd.Args[0]
	s.Cfg.CurrentUserName = username

	// Set the username
	s.Cfg.SetUser()

	fmt.Printf("The user '%v' has been set", username)
	return nil
}
