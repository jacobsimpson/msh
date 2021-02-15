package builtin

import (
	"fmt"
	"os"
)

func PWD() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get the current working directory: %+v\n", err)
		return
	}
	fmt.Printf("%s\n", wd)
}
