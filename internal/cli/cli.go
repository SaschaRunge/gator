package cli

import (
	"errors"
	"fmt"

	"github.com/SaschaRunge/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	callback, exists := c.commands[cmd.name]
	if !exists {
		errMsg := fmt.Sprintf("%s is not a valid command.", cmd.name)
		return errors.New(errMsg)
	}

	return callback(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		errMsg := fmt.Sprintf("Missing argument for command. Usage: %s <username>", cmd.name)
		return errors.New(errMsg)
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	} else {
		fmt.Printf("Current user: %s", cmd.args[0])
	}

	return nil
}
