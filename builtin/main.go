package builtin

type Command interface {
	Execute(args []string) int
	Name() string
	ShortHelp() string
}

var builtins = map[string]Command{}

func init() {
	var l = []Command{
		&cd{},
		&exit{},
		&pwd{},
		&export{},
		&help{},
	}
	for _, b := range l {
		builtins[b.Name()] = b
	}
}

func Get(name string) Command {
	return builtins[name]
}
