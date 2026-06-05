package command

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/database"
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
