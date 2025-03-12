package main

import (
  "fmt"
  "github.com/codybstrange/internal/config"
)

type state struct {
  cfg *Config
}

type command struct {
  name string
  args []string{}
}

type commands struct {
  handlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
  if len(command.args) == 0 {
    return fmt.Errorf("Command must have at least one argument")
  }
  return nil
}

func handlerRegister(s *state, cmd command) error {
  if len(command.args) == 0 {
    return fmt.Errorf("Command must have at least one argument")
  }
  return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
  return
}

func (c *commands) run(s *state, cmd command) error {
  
}
