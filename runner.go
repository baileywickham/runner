package runner

import (
	"bufio"
	"os"
	"reflect"
	"strings"
)

type Shell struct {
	Commands map[string]Command
}
type Command struct {
	Cmd      string
	Callback Callback
	Helptext string
}

// dangerous...
type Callback interface{}

func (s *Shell) Add_command(cmds ...Command) {
	for _, c := range cmds {
		if _, ok := s.Commands[c.Cmd]; ok {
			panic("Dupplicate commands not supported")
		}
		s.Commands[c.Cmd] = c
	}
}

func (s *Shell) Start() {
	println("Entering Runner")
	reader := bufio.NewReader(os.Stdin)
	var history []string
	for {
		print(":|: ")
		// readsting error connot occur for my use case
		text, _ := reader.ReadString('\n')
		history = append(history, text)
		tokens := strings.Fields(text)
		if len(tokens) == 0 {
			s.print_help()
			continue
		}
		if tokens[0] == "help" {
			s.print_help()
			continue
		}

		cmd, ok := s.Commands[tokens[0]]
		if !ok {
			println("Command not found")
		}
		s.call_command(cmd, tokens[1:])
	}
}

func (s *Shell) call_command(cmd Command, strargs []string) {
	// Scary
	f := reflect.TypeOf(cmd.Callback)
	method := reflect.ValueOf(cmd.Callback)
	args := make([]reflect.Value, f.NumIn())

	if len(strargs) != f.NumIn() {
		println("Error calling ", cmd.Cmd)
	}

	for i := 0; i < f.NumIn(); i++ {
		//t := f.In(i).Kind()
		args[i] = reflect.ValueOf(strargs[i])
	}
	method.Call(args)
}

func (s *Shell) print_help() {
	for _, command := range s.Commands {
		println(command.Helptext)
	}
}

func NewShell() Shell {
	m := make(map[string]Command)
	return Shell{m}
}