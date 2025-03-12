package main

import (
  "fmt"
)

func handlerLogin(s *state, cmd command) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("Command must have at least one argument")
  }
  if err := s.cfg.SetUser(cmd.args[0]); err != nil {
    return fmt.Errorf("Error in setting username: %w", err)
  }
  fmt.Println("Username has been set")
  return nil
}

