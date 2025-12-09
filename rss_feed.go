package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

var httpClient = &http.Client{}

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

func (feed *RSSFeed) unescape() {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].unescape()
	}
}

func (item *RSSItem) unescape() {
	item.Title = html.UnescapeString(item.Title)
	item.Link = html.UnescapeString(item.Link)
	item.Description = html.UnescapeString(item.Description)
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var unmarshalledRSS RSSFeed

	//Create Request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	//Execute request with Client.Do
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	defer resp.Body.Close()

	//Check the HTTP status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received non- successful status code: %d", resp.StatusCode)
	}

	//Read body into byte slice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	err = xml.Unmarshal(body, &unmarshalledRSS)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	unmarshalledRSS.unescape()
	return &unmarshalledRSS, nil

}
