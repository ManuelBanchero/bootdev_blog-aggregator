package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/manuelbanchero/blog-aggregator/internal/config"
	"github.com/manuelbanchero/blog-aggregator/internal/database"
)

type State struct {
	Db  *database.Queries
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

	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("The user does not exists, err: %w", err)
	}

	s.Cfg.CurrentUserName = username

	// Set the username
	s.Cfg.SetUser()

	fmt.Printf("The user '%v' has been set", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the register handler expectes a single argument, the username")
	}

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
	})

	if err != nil {
		if pqNotUniqueErr := pq.As(err, pqerror.UniqueViolation); pqNotUniqueErr != nil {
			return fmt.Errorf("The user Name already exists")
		}

		return err
	}

	// set the new username
	s.Cfg.CurrentUserName = user.Name
	if err := s.Cfg.SetUser(); err != nil {
		return err
	}

	fmt.Println("The user was created")
	fmt.Printf("USER DATA\n- ID: %v\n- CreatedAt: %v\n- UpdatedAt: %v\n- Name: %v\n", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)

	return nil
}

func HandlerResetDB(s *State, cmd Command) error {
	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("Error trying to delete users from db, error: %w", err)
	}

	fmt.Println("All rows from User table has been deleted")
	return nil
}
