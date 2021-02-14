package command

import (
	"fmt"
	"os"
	"os/user"
)

func PWD() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get the current working directory: %+v\n", err)
		return
	}
	fmt.Printf("%s\n", wd)
}

func Exit() {
	os.Exit(0)
}

func CD(args []string) {
	dst := ""
	if len(args) == 0 {
		usr, err := user.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to determine home directory: %+v", err)
			return
		}
		dst = usr.HomeDir
	} else {
		dst = args[0]
	}
	if err := os.Chdir(dst); err != nil {
		fmt.Fprintf(os.Stderr, "no such file or directory: %s", dst)
	}
}
