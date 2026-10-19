package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jnsanders1983/blog-aggregator/commands"
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

	s := commands.NewState(&cfg, database.New(db))
	cmds := commands.New()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "a command is required")
		os.Exit(1)
	}

	cmd := commands.Command{Name: os.Args[1], Args: os.Args[2:]}
	if err := cmds.Run(s, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
