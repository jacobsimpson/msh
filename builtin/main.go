package builtin

type Command interface {
	Execute(args []string) int
}

var builtins = map[string]Command{
	"cd":   &cd{},
	"exit": &exit{},
}

func Get(name string) Command {
	return builtins[name]
}
