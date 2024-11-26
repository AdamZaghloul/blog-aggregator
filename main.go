package main

import (
	"blog-aggregator/internal/config"
	"errors"
	"fmt"
	"os"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	ste := state{}
	ste.config = &cfg

	cmds := commands{
		cmds: map[string]func(*state, command) error{
			"login": handlerLogin,
		},
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments.")
		os.Exit(1)
	}

	fmt.Println(args)

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&ste, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username required")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User set to %s\n", cmd.args[0])
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	_, ok := c.cmds[cmd.name]

	if !ok {
		return errors.New("no such command")
	}

	err := c.cmds[cmd.name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}
