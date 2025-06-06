package cli

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/luism2302/gator/internal/config"
	"github.com/luism2302/gator/internal/database"
	"github.com/mmcdole/gofeed"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	Name_to_function map[string]func(s *State, cmd Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if s == nil {
		return fmt.Errorf("nil pointer to state")
	}
	err := c.Name_to_function[cmd.Name](s, cmd)
	return err
}
func (c *Commands) Register(name string, f func(s *State, cmd Command) error) {
	c.Name_to_function[name] = f
}
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("you must provide a username as an argument")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("user %s doesnt exist", cmd.Arguments[0])
	}
	err = s.Cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("logged in as: %s", cmd.Arguments[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("you must provide a username as an argument")
	}
	name := cmd.Arguments[0]
	_, err := s.Db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("user %s already exists", name)
	}
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	newUser, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	err = s.Cfg.SetUser(newUser.Name)
	if err != nil {
		return err
	}
	fmt.Printf("created user %s: %v", name, newUser)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("this command doesnt support arguments")
	}
	err := s.Db.ResetDatabase(context.Background())
	if err != nil {
		return fmt.Errorf("error reseting database: %w", err)
	}
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("this command doesnt support arguments")
	}
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}
	logged_user, err := config.GetLoggedUserName()
	if err != nil {
		return fmt.Errorf("couldnt get logged user: %w", err)
	}
	for _, user := range users {
		if user == logged_user {
			fmt.Printf("* %s (current)\n", user)
			continue
		}
		fmt.Printf("* %s\n", user)
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("you must only provide a duration in string format")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s...\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("you must provide a Name and a URL as arguments")
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
	}
	feed, err := s.Db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}
	fmt.Printf("user %v: created feed %s, with url: %s\n", feed.UserID, feed.Name, feed.Url)
	params_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), params_follow)
	if err != nil {
		return err
	}
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("this command doesnt support arguments")
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Created by: %s\n", feed.UserName)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("you must provide an URL as an argument")
	}
	feed, err := s.Db.GetFeed(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	new_feed_follow, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating new followed feed: %w", err)
	}
	fmt.Printf("Feed: %s\n", new_feed_follow.FeedName)
	fmt.Printf("Followed by: %s\n", new_feed_follow.UserName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("this command doesnt support arguments")
	}
	following, err := s.Db.GetFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for i, followed := range following {
		if i == 0 {
			fmt.Printf("%s follows:\n", followed.UserName)
		}
		fmt.Printf("- %s\n", followed.FeedName)
	}
	return nil
}
func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("you must provide an URL as an argument")
	}
	err := s.Db.DeleteFeedFollow(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("error deleting follow")
	}
	return nil
}

func HandlerBrowse(s *State, cmd Command) error {
	var limit int
	var err error
	if len(cmd.Arguments) == 0 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return err
		}
	}
	posts, err := s.Db.GetPostsForUser(context.Background(), int32(limit))
	if err != nil {
		return err
	}
	for i, _ := range posts {
		posts[i].Title = html.UnescapeString(posts[i].Title)
		posts[i].Url = html.UnescapeString(posts[i].Url)
		fmt.Printf("Title: %s. URL: %s\n", posts[i].Title, posts[i].Url)
	}
	return nil
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.Curr_username)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}

}

func ScrapeFeeds(s *State) error {
	next_feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feed")
	}
	params := database.MarkedFeedFetchedParams{
		ID:            next_feed.ID,
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	err = s.Db.MarkedFeedFetched(context.Background(), params)
	if err != nil {
		return err
	}
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(next_feed.Url)
	if err != nil {
		return err
	}
	for i, item := range feed.Items {
		if i == 0 {
			fmt.Printf("---------------------- %s ----------------------\n", feed.Title)
		}
		params_post := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: *item.PublishedParsed,
			FeedID:      next_feed.ID,
		}
		_, err = s.Db.CreatePost(context.Background(), params_post)
		if err != nil {
			return err
		}
	}
	return nil
}
