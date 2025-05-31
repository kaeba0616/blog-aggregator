package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/kaeba0616/blog-aggregator/internal/config"
	"github.com/kaeba0616/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config....")
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("Error Opening database....")
	}

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.command_register("login", handlerLogin)
	cmds.command_register("register", handlerRegister)
	cmds.command_register("reset", handlerReset)
	cmds.command_register("users", handlerListUsers)
	cmds.command_register("agg", handlerAggregator)
	cmds.command_register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.command_register("feeds", handlerListFeeds)
	cmds.command_register("follow", middlewareLoggedIn(handlerFollow))
	cmds.command_register("following", middlewareLoggedIn(handlerListFollows))
	cmds.command_register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.command_register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.command_run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}
