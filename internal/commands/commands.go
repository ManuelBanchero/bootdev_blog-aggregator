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
	"github.com/manuelbanchero/blog-aggregator/internal/rss"
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

func ListUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error trying to get users, err: %w", err)
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}

		fmt.Printf("* %v\n", user.Name)
	}

	return nil
}

func Agg(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed.Channel.Title)
	fmt.Println(feed.Channel.Description)
	fmt.Println(feed.Channel.Item)
	fmt.Println(feed.Channel.Link)

	return nil
}

func AddFeedHandler(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("commands.go | the AddFeedHandler expects two arguments, the name and the url.")
	}

	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("commands.go - AddFeedHandler() | An error has ocurred trying to get the user.\nError: %w", err)
	}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("commands.go - AddFeedHandler() | An error has ocurred trying to Create Feed.\nError: %w", err)
	}

	fmt.Println("Feed has been created successfully")
	fmt.Println(feed.ID)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)

	return nil
}

func FeedHandler(s *State, cmd Command) error {
	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("commands.go - FeedHandler() | An error has ocurred trying to get all feeds.\nError: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("commands.go - FeedHandler | An error has ocurred tyring to get the user by id.\nError: %w", err)
		}
		fmt.Printf("- Name: %v\n- URL: %v\n- Username: %v", feed.Name, feed.Url, user.Name)
	}

	return nil
}
