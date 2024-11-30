package commands

import (
	"errors"

	"github.com/obanoff/gator/internal/config"
)

type CommandHandler func(s *config.State, args []string) error

type Commands struct {
	commands map[string]CommandHandler
}

func (c *Commands) Register(name string, f CommandHandler) {
	c.commands[name] = f
}

func (c *Commands) Run(s *config.State, args []string) error {
	if len(args) < 1 {
		return errors.New("arguments not provided")
	}

	handler, ok := c.commands[args[0]]
	if !ok {
		return errors.New("invalid command")
	}

	return handler(s, args)
}

func NewCommands() Commands {
	return Commands{
		commands: make(map[string]CommandHandler),
	}
}
