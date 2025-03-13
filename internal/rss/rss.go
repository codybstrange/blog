package rss

import (
  "encoding/xml"
  "html"
  "io"
  "net/http"
  "context"
  "time"
)

type RSSFeed struct {
  Channel struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    Items       []RSSItem `xml:"item"`
  } `xml:"channel"` 
}

type RSSItem struct {
  Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
  req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
  if err != nil {
    return nil, err
  }

  req.Header.Set("User-Agent", "blog-aggregator")

  client := http.Client{
    Timeout: 10 * time.Second,
  }
  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()
  
  dat, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  feed := RSSFeed{}
  if err := xml.Unmarshal(dat, &feed); err != nil {
    return nil, err
  }

  feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
  feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
  for i, _ := range feed.Channel.Items {
    feed.Channel.Items[i].Title = html.UnescapeString(feed.Channel.Items[i].Title)
    feed.Channel.Items[i].Description = html.UnescapeString(feed.Channel.Items[i].Description)
  }  

  return &feed, nil

}
