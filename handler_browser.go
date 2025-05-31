package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kaeba0616/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if limitArg, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = limitArg
		} else {
			return fmt.Errorf("couldn't convert string to int32: %w", err)
		}
	}

	gpparams := database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostForUser(context.Background(), gpparams)
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}
	printPosts(posts)
	return nil
}

func printPosts(posts []database.GetPostForUserRow) {
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("========================================")
	}
}
