package main

import (
  "fmt"
  "github.com/codybstrange/blog/internal/config"
)

type state struct {
  cfg *config.Config
}

type command struct {
  name string
  args []string
}

type commands struct {
  handlers map[string]func(*state, command) error
}

func handlerRegister(s *state, cmd command) error {
  if len(cmd.args) == 0 {
    return fmt.Errorf("Command must have at least one argument")
  }
  return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
  if _, found := c.handlers[name]; !found {
    c.handlers[name] = f
  } else {
    fmt.Printf("Command %s already registered\n", name)
  }
  return
}

func (c *commands) run(s *state, cmd command) error {
  f, found := c.handlers[cmd.name]
  if !found {
    return fmt.Errorf("Command %s not registered yet", cmd.name)
  }
  return f(s, cmd)
}
