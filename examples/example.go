package main

import r "github.com/baileywickham/runner"

func main() {
	shell := r.NewShell()
	c1 := r.Command{
		Cmd:      "echo",
		Callback: echo,
		Helptext: "echo: print a string to stdout"}

	c2 := r.Command{
		Cmd:      "addTen",
		Callback: addTen,
		Helptext: "addTen: takes an int and adds 10"}
	c3 := r.Command{
		Cmd:      "and",
		Callback: and,
		Helptext: "and: prints boolean and of a and b"}
	shell.Add_command(c1, c2, c3) // Add command uses variadic arguments
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
