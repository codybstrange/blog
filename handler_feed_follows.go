package main

import (
  "context"
	"fmt"
	"time"
  "github.com/google/uuid"
  "github.com/codybstrange/blog/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("Not enough arguments in 'follow' command")
  }
  createdAt := time.Now()
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

func handlerListUserFollows(s *state, cmd command, user database.User) error {
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

func handlerUnfollow(s *state, cmd command, user database.User) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("No feed URL provided to unfollow")
  }
  params := database.DeleteFeedFollowParams{
    Name: user.Name,
    Url: cmd.args[0],
  }
  if err := s.db.DeleteFeedFollow(context.Background(), params); err != nil {
    return fmt.Errorf("Error in deleting feed follow: %w", err)
  }
  return nil
}
