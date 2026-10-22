# Gator

Gator is a command-line RSS feed aggregator. Register users, add and follow feeds, let gator collect posts in the background, and browse them in your terminal.

## Requirements

- [Go](https://go.dev/doc/install) (1.22 or newer)
- [PostgreSQL](https://www.postgresql.org/download/)
- [goose](https://github.com/pressly/goose) for database migrations: `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Install

```bash
go install github.com/jnsanders1983/blog-aggregator@latest
```

This puts a statically compiled binary named `blog-aggregator` in `$(go env GOPATH)/bin` (make sure that is on your `PATH`). Once installed, you don't need the Go toolchain to run it. To use the name `gator`, add an alias:

```bash
alias gator=blog-aggregator
```

Below, `gator` means the installed binary. During development you can use `go run .` instead.

## Set up the database

1. Create a database, e.g. `createdb gator`.
2. Run the migrations from the repo's root:

```bash
goose -dir sql/schema postgres "postgres://USER:PASSWORD@localhost:5432/gator?sslmode=disable" up
```

## Configure

Create `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://USER:PASSWORD@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

`current_user_name` is updated automatically by `register` and `login`.

## Usage

```bash
gator register alice                      # create a user and log in
gator addfeed "Hacker News" https://hnrss.org/newest   # add a feed (and follow it)
gator agg 1m                              # collect feeds every minute (runs until Ctrl+C)
gator browse 10                           # show the 10 latest posts from followed feeds (default 2)
```

Run `agg` in one terminal and `browse` in another.

### Commands

| Command | Description |
| --- | --- |
| `register <name>` | Create a user and set as current |
| `login <name>` | Switch current user |
| `users` | List users |
| `addfeed <name> <url>` | Add a feed and follow it |
| `feeds` | List all feeds |
| `follow <url>` / `unfollow <url>` | Follow or unfollow an existing feed |
| `following` | List feeds the current user follows |
| `agg <duration>` | Continuously fetch feeds, e.g. `30s`, `1m` |
| `browse [limit]` | Show latest posts from followed feeds |
| `reset` | Delete all users (and their data) |
