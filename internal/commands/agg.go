package commands

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
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func HandlerAgg(s *State, cmd Command) error {
	testUrl := "https://www.wagslane.dev/index.xml"
	fetchedFeed, err := FetchFeed(context.Background(), testUrl)
	if err != nil {
		return err
	}
	fmt.Println(fetchedFeed)
	return nil
}
func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	client := http.Client{}
	request, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error: couldnt create request: %w", err)
	}
	request.Header.Set("User-Agent", "gator")
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error: couldnt receive response: %w", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error: couldnt read response body: %w", err)
	}
	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("error: couldnt unmarshal xml response: %w", err)
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
