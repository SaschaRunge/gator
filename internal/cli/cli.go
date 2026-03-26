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

type State struct {
	config    *config.Config
	dbQueries *database.Queries
}

type command struct {
	name string
	args []string
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

	return cmd.run(c)
}

func (c *command) run(cli *CLI) error {
	callback, exists := cli.commands[c.name]
	if !exists {
		return fmt.Errorf("%s is not a valid command.", c.name)
	}

	return callback(cli.state, *c)
}

func NewState(cfg *config.Config, dbQueries *database.Queries) State {
	return State{cfg, dbQueries}
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

//func (c *command)

func registerCommands() map[string]func(State, command) error {
	return map[string]func(State, command) error{
		"addfeed":   middlewareLoggedIn(handlerAddFeed),
		"agg":       handlerAgg,
		"feeds":     handlerFeeds,
		"follow":    middlewareLoggedIn(handlerFollow),
		"following": handlerFollowing,
		"login":     handlerLogin,
		"register":  handlerRegister,
		"reset":     handlerReset,
		"unfollow":  middlewareLoggedIn(handlerUnfollow),
		"users":     handlerUsers,
	}
}
