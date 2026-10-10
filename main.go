package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/paulhoffM/blog/internal/config"
	"github.com/paulhoffM/blog/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Could not exeucte read from main: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal("Failed to open database")
	}
	defer db.Close()
	dbQueries := database.New(db)
	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := &commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowList))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", handlerBrowse)
	if len(os.Args) < 2 {
		log.Fatal("Arguments in command missing")
	}
	cmdStr := command{}
	cmdStr.Name = os.Args[1]
	cmdStr.Args = os.Args[2:]
	err = cmds.run(programState, cmdStr)
	if err != nil {
		log.Fatal(err)
	}

}
