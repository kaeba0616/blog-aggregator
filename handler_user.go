package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kaeba0616/blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), name)

	if err != nil {
		fmt.Println("Didn't find the user")
		log.Fatal(err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user %w", err)
	}

	fmt.Println("User switched successfully!!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	name := cmd.Args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatal("already exists the name...")
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("User are created successfully!")
	printUser(user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf(" - %v (current)\n", user.Name)
		} else {
			fmt.Printf(" - %v\n", user.Name)
		}
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" - ID:       %v\n", user.ID)
	fmt.Printf(" - Name:     %v\n", user.Name)
}
