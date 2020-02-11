package main

import r "github.com/baileywickham/runner"

func main() {
	shell := r.NewShell()
	c1 := r.Command{"echo", echo, "echo: print a string to stdout"}
	c2 := r.Command{"addTen", addTen, "addTen: takes and int and adds 10"}
	shell.Add_command(c1)
	shell.Add_command(c2)
	shell.Start()

}

func echo(s string) {
	println(s)
}

func addTen(i int) {
	println(i + 10)
}
