package main
 
import (
  "github.com/codybstrange/blog/internal/config"
  "github.com/codybstrange/blog/internal/database"
  "fmt"
  "os"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
)

type state struct {
  db  *database.Queries
  cfg *config.Config
}

func main() {
  cfg, err := config.Read()
  if err != nil {
    fmt.Printf("Error in Read function: %v", err)
    return
  }
  db, err := sql.Open("postgres", cfg.DBUrl)
  if err != nil {
    log.Fatal(err)
  }

  dbQueries := database.New(db)

  s := &state{db : dbQueries, cfg: &cfg}
  commands := commands {
    handlers: make(map[string]func(*state, command) error),
  }

  commands.register("login",    handlerLogin)
  commands.register("register", handlerRegister)
  commands.register("reset", handlerReset)
  commands.register("users", handlerListUsers)
  commands.register("agg", handlerAgg)
  commands.register("addfeed",  middlewareLoggedIn(handlerAddFeed))
  commands.register("feeds", handlerListFeeds)
  commands.register("follow",  middlewareLoggedIn(handlerFollow))
  commands.register("following",  middlewareLoggedIn(handlerListUserFollows))
  commands.register("unfollow",  middlewareLoggedIn(handlerUnfollow))
  commands.register("browse", middlewareLoggedIn(handlerBrowse))

  args := os.Args[1:]
  if len(args) < 1 {
    log.Fatal("Usage: cli <command> [args...]")
  }
  cmd := command {name: args[0], args: args[1:]}
  if err := commands.run(s, cmd); err != nil {
    log.Fatal(err)
  }

  os.Exit(0)
}
