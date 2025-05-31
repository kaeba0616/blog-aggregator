package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kaeba0616/blog-aggregator/internal/database"
)

func handlerAggregator(s *state, cmd command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenReqs := cmd.Args[0]
	duration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("couldn't parse to time <time_between_reqs>: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", duration)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {

	feed, err := s.db.GetNextFeedFetched(context.Background())
	if err != nil {
		log.Printf("couldn't get next feed fetched: %v", err)
		return
	}

	log.Println("Found a feed to fetch!")
	scrapeFeed(feed, s.db)
}

func scrapeFeed(feed database.Feed, db *database.Queries) {
	_, err := db.UpdateFeedFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("couldn't mark feed fetched: %v", err)
		return
	}

	fetch, err := fetched(context.Background(), feed.Url)
	if err != nil {
		log.Printf("couldn't fetch the feed with <%s>: %v", feed.Url, err)
		return
	}

	for _, item := range fetch.Channel.Item {
		publishedAt := time.Time{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = t
		}

		cpparams := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}
		_, err := db.CreatePost(context.Background(), cpparams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("couldn't create a post: %v", err)
			continue
		}
	}
	log.Printf("Feeds %s collected, %v posts found", feed.Name, len(fetch.Channel.Item))
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	new_params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), new_params)
	if err != nil {
		return fmt.Errorf("couldn't create feed_follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println("========================================")

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	fmt.Printf("Found %d feeds:\n", len(feeds))

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't find user: %w", err)
		}
		printFeed(feed, user)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("no feed found: %w", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("no follow found: %w", err)
	}
	printFollow(follow.UserName, follow.FeedName)
	return nil

}

func handlerListFollows(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("No follow found")
		return nil
	}

	for _, follow := range follows {
		printFollow(follow.UserName, follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get the feed: %w", err)
	}

	dfparams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), dfparams)
	if err != nil {
		return fmt.Errorf("couldn't delete feed-follow: %w", err)
	}

	return nil
}

func printFollow(username, feedname string) {
	fmt.Printf(" - %s\n", username)
	fmt.Printf(" - %s\n", feedname)
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf(" - ID:			%s\n", feed.ID)
	fmt.Printf(" - CreatedAt:	%s\n", feed.CreatedAt)
	fmt.Printf(" - UpdatedAt:	%s\n", feed.UpdatedAt)
	fmt.Printf(" - Name:		%s\n", feed.Name)
	fmt.Printf(" - Url:			%s\n", feed.Url)
	fmt.Printf(" - UserID:		%s\n", feed.UserID)
	fmt.Printf(" - User:		%s\n", user.Name)
}
