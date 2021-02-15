package builtin

import (
	"fmt"
)

type help struct{}

func (*help) Execute([]string) int {
	fmt.Printf("These shell commands are defined internally. Type `help` to see this list.\n")
	return 0
}
