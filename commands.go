package main

import(
	"errors"
	"context"
	"github.com/paulhoffM/blog/internal/database"
)


type command struct {
	Name string 
	Args []string 
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}


func (c *commands) register(name string, f func(*state, command) error){
	c.registeredCommands[name] = f
	}

func (c *commands) run(programState *state, cmd command) error{
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(programState, cmd)
}

func middlewareLoggedIn(handler func(programState *state, cmd command, user database.User) error) func(*state, command) error{
	return func(programState *state, cmd command) error {
		user, err := programState.db.GetUser(context.Background(), programState.cfg.CurrentUserName)
			if err != nil {
				return err
			}
		return handler(programState, cmd, user)	
	}
}