package main

import (
	"database/sql"
	"fmt"
	"gator/internal/command"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/state"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.Db_url)
	db_queries := database.New(db)
	s := state.State{Config_ptr: &cfg, Db_ptr: db_queries}

	cmds := command.Commands{Cmd_map: make(map[string]func(*state.State, command.Command) error)}
	cmds.Register("login", command.HandlerLogin)
	cmds.Register("register", command.HandlerRegister)

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Too few arguments")
		os.Exit(1)
	}

	cmd := command.Command{Name: args[0], Args: args[1:]}

	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Println("ERROR!: ", err)
		os.Exit(1)
	}
}
