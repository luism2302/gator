package commands

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/luism2302/gator/internal/database"
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
	if len(cmd.Args) < 1 {
		return fmt.Errorf("error: agg command requires timeBetweenReqs argument")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error: unsupported duration. Please use s, m or h as units")
	}
	fmt.Printf("Collecting feeds every %s\n", cmd.Args[0])
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

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

func scrapeFeeds(s *State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error: couldnt get next feed to fetch: %w", err)
	}
	lastFetchedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	markFeedParams := database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: lastFetchedAt,
		ID:            nextFeed.ID,
	}
	err = s.Db.MarkFeedFetched(context.Background(), markFeedParams)
	if err != nil {
		return fmt.Errorf("error: couldnt mark feed as fetched: %w", err)
	}
	fetchedFeed, err := FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	fmt.Printf("======Succesfully fetched feed: %s======\n", fetchedFeed.Channel.Title)
	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf("-\t%s\n", item.Title)
	}
	return nil
}
