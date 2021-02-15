package builtin

type Command interface {
	Execute(args []string) int
}

func Get(name string) Command {
	return nil
}
