package main

import (
  "context"
  "github.com/codybstrange/blog/internal/rss"
  "github.com/codybstrange/blog/internal/database"
  "fmt"
  "time"
  "github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
  feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
  if err != nil {
    return err
  }
  fmt.Printf("%v", feed)
  return nil
}

func handlerAddFeed(s *state, cmd command) error {
  if len(cmd.args) != 2 {
    return fmt.Errorf("Incorrect number of arguments to add feed.")
  }
  user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
  if err != nil {
    return err
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

func handlerFollow(s *state, cmd command) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("Not enough arguments in follow command")
  }
  createdAt := time.Now()
  user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
  if err != nil {
    return err
  }
  url := cmd.args[0]
  feed, err := s.db.GetFeedByURL(context.Background(), url)
  params := database.CreateFeedFollowParams{
    ID: uuid.New(),
    CreatedAt: createdAt,
    UpdatedAt: createdAt,
    UserID: user.ID,
    FeedID: feed.ID,
  }
  if _, err = s.db.CreateFeedFollow(context.Background(), params); err != nil {
    return err
  }
  fmt.Printf("Now following - Feed Name: %v; User: %v\n", user.Name, feed.Name)
  return nil
}

func handlerListUserFollows(s *state, cmd command) error {
  feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
  if err != nil {
    return err
  }
  fmt.Printf("User %v is following feeds:\n", s.cfg.CurrentUserName)
  for _, feedFollow := range feedFollows {
    fmt.Printf("%v\n", feedFollow.FeedName)
  }
  return nil
}
