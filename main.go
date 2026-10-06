package main

import (
	"github.com/paulhoffM/blog/internal/config"
	"github.com/paulhoffM/blog/internal/database"
	"log"
	"os"
	"database/sql"
	"fmt"
)
import _ "github.com/lib/pq"

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main () {
	cfg, err := config.Read()
		if err != nil{
			log.Fatalf("Could not exeucte read from main: %v", err)
		}
	fmt.Println(cfg.CurrentUserName)
	db, err := sql.Open("postgres", cfg.DbUrl)
		if err != nil {
			log.Fatal("Failed to open database")
		}
	defer db.Close()
	dbQueries := database.New(db)
	programState := &state{
		db: dbQueries,
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
	cmds.register("addfeed", handlerAddFeed)
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

