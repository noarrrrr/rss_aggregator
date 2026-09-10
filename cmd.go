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

func HandleAggregate(s *state, cmd command) error {
	bg := context.Background()
	feed, err := fetchFeed(bg, "https://www.wagslane.dev/index.xml")
	handle(err)
	fmt.Println(feed)
	return nil
}

func HandleAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("Please provide the name (one word) and the url")
	}
	bg := context.Background()
	user, err := s.db.GetUser(bg, s.cfg.Current_username)
	handle(err)
	userID := user.ID
	name := cmd.args[0]
	url := cmd.args[1]
	params := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    userID,
	}
	feed, err := s.db.AddFeed(bg, params)
	handle(err)
	fmt.Println(feed)
	follow_cmd := command{
		name: "temp",
		args: cmd.args[1:],
	}
	HandleFollow(s, follow_cmd)
	return nil
}

func HandleFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("This command takes no arguments")
	}
	bg := context.Background()
	feeds, err := s.db.ListFeeds(bg)
	handle(err)
	for i, feed := range feeds {
		username, err := s.db.GetNameByID(bg, feed.UserID)
		handle(err)
		fmt.Printf("\nFeed #%d: %v\nURL: %v\nFeed Creater: %v\n\n", i+1, feed.Name, feed.Url, username)
	}
	return nil
}

func HandleFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Please provide a url to follow")
	}

	bg := context.Background()

	user, err := s.db.GetUser(bg, s.cfg.Current_username)
	handle(err)
	user_id := user.ID

	feed, err := s.db.GetFeedByURL(bg, cmd.args[0])
	handle(err)
	feed_id := feed.ID

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user_id,
		FeedID:    feed_id,
	}

	follow_info, err := s.db.CreateFeedFollow(bg, params)
	handle(err)

	fmt.Printf("You (%v) have started following %v\n", follow_info.UserName, follow_info.FeedName)
	return nil
}

func HandleFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("This command takes no arguments")
	}

	bg := context.Background()

	user := s.cfg.Current_username

	followed_feeds, err := s.db.GetFeedFollowsForUser(bg, user)
	handle(err)

	fmt.Printf("Feeds that you (%v) are following:\n", user)

	for _, feed := range followed_feeds {
		fmt.Println(feed)
	}
	return nil
}
