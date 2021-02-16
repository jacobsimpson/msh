package parser

type Program struct {
	Command *Command
}

type Command struct {
	Name        string
	Arguments   []string
	Redirection *Redirection
}

type Type int

const (
	Truncate Type = iota
)

type Redirection struct {
	Type   Type
	Target string
}

func getRedirection(stdout interface{}) *Redirection {
	if stdout == nil {
		return nil
	}
	a := stdout.([]interface{})
	return a[1].(*Redirection)
}
