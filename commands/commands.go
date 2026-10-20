package commands

import "fmt"

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	commandHandler map[string]func(*State, Command) error
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
	cmds.register("addfeed", addfeed)
	cmds.register("feeds", feeds)
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
