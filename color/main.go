package color

import (
	"github.com/fatih/color"
)

var Active = true

func Blue(s string) string {
	if Active {
		return color.BlueString(s)
	}
	return s
}

func Cyan(s string) string {
	if Active {
		return color.BlueString(s)
	}
	return s
}

func Red(s string) string {
	if Active {
		return color.RedString(s)
	}
	return s
}

func Yellow(s string) string {
	if Active {
		return color.YellowString(s)
	}
	return s
}
