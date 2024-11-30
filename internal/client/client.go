package client

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/obanoff/gator/internal/config"
	"github.com/obanoff/gator/internal/database"
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

func FetchFeed(ctx context.Context, feedUL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("user-agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-2xx status code: %d", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("content-type"), "xml") {
		return nil, errors.New("not xml reponse body")
	}

	var feed RSSFeed
	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func ScrapeFeeds(s *config.State) error {
	feed, err := s.DB.Queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.DB.Queries.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	fetched, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, item := range fetched.Channel.Item {
		timestamp := parseDate(item.PubDate)

		err = s.DB.Queries.CreatePost(context.Background(), database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			FeedID: sql.NullInt32{
				Int32: feed.ID,
				Valid: true,
			},
			PublishedAt: timestamp,
		})
		if err != nil {
			s.Logger.Error(err)
		}
	}

	return nil
}

func parseDate(dateStr string) sql.NullTime {
	layouts := []string{
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02",
	}

	var timestamp sql.NullTime
	var err error

	for _, l := range layouts {
		timestamp.Time, err = time.Parse(l, dateStr)
		if err == nil {
			timestamp.Valid = true
			return timestamp
		}
	}

	timestamp.Valid = false
	return timestamp
}
