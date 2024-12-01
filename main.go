package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
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
			"login":    handlerLogin,
			"register": handlerRegister,
			"reset":    handlerReset,
			"users":    handlerUsers,
			"agg":      handlerAgg,
			"addfeed":  handlerAddFeed,
		},
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments.")
		os.Exit(1)
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username required")
	}

	ctx := context.Background()

	if _, err := s.db.GetUser(ctx, cmd.args[0]); err != nil {
		if err != nil {
			return errors.New("user does not exist")
		}
		return errors.New("user already exists")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User set to %s\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username required")
	}

	ctx := context.Background()

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	usr, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	err = s.config.SetUser(usr.Name)
	if err != nil {
		return err
	}

	fmt.Printf("New user created\nName:%s\nCreatedAt:%v\nUpdatedAt:%v\nID:%v\n", usr.Name, usr.CreatedAt, usr.UpdatedAt, usr.ID.ID())

	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	err := s.db.Reset(ctx)
	if err != nil {
		return err
	}

	fmt.Println("reset successful")

	return nil
}

func handlerUsers(s *state, cmd command) error {

	ctx := context.Background()

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {

		fmt.Printf("* %s ", user)
		if user == s.config.CurrentUser {
			fmt.Println("(current)")
		} else {
			fmt.Println()
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()

	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("need name and url args")
	}

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.config.CurrentUser)
	if err != nil {
		return err
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed := database.Feed{}

	feed, err = s.db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}

	fmt.Println(feed)

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
