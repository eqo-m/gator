package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no username provided")
	}
	userName := cmd.args[0]
	ctx := context.Background()

	_, err := s.db.GetUser(ctx, userName)

	if err != nil {
		return fmt.Errorf("user %s not found", userName)
	}

	err = s.config.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to %s\n", userName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no username provided")
	}
	userName := cmd.args[0]
	ctx := context.Background()

	_, err := s.db.GetUser(ctx, userName)
	if err == nil {
		return errors.New("user already exists")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	userID := uuid.New()
	now := time.Now().UTC()

	params := database.CreateUserParams{
		ID:        userID,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      userName,
	}

	_, err = s.db.CreateUser(ctx, params)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.config.SetUser(userName)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User created and logged in: %s\n", userName)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	current := s.config.CurrentUserName

	list, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range list {
		msg := "* " + user
		if user == current {
			msg += " (current)"
		}
		fmt.Println(msg)
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.ResetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to reset users: %w", err)
	}
	return nil
}
