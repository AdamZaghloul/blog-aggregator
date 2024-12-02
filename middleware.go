package main

import (
	"blog-aggregator/internal/database"
	"context"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		ctx := context.Background()

		user, err := s.db.GetUser(ctx, s.config.CurrentUser)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
