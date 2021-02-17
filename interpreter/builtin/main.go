package builtin

import (
	"io"
	"strings"
)

var Version string

type Command interface {
	Execute(in io.ReadCloser, out, err io.WriteCloser, args []string) <-chan int
	Name() string
	ShortHelp() string
	LongHelp() string
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

func wrap(content string) []string {
	result := []string{}
	for {
		if len(content) < 80 {
			result = append(result, strings.TrimSpace(content))
			break
		}
		i := strings.LastIndex(content[0:80], " ")
		if i < 1 {
			result = append(result, content[0:80])
			content = content[80:]
			continue
		}
		result = append(result, content[0:i])
		content = content[i+1:]
	}
	return result
}

func done(status int) <-chan int {
	c := make(chan int)
	go func() {
		c <- status
		close(c)
	}()
	return c
}
