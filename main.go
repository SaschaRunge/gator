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
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error loading config: %s", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		fmt.Printf("Error opening SQL-DB: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	state := cli.NewState(&cfg, dbQueries)
	newCli := cli.New(state)
	if err := newCli.Run(); err != nil {
		fmt.Printf("error running CLI: %s", err)
		os.Exit(1)
	}

	//fmt.Printf("%+v", cfg)
}
