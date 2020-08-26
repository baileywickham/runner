package runner

import (
	"bufio"
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
		switch t.Kind().String() {
		case "bool", "Bool":
			arg, err := strconv.ParseBool(strargs[i])
			args[i] = reflect.ValueOf(arg).Convert(t)
			if err != nil {
				goto Err
			}
		case "int", "Int", "Int8", "Int16", "Int32", "Int64":
			arg, err := strconv.ParseInt(strargs[i], 0, 0) // need to change to size
			if err != nil {
				goto Err
			}
			// needs to be changed
			args[i] = reflect.ValueOf(int(arg))
		case "uint", "Uint", "Uint8", "Uint16", "Uint32", "Uint64":
			arg, err := strconv.ParseUint(strargs[i], 0, 0)
			if err != nil {
				goto Err
			}
			args[i] = reflect.ValueOf(arg).Convert(t)
		case "float", "Float32", "Float64":
			arg, err := strconv.ParseFloat(strargs[i], 0)
			if err != nil {
				goto Err
			}
			args[i] = reflect.ValueOf(arg).Convert(t)
		default:
			// panics on failure, probably not good
			args[i] = reflect.ValueOf(strargs[i]).Convert(t)
		}
	}
	method.Call(args)
	return

Err:
	println("ERROR CONVERTING TYPES")
}

func (s *Shell) print_help() {
	for _, command := range s.Commands {
		println(command.Cmd+":", command.Helptext)
	}
}

func NewShell() Shell {
	m := make(map[string]Command)
	c1 := Command{"exit", os.Exit, "exit runner"}
	s := Shell{m}
	s.Add_command(c1)
	s.Add_command(Command{"help", s.print_help, "print list of cammands"})
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
