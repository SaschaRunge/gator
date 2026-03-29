package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/SaschaRunge/gator/internal/database"
	"github.com/SaschaRunge/gator/internal/rss"
)

func scrapeFeeds(s State) error {
	ctx := context.Background()

	feedItem, err := s.dbQueries.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch from database: %w", err)
	}

	markFeedFetchedParams := database.MarkFeedFetchedParams{
		ID:        feedItem.ID,
		UpdatedAt: time.Now(),
	}
	s.dbQueries.MarkFeedFetched(ctx, markFeedFetchedParams)

	fmt.Printf("\nFetching feeds for %s... .\n", feedItem.Name)
	rssFeed, err := rss.FetchFeed(ctx, feedItem.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf(" - %s\n", item.Title)
	}
	fmt.Println("==========")

	return nil
}

func middlewareLoggedIn(handler func(s State, cmd command, user database.User) error) func(State, command) error {
	return func(s State, cmd command) error {
		ctx := context.Background()
		user, err := s.dbQueries.GetUser(ctx, s.config.Current_user_name)
		if err != nil {
			return fmt.Errorf("current user %s not found in database: %w", s.config.Current_user_name, err)
		}
		return handler(s, cmd, user)
	}
}
