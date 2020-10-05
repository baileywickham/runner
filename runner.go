package runner

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"unicode"
)

type Shell struct {
	Commands map[string]Command
}
type Command struct {
	Cmd      string
	Callback Callback
	Helptext string
}

var prompt string = ":|: "

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

func (s *Shell) Flags() {
	if len(os.Args) == 1 {
		println("Must specify command in Flag mode")
		s.print_help()
		return
	}
	cmd, ok := s.Commands[os.Args[1]]
	if !ok {
		println(os.Args[1], "Command not found")
		return
	}
	if len(os.Args) > 2 {
		s.call_command(cmd, os.Args[2:])
	} else {
		s.call_command(cmd, nil)
	}
}

func (s *Shell) Start() {
	println("Entering Runner")
	s.Add_command(Command{"exit", func() { os.Exit(0) }, "exit runner"})
	reader := bufio.NewReader(os.Stdin)
	var history []string // Not yet implimented
	for {
		print(prompt)
		// readsting error connot occur for my use case
		text, _ := reader.ReadString('\n')
		history = append(history, text)

		tokens := parseLine(text)
		if len(tokens) == 0 {
			s.print_help()
			continue
		}

		cmd, ok := s.Commands[tokens[0]]
		if !ok {
			println(tokens[0], "Command not found")
			continue
		}
		s.call_command(cmd, tokens[1:])
	}
}

func (s *Shell) call_command(cmd Command, strargs []string) {
	// Scary
	f := reflect.TypeOf(cmd.Callback)
	method := reflect.ValueOf(cmd.Callback)
	args := make([]reflect.Value, f.NumIn())

	if len(strargs) < f.NumIn() {
		println("Error calling, invalid paramaters", cmd.Cmd)
		return
	}

	for i := 0; i < f.NumIn(); i++ {
		t := f.In(i)
		v, err := convert_types(t, strargs[i])
		if err != nil {
			println(err.Error())
			continue
		}
		args[i] = v
	}
	method.Call(args)
}

func convert_types(to reflect.Type, arg string) (reflect.Value, error) {
	switch to.Kind().String() {
	case "bool", "Bool":
		arg, err := strconv.ParseBool(arg)
		if err != nil {
			goto Err
		}
		return reflect.ValueOf(arg).Convert(to), nil
	case "int", "Int", "Int8", "Int16", "Int32", "Int64":
		arg, err := strconv.ParseInt(arg, 0, 0) // need to change to size
		if err != nil {
			goto Err
		}
		// needs to be changed
		return reflect.ValueOf(int(arg)), nil
	case "uint", "Uint", "Uint8", "Uint16", "Uint32", "Uint64":
		arg, err := strconv.ParseUint(arg, 0, 0)
		if err != nil {
			goto Err
		}
		return reflect.ValueOf(arg).Convert(to), nil
	case "float", "Float32", "Float64":
		arg, err := strconv.ParseFloat(arg, 0)
		if err != nil {
			goto Err
		}
		return reflect.ValueOf(arg).Convert(to), nil
	default:
		// panics on failure, probably not good
		return reflect.ValueOf(arg).Convert(to), nil
	}

Err:
	return reflect.Value{}, errors.New("Connot convert type to argument")
}

func (s *Shell) print_help() {
	for _, command := range s.Commands {
		println(command.Cmd+":", command.Helptext)
	}
}

func NewShell() Shell {
	m := make(map[string]Command)
	s := Shell{m}
	s.Add_command(
		Command{"help", s.print_help, "print list of commands"})
	return s
}

func parseLine(str string) []string {
	parsed := make([]string, 0)
	token := ""
	quoted := false
	for _, char := range str {

		if char == '"' {
			if quoted {
				parsed = append(parsed, token)
			} else {
				quoted = true
				continue
			}
		}

		if unicode.IsSpace(rune(char)) {
			if !quoted && token != "" {
				parsed = append(parsed, token)
				token = ""
				continue
			}
		}
		token = token + string(char)
	}
	return parsed
}
