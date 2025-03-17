package main

import (
  "database/sql"
  "context"
  "github.com/codybstrange/blog/internal/database"
  "github.com/codybstrange/blog/internal/rss"
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
  feed, err := rss.FeedByURL(context.Background(), nextfeed.Url)
  if err != nil {
    return fmt.Errorf("Error in fetching feed by URL %w", err)
  }
  created_at := time.Now()
  const longForm := "Jan 2, 2006 at 3:04pm (UTC)"
  pubdate, _ := time.Parse(longForm, feed.PubDate)
  params := database.CreatePostParams {
    ID: uuid.New(),
    CreatedAt: created_at,
    UpdatedAt: created_at,
    Title: feed.Title,
    Description: feed.Description,
    Url: feed.Link,
    PublishedAt: pubdate,
    FeedID: nextfeed.ID,
  }
  if _, err := CreatePost(context.Background(), params); err != nil {
    return fmt.Errorf("Error in creating post: %w", err)
  }
  return nil
}

func browse(s *state, cmd command, user database.User) error {
  limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}

}
