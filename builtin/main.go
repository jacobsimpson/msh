package builtin

type Command interface {
	Execute(args []string) int
}

var builtins = map[string]Command{
	"cd":     &cd{},
	"exit":   &exit{},
	"pwd":    &pwd{},
	"export": &export{},
	"help":   &help{},
}

func Get(name string) Command {
	return builtins[name]
}
