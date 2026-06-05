package command

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/database"
	"gator/internal/rss"
	"gator/internal/state"
	"time"

	"github.com/google/uuid"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Cmd_map map[string]func(*state.State, Command) error
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	if f, ok := c.Cmd_map[cmd.Name]; ok {
		err := f(s, cmd)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Unknown command")
	}
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.Cmd_map[name] = f
}

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Incorrect amount of args passed to login")
	}

	_, err := s.Db_ptr.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.Config_ptr.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}

func HandlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("Incorrect amount of args passed to register")
	}
	current_time := time.Now()
	input_args := database.CreateUserParams{ID: uuid.New(), CreatedAt: current_time, UpdatedAt: current_time, Name: cmd.Args[0]}
	usr, err := s.Db_ptr.CreateUser(context.Background(), input_args)
	if err != nil {
		return err
	}

	err = HandlerLogin(s, cmd)
	if err != nil {
		return err
	}

	fmt.Println("User has been created: ", usr)
	return nil
}

func HandlerReset(s *state.State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("Reset does not take any arguments")
	}
	err := s.Db_ptr.Reset(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Reset successfully")
	return nil
}

func HandlerUsers(s *state.State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("Users does not take any arguments")
	}

	usrs, err := s.Db_ptr.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, usr := range usrs {
		status := ""
		if usr.Name == s.Config_ptr.Current_user_name {
			status += " (current)"
		}
		fmt.Println("* " + usr.Name + status)
	}
	return nil
}

func HandlerAgg(s *state.State, cmd Command) error {
	// Will need to change this in the futur to accept args
	if len(cmd.Args) != 0 {
		return errors.New("Agg does not take any arguments")
	}

	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func HandlerAddFeed(s *state.State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return errors.New("Incorrect amount of arguments passed into addfeed")
	}
	usr, err := s.Db_ptr.GetUser(context.Background(), s.Config_ptr.Current_user_name)
	if err != nil {
		return err
	}

	current_time := time.Now()
	input_args := database.CreateFeedParams{ID: uuid.New(), CreatedAt: current_time, UpdatedAt: current_time, Name: cmd.Args[0], Url: cmd.Args[1], UserID: usr.ID}
	feed, err := s.Db_ptr.CreateFeed(context.Background(), input_args)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func HandlerFeeds(s *state.State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("Feeds does not take any arguments")
	}

	feeds, err := s.Db_ptr.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		name, err := s.Db_ptr.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(feed.Name + " " + feed.Url + " " + name)
	}
	return nil
}
