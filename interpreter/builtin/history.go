package builtin

import (
	"fmt"
	"os"
)

type history struct {
	directories []string
	current     int
}

func newHistory() (*history, error) {
	d, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get the current directory: %+v", err)
	}
	return &history{
		directories: []string{d},
		current:     0,
	}, nil
}

func (h *history) add(path string) {
	if h.current != len(h.directories)-1 {
		h.directories = h.directories[0 : h.current+1]
	}
	h.directories = append(h.directories, path)
	h.current = len(h.directories) - 1
}

func (h *history) get() string {
	return h.directories[h.current]
}

func (h *history) back() string {
	if h.current > 0 {
		h.current--
	}
	return h.directories[h.current]
}

func (h *history) forward() string {
	if h.current < len(h.directories)-1 {
		h.current++
	}
	return h.directories[h.current]
}

func (h *history) String() string {
	result := ""
	for i, d := range h.directories {
		marker := " "
		if i == h.current {
			marker = "*"
		}
		result += fmt.Sprintf("%s %s\n", marker, d)
	}
	return result
}
