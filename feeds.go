package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/database"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no argument provided")
	}
	timeDelta, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to parse argument :%w", err)
	}
	ticker := time.NewTicker(timeDelta)
	fmt.Printf("Collecting feeds every %s\n", timeDelta)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println("error scraping feeds:\n", err)
		}
	}

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("must provide feed name and URL")
	}
	ctx := context.Background()
	name := cmd.args[0]
	url := cmd.args[1]

	feedID := uuid.New()
	now := time.Now().UTC()
	userID := user.ID

	params := database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		UserID:    userID,
		Url:       url,
	}

	feed, err := s.db.CreateFeed(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create Feed : %w", err)
	}
	fmt.Printf("Feed Created:\n")
	printFeed(feed)

	//creating a feed_follow record
	followID := uuid.New()
	followParams := database.CreateFeedFollowParams{
		ID:        followID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return fmt.Errorf("failed to follow Feed :%w", err)
	}

	fmt.Printf("Current user : %s\n", user.Name)
	fmt.Printf("Following feed : %s\n", feed.Name)

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("no url provided")
	}
	ctx := context.Background()
	url := cmd.args[0]

	feedFollowID := uuid.New()
	now := time.Now().UTC()

	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("feed %s not found in database", url)
		}
		return fmt.Errorf("database lookup error: %w", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        feedFollowID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to follow Feed :%w", err)
	}

	fmt.Printf("Current user : %s\n", user.Name)
	fmt.Printf("Following feed : %s\n", feed.Name)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	ctx := context.Background()
	list, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("failed to get list :%w", err)
	}
	for _, feed := range list {
		fmt.Printf("Feed: %s Created By: %s URL: %s\n", feed.FeedName, feed.UserName, feed.Url)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("please provide a URL")
	}
	url := cmd.args[0]

	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, url)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error : Feed %s not found", url)
	}
	if err != nil {
		return fmt.Errorf("failed to fetch Feed :%w", err)
	}
	params := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	_, err = s.db.UnfollowFeed(ctx, params)
	if err != nil {
		return fmt.Errorf("unfollowing failed: %w", err)
	}
	return nil
}

func handlerFollowList(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	list, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch followed feeds: %w", err)
	}
	fmt.Printf("User %s is following:\n", user.Name)
	for _, feed := range list {
		fmt.Printf("Feed: %s \n", feed.FeedName)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	limit := 2
	if len(cmd.args) == 1 {
		n, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = n
	}

	params := database.GetPostsUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsUser(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to fetch posts : %w", err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.Title.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("-----")
	}
	return nil
}
