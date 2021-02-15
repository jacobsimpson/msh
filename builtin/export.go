package builtin

import (
	"fmt"
	"os"
)

func Export() {
	for _, e := range os.Environ() {
		fmt.Printf("%s\n", e)
	}
}
