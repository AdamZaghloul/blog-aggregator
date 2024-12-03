package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
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

	db, err := sql.Open("postgres", cfg.DbUrl)
	dbQueries := database.New(db)
	ste.db = dbQueries

	cmds := commands{
		cmds: map[string]func(*state, command) error{
			"login":     handlerLogin,
			"register":  handlerRegister,
			"reset":     handlerReset,
			"users":     handlerUsers,
			"agg":       handlerAgg,
			"addfeed":   middlewareLoggedIn(handlerAddFeed),
			"feeds":     handlerFeeds,
			"follow":    middlewareLoggedIn(handlerFollow),
			"following": middlewareLoggedIn(handlerFollowing),
			"unfollow":  middlewareLoggedIn(handlerUnfollow),
			"browse":    middlewareLoggedIn(handlerBrowse),
		},
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments.")
		os.Exit(1)
	}

	if args[1] == "help" {
		printHelp(cmds)
		return
	}

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

func printHelp(cmds commands) {
	fmt.Println()
	for key, _ := range cmds.cmds {
		fmt.Println(key)
	}
	fmt.Println()
}
