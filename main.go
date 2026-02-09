package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Printf("error: db")
		os.Exit(1)
	}
	dbQueries := database.New(db)

	st := &state{
		db:     dbQueries,
		config: cfg,
	}

	cmds := commands{
		commands: make(map[string]func(*state, command) error, 0),
	}

	userInput := os.Args
	if len(userInput) < 2 {
		fmt.Printf("not enough args\n")
		os.Exit(1)
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)

	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	cmd := command{
		name:      userInput[1],
		arguments: userInput[2:],
	}

	err = cmds.run(st, cmd)
	if err != nil {
		// check if it's a unique violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// duplicate URL – ignore
			//do literally nothing
		} else {
			fmt.Printf("cmd run error: %v", err)
		}
	}
}
