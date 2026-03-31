package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq/pqerror"

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

	count := 0
	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			publishedAt = time.Time{}
		}

		createPostPararms := database.CreatePostsParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       toNullString(item.Title),
			Url:         item.Link,
			Description: toNullString(item.Description),
			PublishedAt: toNullTime(publishedAt),
			FeedID:      feedItem.ID,
		}
		if _, err := s.dbQueries.CreatePosts(ctx, createPostPararms); err != nil {
			dbErr, ok := err.(*pq.Error)
			if !ok || dbErr.Code != "23505" {
				return fmt.Errorf("unexpected error encountered: %w", err)
			}
		} else {
			count++

		}
		//fmt.Printf("failed to create post %s: %s", item.Title, err)
		//fmt.Printf(" - %s\n", item.Title)
	}
	fmt.Printf("Added %d new feeds.\n", count)

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

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func toNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: !t.IsZero()}
}
