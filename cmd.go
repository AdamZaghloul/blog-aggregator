package main

import (
	"blog-aggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"strconv"
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
	if len(cmd.args) < 1 {
		return errors.New("need duration like 1s, 1m, 1h, etc")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests.String())

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("need name and url args")
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	_, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	cmd.args = cmd.args[1:]

	handlerFollow(s, cmd, user)

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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("not enough arguments; need url to follow")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	newFeed := database.CreateFeedFollowRow{}

	newFeed, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%s is now following %s\n", user.Name, newFeed.FeedName)
	fmt.Println()

	return nil
}
func handlerFollowing(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println()
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	fmt.Println()

	return nil

}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 1 {
		return errors.New("not enough arguments; need url to unfollow")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%s is now following %s\n", user.Name, feed.Name)
	fmt.Println()

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	var err error

	if len(cmd.args) > 0 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			limit = 2
		}
	}

	params := database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println()

	for _, post := range posts {
		fmt.Printf("%s - %s - %v\n", post.Title, post.Name.String, post.PublishedAt)
	}

	fmt.Println()

	return nil
}
