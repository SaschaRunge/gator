package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/SaschaRunge/gator/internal/database"
	"github.com/SaschaRunge/gator/internal/rss"
)

func handlerAddFeed(s State, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	name := cmd.args[0]
	url := cmd.args[1]
	ctx := context.Background()

	user, err := s.dbQueries.GetUser(ctx, s.config.Current_user_name)
	if err != nil {
		return fmt.Errorf("unable to add feed, user %s not found in database: %w", s.config.Current_user_name, err)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	_, err = s.dbQueries.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("unable to add feed, possible duplicate? %w", err)
	}

	return nil
}

func handlerAgg(s State, cmd command) error {
	rssFeed, err := rss.FetchFeed(context.Background(), FeedURL)
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %w", err)
	}

	fmt.Printf("%v\n", rssFeed)
	return nil
}

func handlerFeeds(s State, cmd command) error {
	feeds, err := s.dbQueries.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.dbQueries.GetUserFromUUID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("no matching uuid found in table users: %w", err)
		}
		fmt.Printf("name: %s | url: %s | creator: %s \n", feed.Name, feed.Url, user.Name)
	}

	return nil
}
func handlerLogin(s State, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.name)
	}

	if _, err := s.dbQueries.GetUser(context.Background(), cmd.args[0]); err != nil {
		return fmt.Errorf("user '%s' does not exist!", cmd.args[0])
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	} else {
		fmt.Printf("Current user: %s\n", cmd.args[0])
	}

	return nil
}

func handlerRegister(s State, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.name)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	_, err := s.dbQueries.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	} else {
		fmt.Printf("Current user: %s\n", cmd.args[0])
	}

	return nil
}

func handlerReset(s State, cmd command) error {
	if err := s.dbQueries.ResetUsers(context.Background()); err != nil {
		return err
	}

	fmt.Println("Reset table 'users'.")
	return nil
}

func handlerUsers(s State, cmd command) error {
	users, err := s.dbQueries.GetUsers(context.Background())
	if err != nil {
		return err
	}

	currentUser := s.config.Current_user_name
	for _, user := range users {
		msg := fmt.Sprintf("* %s", user.Name)
		if user.Name == currentUser {
			msg = fmt.Sprintf("%s (current)", msg)
		}
		fmt.Println(msg)
	}
	return nil
}
