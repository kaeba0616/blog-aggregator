package main

import (
	"context"
	"fmt"

	"github.com/kaeba0616/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("no user found: %w", err)
		}

		return handler(s, cmd, user)
	}
}
