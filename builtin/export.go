package builtin

import (
	"fmt"
	"os"
)

type export struct{}

func (*export) Execute([]string) int {
	for _, e := range os.Environ() {
		fmt.Printf("%s\n", e)
	}
	return 0
}

func (*export) Name() string { return "export" }

func (*export) ShortHelp() string { return "export" }
