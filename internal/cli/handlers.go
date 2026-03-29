package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/SaschaRunge/gator/internal/database"
)

func handlerAddFeed(s State, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	name := cmd.args[0]
	url := cmd.args[1]
	ctx := context.Background()

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	_, err := s.dbQueries.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("unable to add feed, possible duplicate? %w", err)
	}

	cmd.args = cmd.args[1:]
	return handlerFollow(s, cmd, user)
}

func handlerAgg(s State, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <period>", cmd.name)
	}

	timePeriod, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("argument <period> could not be parsed into a time value, use e.g, '5s' or '1m': %w", err)
	}

	ticker := time.NewTicker(timePeriod)
	fmt.Printf("Collecting feeds every %v.\n\n", timePeriod)

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println(err)
		}
	}
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

func handlerFollow(s State, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.name)
	}

	ctx := context.Background()

	url := cmd.args[0]
	feed, err := s.dbQueries.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("unable to find feed with url '%s': %w", url, err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.dbQueries.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("User %s is now following feed %s.\n", follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing(s State, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("command %s does not support arguments", cmd.name)
	}

	ctx := context.Background()
	follows, err := s.dbQueries.GetFeedFollowsForUser(ctx, s.config.Current_user_name)
	if err != nil {
		return fmt.Errorf("unable to retrieve followed feeds: %w", err)
	}

	fmt.Printf("You are currently following these feeds:\n")
	for _, follow := range follows {
		fmt.Printf("- %s \n", follow.FeedName)
	}
	fmt.Printf("==========\n")

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

func handlerUnfollow(s State, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	ctx := context.Background()
	url := cmd.args[0]
	feed, err := s.dbQueries.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("unable to find feed with url '%s': %w", url, err)
	}

	deleteFollowParams := database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	return s.dbQueries.DeleteFollow(ctx, deleteFollowParams)
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
