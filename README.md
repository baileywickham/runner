# runner
CLI runner for golang programs

![Example Run](https://raw.githubusercontent.com/baileywickham/runner/master/runner.png)


## Usage
See `example/` for example program

Support for argument parsing is extreamly limited. Working on that now...

Here is a default configuration:
```golang
import r "github.com/baileywickham/runner"

func main() {
        shell := r.NewShell()
        c1 := r.Command{
                Cmd:      "echo",
                Callback: echo,
                Helptext: "print a string to stdout"}

        c2 := r.Command{
                Cmd:      "addTen",
                Callback: addTen,
                Helptext: "takes and int and adds 10"}

        shell.Add_command(c1, c2) // Add command uses variadic arguments
        shell.Flags() // Parse from cli flags instead of runner
        shell.Start() // Parse from runner, a shell like interface

}

func echo(s string) {
        println(s)
}

func addTen(i int) {
        println(i + 10)
}
```



