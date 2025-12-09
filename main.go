package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"strings"
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
	mappings map[string]func(*state, command) error
}

func printFeed(feed database.Feed) {
	fmt.Println("Feed Created:")
	fmt.Printf("  ID:         %s\n", feed.ID)
	fmt.Printf("  Name:       %s\n", feed.Name)
	fmt.Printf("  URL:        %s\n", feed.Url)
	fmt.Printf("  User ID:    %s\n", feed.UserID)
	fmt.Printf("  Created At: %v\n", feed.CreatedAt)
	fmt.Printf("  Updated At: %v\n", feed.UpdatedAt)
}

func (s *state) getCurrentUserByName(ctx context.Context) (database.User, error) {
	user, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if errors.Is(err, sql.ErrNoRows) {
		return database.User{}, fmt.Errorf("user %s not found - please login first")
	}
	if err != nil {
		return database.User{}, fmt.Errorf("database error: %w", err)
	}
	return user, nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	_, err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		return fmt.Errorf("failed marking feed: %w", err)
	}
	feed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed, %w", err)
	}
	for _, item := range feed.Channel.Item {
		fmt.Printf(">>%s %s\n", feed.Channel.Title, item.Title)

		publishedTime, err := parsePublishedTimeAt(item.PubDate)
		if err != nil {
			fmt.Printf("failed to parse publish date. skipping post %s : %v\n", item.Title, err)
			continue
		}

		now := time.Now().UTC()
		params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Title: sql.NullString{
				String: item.Title,
				Valid:  item.Title != "",
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			Published: sql.NullTime{
				Time:  *publishedTime,
				Valid: true,
			},
			FeedID: nextFeed.ID,
		}
		_, err = s.db.CreatePost(ctx, params)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
				fmt.Printf("Skipping - post %s already exists \n", item.Title)
				continue
			}
			return fmt.Errorf("failed to save Post : %v", err)
		}
		fmt.Printf("Saved Post %s\n", item.Title)
	}
	return nil
}

func parsePublishedTimeAt(pubDate string) (*time.Time, error) {
	if pubDate == "" {
		return nil, errors.New("pubDate field is empty")
	}
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		"2006-01-02T15:04:05Z",
		time.RFC3339,
	}
	for _, format := range formats {
		t, err := time.Parse(format, pubDate)
		if err == nil {

			return &t, nil
		}
	}
	return nil, fmt.Errorf("unsupported date format: %s", pubDate)
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func (c *commands) run(s *state, cmd command) error {
	run, ok := c.mappings[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return run(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.mappings[name] = f
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}
	dbURL := cfg.DBURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)
	s := &state{db: dbQueries, config: &cfg}
	c := &commands{
		mappings: make(map[string]func(*state, command) error),
	}

	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAggregate)

	c.register("feeds", handlerGetFeeds)

	// using 'middleware'
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowList))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	err = c.run(s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
