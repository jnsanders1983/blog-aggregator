package commands

import (
	"context"
	"fmt"

	"github.com/jnsanders1983/blog-aggregator/internal/database"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	commandHandler map[string]func(*State, Command) error
}

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.db.GetUser(context.Background(), s.c.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func New() *Commands {
	cmds := &Commands{
		commandHandler: make(map[string]func(*State, Command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(addfeed))
	cmds.register("feeds", feeds)
	cmds.register("follow", middlewareLoggedIn(follow))
	cmds.register("following", middlewareLoggedIn(following))
	cmds.register("unfollow", middlewareLoggedIn(unfollow))
	return cmds
}

func (c *Commands) Run(s *State, cmd Command) error {
	f := c.commandHandler[cmd.Name]
	if f == nil {
		return fmt.Errorf("invalid command: %s", cmd.Name)
	}
	return f(s, cmd)
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commandHandler[name] = f
}
