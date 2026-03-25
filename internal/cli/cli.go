package cli

import (
	"fmt"
	"os"

	"github.com/SaschaRunge/gator/internal/config"
	"github.com/SaschaRunge/gator/internal/database"
)

type CLI struct {
	commands map[string]func(State, command) error
	state    State
}

func New(s State) *CLI {
	return &CLI{
		commands: registerCommands(),
		state:    s,
	}
}

func (c *CLI) Run() error {
	args := os.Args
	if len(args) < 2 {
		return fmt.Errorf("missing command name in input.")
	}

	cmd, err := newCommand(args[1:])
	if err != nil {
		return fmt.Errorf("error parsing command: %w", err)
	}

	return c.runCommand(cmd)
}

func (c *CLI) runCommand(cmd command) error {
	callback, exists := c.commands[cmd.name]
	if !exists {
		return fmt.Errorf("%s is not a valid command.", cmd.name)
	}

	return callback(c.state, cmd)
}

type State struct {
	config    *config.Config
	dbQueries *database.Queries
}

func NewState(cfg *config.Config, dbQueries *database.Queries) State {
	return State{cfg, dbQueries}
}

type command struct {
	name string
	args []string
}

func newCommand(args []string) (command, error) {
	if len(args) == 0 {
		return command{}, fmt.Errorf("no arguments supplied for creation of command.")
	}

	return command{
		name: args[0],
		args: args[1:],
	}, nil
}

func registerCommands() map[string]func(State, command) error {
	return map[string]func(State, command) error{
		"addfeed":  handlerAddFeed,
		"agg":      handlerAgg,
		"login":    handlerLogin,
		"register": handlerRegister,
		"reset":    handlerReset,
		"users":    handlerUsers,
	}
}
