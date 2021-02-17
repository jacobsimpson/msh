package parser

type Program struct {
	Command Command
}

type Command interface{}

type Exec struct {
	Name      string
	Arguments []string
}

type Type int

const (
	Truncate Type = iota
	TruncateAll
	Append
)

type Redirection struct {
	Type    Type
	Target  string
	Command Command
}

type Pipe struct {
	Src Command
	Dst Command
}

func getRedirection(stdout interface{}) *Redirection {
	if stdout == nil {
		return nil
	}
	a := stdout.([]interface{})
	return a[1].(*Redirection)
}
