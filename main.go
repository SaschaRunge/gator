package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/SaschaRunge/gator/internal/cli"
	"github.com/SaschaRunge/gator/internal/config"
	"github.com/SaschaRunge/gator/internal/database"
)

func main() {
	cfg, _ := config.Read()

	db, err := sql.Open("postgres", cfg.Db_url)
	dbQueries := database.New(db)

	state := cli.NewState(&cfg, dbQueries)
	cmds := cli.NewCommands()

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Missing command name in input. Exiting.\n")
		os.Exit(1)
	}

	cmd, _ := cli.NewCommand(args[1:])
	err = cmds.Run(&state, cmd)
	if err != nil {
		fmt.Printf("Could not run command: '%s'\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", cfg)
}
