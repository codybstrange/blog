package main

import (
  "database/sql"
  "context"
  "github.com/codybstrange/blog/internal/database"
  "fmt"
  "time"
  "github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("Usage: agg <time_between_reqs>")
  }
  interval, err := time.ParseDuration(cmd.args[0])
  if err != nil {
    return fmt.Errorf("Error in parsing duration %w", err)
  }
  fmt.Printf("Collecting feeds every %v\n", interval)
  ticker := time.NewTicker(interval)
  for ; ; <-ticker.C{
    scrapeFeeds(s)
  }

  return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
  if len(cmd.args) != 2 {
    return fmt.Errorf("Incorrect number of arguments to add feed.")
  }
 
  name := cmd.args[0]
  url := cmd.args[1]
  created_at := time.Now()
  params := database.AddFeedParams{
    ID: uuid.New(),
    CreatedAt: created_at,
    UpdatedAt: created_at,
    Name: name,
    Url: url,
    UserID: user.ID,
  }
  feed, err := s.db.AddFeed(context.Background(), params)
  if err != nil {
    return err
  }
  fmt.Println("Added feed to user")
  fmt.Printf("%v\n", feed)
  created_at = time.Now()
  feedFollowParams := database.CreateFeedFollowParams{
    ID: uuid.New(),
    CreatedAt: created_at,
    UpdatedAt: created_at,
    UserID: user.ID,
    FeedID: params.ID,
  }
  if _, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams); err != nil {
    return err
  }
  return nil
}

func handlerListFeeds(s *state, cmd command) error {
  feeds, err := s.db.GetAllFeeds(context.Background())
  if err != nil {
    return err
  }
  for _, f := range feeds {
    fmt.Printf("%v\n", f.Name)
    fmt.Printf("%v\n", f.Url)
    name, err := s.db.GetUserByID(context.Background(), f.UserID)
    if err != nil {
      return err
    }
    fmt.Printf("%v\n", name)
  }
  return nil
}

func scrapeFeeds(s *state) error {
  nextfeed, err := s.db.GetNextFeedToFetch(context.Background())
  if err != nil {
    return fmt.Errorf("Issue with getting the next feed: %w", err)
  }
  err = s.db.MarkFeedFetched(context.Background(),
    database.MarkFeedFetchedParams{
      ID: nextfeed.ID,
      LastFetchedAt: sql.NullTime{
        Time: time.Now(),
        Valid: true,
      },
    })
  if err != nil {
    return fmt.Errorf("Error with marking feed as fetched %w", err)
  }
  feed, err := s.db.GetFeedByURL(context.Background(), nextfeed.Url)
  if err != nil {
    return fmt.Errorf("Error in fetching feed by URL %w", err)
  }
  print(feed)
  return nil
}

func print(f database.Feed) {
  fmt.Println("Feed:")
  fmt.Printf("ID: %v\n", f.ID)
  fmt.Printf("Created At: %v\n", f.CreatedAt)
  fmt.Printf("Updated At: %v\n", f.UpdatedAt)
  fmt.Printf("Name: %v\n", f.Name)
  fmt.Printf("URL: %v\n", f.Url)
  fmt.Printf("User ID: %v\n", f.UserID)
  fmt.Printf("Last Fetched At: %v\n", f.LastFetchedAt)
}
