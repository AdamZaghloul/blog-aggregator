package main

import (
	"blog-aggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

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

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	fmt.Println()

	for _, feed := range feeds {
		fmt.Printf("%s - %s - %s\n", feed.Name, feed.Url, feed.UserName.String)
	}

	fmt.Println()

	return nil
}
