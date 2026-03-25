package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/SaschaRunge/gator/internal/database"
	"github.com/SaschaRunge/gator/internal/rss"
)

func handlerAgg(s State, cmd command) error {
	rssFeed, err := rss.FetchFeed(context.Background(), FeedURL)
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %w", err)
	}

	fmt.Printf("%v\n", rssFeed)
	return nil
}

func handlerLogin(s State, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument for command. Usage: %s <username>\n", cmd.name)
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
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument for command. Usage: %s <username>\n", cmd.name)
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
