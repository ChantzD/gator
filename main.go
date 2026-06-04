package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"github.com/google/uuid"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type state struct {
	config_ptr *config.Config
	db_ptr     *database.Queries
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	cmd_map map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Incorrect amount of args passed to login")
	}

	_, err := s.db_ptr.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.config_ptr.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.cmd_map[cmd.Name]; ok {
		err := f(s, cmd)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Unknown command")
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd_map[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("Incorrect amount of args passed to register")
	}
	current_time := time.Now()
	input_args := database.CreateUserParams{ID: uuid.New(), CreatedAt: current_time, UpdatedAt: current_time, Name: cmd.Args[0]}
	usr, err := s.db_ptr.CreateUser(context.Background(), input_args)
	if err != nil {
		return err
	}

	err = handlerLogin(s, cmd)
	if err != nil {
		return err
	}

	fmt.Println("User has been created: ", usr)
	return nil
}

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.Db_url)
	db_queries := database.New(db)
	s := state{config_ptr: &cfg, db_ptr: db_queries}

	cmds := commands{cmd_map: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Too few arguments")
		os.Exit(1)
	}

	cmd := command{Name: args[0], Args: args[1:]}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println("ERROR!: ", err)
		os.Exit(1)
	}
}
