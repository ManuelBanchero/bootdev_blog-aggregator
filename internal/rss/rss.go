package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		feedURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("rss.go | An error has ocurred trying to create the http request object.\nError: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rss.go | An error has ocurred trying to GET from [%v].\nError: %w", feedURL, err)
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("rss.go | An error has ocurred trying read the response.\nError: %w", err)
	}

	var feed *RSSFeed
	if err := xml.Unmarshal(resBytes, &feed); err != nil {
		return nil, fmt.Errorf("rss.go | An error has ocurred trying to unmarshal the response.\nError: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	return feed, nil
}
