package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (cs *commands) command_register(name string, f func(*state, command) error) {
	cs.registeredCommands[name] = f
}

func (cs *commands) command_run(s *state, cmd command) error {
	f, ok := cs.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
