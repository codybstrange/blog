package main

import (
  "fmt"
  "context"
  "github.com/google/uuid"
  "time"
  "github.com/codybstrange/blog/internal/database"
)

func handlerLogin(s *state, cmd command) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("Command must have at least one argument")
  }
  if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
    return fmt.Errorf("Cannot get user from database: %w", err)
  }
  fmt.Println("User found")
  if err := s.cfg.SetUser(cmd.args[0]); err != nil {
    return fmt.Errorf("Error in setting username: %w", err)
  }
  fmt.Println("Username has been set")
  return nil
}

func handlerRegister(s *state, cmd command) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("Command must have at least one argument")
  }
  name := cmd.args[0]
  id := uuid.New()
  created_at := time.Now()
  userParams := database.CreateUserParams{
    ID: id,
    CreatedAt: created_at,
    UpdatedAt: created_at,
    Name: name,
  }
  if _, err := s.db.CreateUser(context.Background(), userParams); err != nil {
    return fmt.Errorf("Error in creating user: %w", err)
  }
  fmt.Println("User has been created")
  if err := s.cfg.SetUser(name); err != nil {
    return fmt.Errorf("couldn't set user: %w", err)
  } 
  return nil
}
