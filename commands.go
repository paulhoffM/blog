package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/paulhoffM/blog/internal/database"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(programState *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(programState, cmd)
}

func middlewareLoggedIn(handler func(programState *state, cmd command, user database.User) error) func(*state, command) error {
	return func(programState *state, cmd command) error {
		user, err := programState.db.GetUser(context.Background(), programState.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(programState, cmd, user)
	}
}

func scrapeFeeds(programState *state) error {
	feed, err := programState.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	feedUpdated, err := programState.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}
	rss, err := fetchFeed(context.Background(), feedUpdated.Url)
	if err != nil {
		return err
	}
	for _, item := range rss.Channel.Item {
		timeArg := time.Now()
		NewID := uuid.New()
		post, err := programState.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          NewID,
			CreatedAt:   timeArg,
			UpdatedAt:   timeArg,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: item.PubDate,
			FeedID:      feed.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't store post: %w", err)
		}
		fmt.Printf("Post '%s' was saved", post.Title)
	}
	return nil
}
