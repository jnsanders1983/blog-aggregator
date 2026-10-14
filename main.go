package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jnsanders1983/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s := state{
		c: &cfg,
	}

	fmt.Printf("%+v\n", *s.c)

	cmds := commands{
		commandHandler: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

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
