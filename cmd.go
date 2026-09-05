package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/noarrrrr/rss_aggregator/internal/config"
	"github.com/noarrrrr/rss_aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	return c.handlers[cmd.name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Please provide one username")
	}

	bg := context.Background()

	user, err := s.db.GetUser(bg, cmd.args[0])
	handle(err)

	s.cfg.Current_username = user.Name
	err = s.cfg.SetUser()
	if err != nil {
		return err
	}
	fmt.Printf("User '%v' has been logged in\n", cmd.args[0])
	return nil
}

func HandlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Please provide one username")
	}
	bg := context.Background()
	id := uuid.New()

	user, err := s.db.CreateUser(bg, database.CreateUserParams{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	handle(err)

	s.cfg.Current_username = user.Name
	err = s.cfg.SetUser()
	handle(err)

	fmt.Printf("User '%v' has been created and logged in\n", user.Name)
	return nil
}

func HandlerReset(s *state, cmd command) error {
	bg := context.Background()
	err := s.db.Reset(bg)
	handle(err)
	return nil
}

func HandlerUsers(s *state, cmd command) error {
	bg := context.Background()
	users, err := s.db.GetUsers(bg)
	handle(err)
	for _, user := range users {
		if user == s.cfg.Current_username {
			fmt.Printf("* %v (current)\n", user)
		} else {
			fmt.Println("* " + user)
		}
	}
	return nil
}
