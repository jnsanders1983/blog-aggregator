package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jnsanders1983/blog-aggregator/internal/config"
	"github.com/jnsanders1983/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DataBaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	s := state{
		c:  &cfg,
		db: dbQueries,
	}

	cmds := commands{
		commandHandler: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "a command is required")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(&s, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
