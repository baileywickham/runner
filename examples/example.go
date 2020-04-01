package main

import r "github.com/baileywickham/runner"

func main() {
	shell := r.NewShell()
	shell.Add_command(r.Command{
		Cmd:      "echo",
		Callback: echo,
		Helptext: "print a string to stdout"},
		r.Command{
			Cmd:      "addTen",
			Callback: addTen,
			Helptext: "takes an int and adds 10"},
		r.Command{
			Cmd:      "and",
			Callback: and,
			Helptext: "prints boolean and of a and b"},
		r.Command{
			Cmd: "hello-world",
			Callback: func() {
				println("Hello World!")
			},
			Helptext: "prints hello world!"})
	shell.Flags()
	shell.Start()

}

func echo(s string) {
	println(s)
}

func addTen(i int) {
	println(i + 10)
}

func and(a, b bool) {
	println(a && b)
}
