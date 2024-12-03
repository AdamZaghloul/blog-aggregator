package main

import (
	"blog-aggregator/internal/database"
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	feed := RSSFeed{}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &feed, err
	}

	req.Header.Add("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return &feed, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return &feed, err
	}

	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &feed, err
	}

	unescapeStrings(&feed)

	return &feed, nil
}

func unescapeStrings(feed *RSSFeed) {

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString((feed.Channel.Title))
		item.Description = html.UnescapeString((feed.Channel.Description))
	}

}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{time.Now(), true},
		UpdatedAt:     time.Now(),
		ID:            nextFeed.ID,
	}

	_, err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("Saving updates from %s\n", feed.Channel.Title)
	fmt.Println()

	for _, item := range feed.Channel.Item {
		fmt.Printf("%s\n", item.Title)
		pubTime, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return err
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{item.Description, true},
			PublishedAt: pubTime,
			FeedID:      nextFeed.ID,
		}

		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				//Do nothing because we will leave the post as one entry in the DB.
			} else {
				return err
			}
		}
	}

	return nil
}
