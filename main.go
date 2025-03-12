package main

import (
  "github.com/codybstrange/blog/internal/config"
  "fmt"
  "os"
  "log"
)


func main() {
  cfg, err := config.Read()
  if err != nil {
    fmt.Printf("Error in Read function: %v", err)
    return
  }
  s := &state{cfg: &cfg}
  commands := commands {
    handlers: make(map[string]func(*state, command) error),
  }

  commands.register("login", handlerLogin)

  args := os.Args[1:]
  if len(args) < 2 {
    log.Fatal("Usage: cli <command> [args...]")
    os.Exit(1)
  }
  cmd := command {name: args[0], args: args[1:]}
  if err := commands.run(s, cmd); err != nil {
    log.Fatal(err)
  }

  os.Exit(0)
}
