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
