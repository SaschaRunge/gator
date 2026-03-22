package cli

import (
	"errors"
	"fmt"

	"github.com/SaschaRunge/gator/internal/config"
	"github.com/SaschaRunge/gator/internal/database"
)

type state struct {
	config    *config.Config
	dbQueries *database.Queries
}

func NewState(cfg *config.Config, dbQueries *database.Queries) state {
	return state{cfg, dbQueries}
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func NewCommand(args []string) (command, error) {
	if len(args) == 0 {
		return command{}, errors.New("No arguments supplied for creation of command.\n")
	}

	if len(args) == 1 {
		return command{
			name: args[0],
			args: []string{},
		}, nil
	} else {
		return command{
			name: args[0],
			args: args[1:],
		}, nil
	}
}

func NewCommands() commands {
	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	return cmds
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func (c *commands) Run(s *state, cmd command) error {
	callback, exists := c.commands[cmd.name]
	if !exists {
		errMsg := fmt.Sprintf("%s is not a valid command.\n", cmd.name)
		return errors.New(errMsg)
	}

	return callback(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		errMsg := fmt.Sprintf("Missing argument for command. Usage: %s <username>\n", cmd.name)
		return errors.New(errMsg)
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	} else {
		fmt.Printf("Current user: %s\n", cmd.args[0])
	}

	return nil
}
