package commands

import (
	"github.com/jnsanders1983/blog-aggregator/internal/config"
	"github.com/jnsanders1983/blog-aggregator/internal/database"
)

type State struct {
	c  *config.Config
	db *database.Queries
}

func NewState(c *config.Config, db *database.Queries) *State {
	return &State{c: c, db: db}
}
