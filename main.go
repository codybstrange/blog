package main

import (
  "github.com/codybstrange/blog-aggregator/internal/config"
  "fmt"
)


func main() {
  cfg, err := config.Read()
  if err != nil {
    fmt.Printf("Error in Read function: %v", err)
    return
  }
  state := state{cfg: &cfg}
  
  cfg, err = config.Read()
  if err != nil {
    fmt.Printf("Error in Read function: %v", err)
    return
  }
  fmt.Printf("db_url: %v\n", cfg.DBUrl)
  fmt.Printf("current_user_name: %v\n", cfg.CurrentUserName)
  
  return
}
