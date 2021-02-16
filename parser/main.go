package parser

type Program struct {
	Command *Command
}

type Command struct {
	Name      string
	Arguments []string
	Stdout    string
}

func getStdout(stdout interface{}) string {
	if stdout == nil {
		return ""
	}
	a := stdout.([]interface{})
	return a[1].(string)
}
