package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/google/uuid"
	"github.com/paulhoffM/blog/internal/database"
)


func handlerFollowList(programState *state, cmd command) error {
	ctx := context.Background()
	user, err := programState.db.GetUser(ctx, programState.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't retrieve user: %w", err)
		}
	feedFollows, err := programState.db.GetFeedFollowsForUser(ctx, user.ID)
	for _, feed := range feedFollows{
		fmt.Println(feed.FeedName)
	}
return nil
}

func handlerFollow(programState *state, cmd command) error {
if len(cmd.Args) != 1 {
	log.Fatalf("Provide URL only")
}
ctx := context.Background()
timeArg := time.Now()
user, err := programState.db.GetUser(ctx, programState.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieve user: %w", err)
	}
feed, err := programState.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't retrieve feed: %w", err)
	}
_, err = programState.db.CreateFeedFollow(ctx,  database.CreateFeedFollowParams{
	ID : uuid.New(),
	CreatedAt : timeArg,
	UpdatedAt : timeArg,
	FeedID : feed.ID,
	UserID : user.ID,
})
	if err != nil {
		return fmt.Errorf("couldn't follow: %w", err)
	}
fmt.Println(feed.Name)
fmt.Println(programState.cfg.CurrentUserName)
return nil
}

func handlerListFeeds(programState *state, cmd command) error{
	feeds, err := programState.db.ListFeeds(context.Background())
	if err != nil {
            return fmt.Errorf("couldn't retrieve feeds: %w", err)
        }
    for _, feed := range feeds {
        fmt.Println(feed.FeedName)
        fmt.Println(feed.Url)
        fmt.Println(feed.UserName)
    }
    return nil
}

func handlerAddFeed(programState *state, cmd command) error {
	ctx := context.Background()
	timeArg := time.Now()
	NewID := uuid.New()
	user, err := programState.db.GetUser(ctx, programState.cfg.CurrentUserName)
		if err!= nil{
			log.Fatal("User does not exit")
		}

	if len(cmd.Args) != 2 {
		log.Fatal("Name or URL mssing")
	}

	feed , err := programState.db.CreateFeed(ctx,  database.CreateFeedParams{
			ID : NewID,
			CreatedAt : timeArg,
			UpdatedAt : timeArg,
			Name : cmd.Args[0],
			Url : cmd.Args[1],
			UserID : user.ID,
	})
		if err != nil {
			return fmt.Errorf("couldn't create feed: %w", err)
		}
	
	_ , err = programState.db.CreateFeedFollow(ctx,  database.CreateFeedFollowParams{
			ID : uuid.New(),
			CreatedAt : timeArg,
			UpdatedAt : timeArg,
			FeedID : NewID,
			UserID : user.ID,
	})
		if err != nil {
			return fmt.Errorf("couldn't create feed follow: %w", err)
		}
	

	fmt.Println("Feed created successfully:")
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}


func handlerAgg(programState *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	feed, err := fetchFeed(ctx, feedURL)
		if err != nil{
			log.Fatalf("There is no feed: %v", err)
		}
	fmt.Println(feed)
	return nil
}

func handlerGetUsers(programState *state, cmd command) error {
	list, err := programState.db.List(context.Background())
		if err != nil {
			log.Fatal("List collapsed")	
		}
	for _, name := range list {
		if name == programState.cfg.CurrentUserName{
			{fmt.Printf("%v (current)\n",name)}
		} else {fmt.Println(name)}
	}
	return nil
}

func handlerReset(programState *state, cmd command) error {
	err := programState.db.Reset(context.Background())
		if err != nil {
			log.Fatal("Could not reset")	
		}
	return nil
}

func handlerRegister(programState *state, cmd command) error {
	if len(cmd.Args) != 1 {
		log.Fatal("Login Name is missing or to long")
	}

	ctx := context.Background()
	timeArg := time.Now()

	_ , err := programState.db.CreateUser(ctx,  database.CreateUserParams{
			ID : uuid.New(),
			CreatedAt : timeArg,
			UpdatedAt : timeArg,
			Name : cmd.Args[0],
	})
		if err != nil {
			return fmt.Errorf("couldn't create user: %w", err)
		}
	err = programState.cfg.SetUser(cmd.Args[0])
		if err != nil{
			fmt.Errorf("Could not set username  from main: %w", err)
		}
	return nil
}

func handlerLogin(programState *state, cmd command) error{
	if len(cmd.Args) == 0 {
		log.Fatalf("Login Name is missing")
	}
	ctx := context.Background()
	_, err := programState.db.GetUser(ctx, cmd.Args[0]) 
		if err != nil {
			log.Fatal("User does not exit")
		}

	err = programState.cfg.SetUser(cmd.Args[0])
	if err != nil{
		fmt.Errorf("Could not set username  from main: %w", err)
	}
	fmt.Println("User switched successfully")
	return nil
}

