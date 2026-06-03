package main

import "gator/internal/config"
import "fmt"
import "errors"
import "os"

type state struct {
	config_ptr *config.Config
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

	err := s.config_ptr.SetUser(cmd.Args[0])
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

func main() {
	cfg := config.Read()
	s := state{config_ptr: &cfg}
	cmds := commands{cmd_map: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Too few arguments")
		os.Exit(1)
	}

	cmd := command{Name: args[0], Args: args[1:]}

	err := cmds.run(&s, cmd)
	if err != nil {
		fmt.Println("ERROR!: ", err)
		os.Exit(1)
	}
}
