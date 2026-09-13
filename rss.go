package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/noarrrrr/rss_aggregator/internal/database"
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
	res, err := http.Get(feedURL)
	if res.StatusCode == 404 {
		return &RSSFeed{}, errors.New("not found")
	}
	handle(err)

	bytes, err := io.ReadAll(res.Body)
	var feed RSSFeed

	err = xml.Unmarshal(bytes, &feed)
	handle(err)

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("This command takes no arguments")
	}
	bg := context.Background()

	feed, err := s.db.NextFeedToFetch(bg)
	handle(err)

	fetchTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	params := database.MarkFeedFetchedParams{
		LastFetchedAt: fetchTime,
		ID:            feed.ID,
	}
	s.db.MarkFeedFetched(bg, params)

	rss, err := fetchFeed(bg, feed.Url)
	handle(err)

	for _, item := range rss.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}
