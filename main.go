package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/noarrrrr/rss_aggregator/internal/database"

	_ "github.com/lib/pq"
	"github.com/noarrrrr/rss_aggregator/internal/config"
)

func handle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	cfg, err := config.Read()
	handle(err)

	db, err := sql.Open("postgres", cfg.Db_url)
	dbQueries := database.New(db)

	st := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", HandlerRegister)
	cmds.register("reset", HandlerReset)
	cmds.register("users", HandlerUsers)
	cmds.register("agg", HandleAggregate)
	cmds.register("addfeed", HandleAddFeed)
	cmds.register("feeds", HandleFeeds)
	cmds.register("follow", HandleFollow)
	cmds.register("following", HandleFollowing)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	cmd := command{
		args[1],
		args[2:],
	}

	err = cmds.run(&st, cmd)
	handle(err)
}
