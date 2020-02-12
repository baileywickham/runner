# runner
CLI runner for golang programs
![Example Run](https://raw.githubusercontent.com/baileywickham/runner/master/runner.png)


## Usage
See `example/` for example program

Here is a default configuration:
```golang
import r "github.com/baileywickham/runner"

func main() {
        shell := r.NewShell()
        c1 := r.Command{"echo", echo, "echo: print a string to stdout"}
        c2 := r.Command{"addTen", addTen, "addTen: takes and int and adds 10"}
        shell.Add_command(c1, c2) // Add command uses variadic arguments
        shell.Start()

}

func echo(s string) {
        println(s)
}

func addTen(i int) {
        println(i + 10)
}
```



