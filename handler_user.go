package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/database"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Printf("no username")
		os.Exit(1)
	}
	username := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("error %s", err)
		os.Exit(1)
	}

	if err := s.config.SetUser(username); err != nil {
		return fmt.Errorf("error setting user, handlerLogin")
	}
	fmt.Printf("user set to %s", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no username")
	}
	name := cmd.arguments[0]
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	}
	_, err := s.db.GetUser(context.Background(), userParams.Name)
	if err == nil {
		os.Exit(1)
	}
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Printf("user %s was created successfully\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Printf("reset failed")
		return err
	}
	fmt.Printf("reset successful\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error retrieving users")
		return err
	}

	for _, user := range users {
		name := user.Name
		if name == s.config.CurrentUserName {
			name = name + " (current)"
		}
		fmt.Println(name)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {

	TimeBetweenRequests, err := time.ParseDuration("1m")
	if err != nil {
		return err
	}
	ticker := time.NewTicker(TimeBetweenRequests)
	for ; ; <-ticker.C {
		fmt.Printf("\nfetching...\n\n")
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		fmt.Printf("No arguments. Need [name] and [url]")
		os.Exit(1)
	}
	if len(cmd.arguments) == 1 {
		fmt.Printf("Missing argument.Need [url]")
		os.Exit(1)
	}

	args := cmd.arguments
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      args[0],
		Url:       args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	followRes, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("CreateFeedFollow callback %v", err)
	}

	fmt.Printf("created feed %s for %s", followRes.FeedName, followRes.UserName)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("need url")
	}
	url := cmd.arguments[0]

	feed, err := s.db.GetFeedUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("GetFeedUrl callback %v", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	followRes, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("CreateFeedFollow callback %v", err)
	}

	fmt.Println(followRes.FeedName)
	fmt.Println(followRes.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error %v, Callback GetFeedFollowsForUser", err)
	}

	for _, feed := range feeds {
		fmt.Println(feed.UserName, feed.FeedName)

	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("need url")
	}
	url := cmd.arguments[0]
	feed, err := s.db.GetFeedUrl(context.Background(), url)
	if err != nil {
		return err
	}
	unfollowParams := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.Unfollow(context.Background(), unfollowParams)
	if err != nil {
		return err
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.arguments) > 0 {
		intLimit, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return err
		}
		limit = int32(intLimit)
	}
	getPostParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), getPostParams)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println(post.Description.String)
		fmt.Println(post.Url)

	}
	return nil
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

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	rssfeed, err := fetchFeed(context.Background(), nextFeed)
	if err != nil {
		return err
	}

	feed, err := s.db.GetFeedUrl(context.Background(), nextFeed)
	if err != nil {
		return err
	}

	_, err = s.db.MarkfeedFetched(context.Background(), feed.ID)

	for _, item := range rssfeed.Channel.Item {
		fmt.Println(item.Title)
		fmt.Println(item.PubDate)

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Description, Valid: true},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{},
			FeedID:      feed.ID,
		}

		_, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			// check if it's a unique violation
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				// duplicate URL – ignore
				//do literally nothing
			} else {
				fmt.Printf("Create post error: %v", err)
			}
		}
	}
	return nil
}
