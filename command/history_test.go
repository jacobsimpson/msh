package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistory(t *testing.T) {
	assert := assert.New(t)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to get the current working directory: %+v", err)
	}

	h, err := newHistory()
	if err != nil {
		t.Fatalf("unable to initialize history: %+v", err)
	}

	h.add("this")
	h.add("that")
	h.add("theother")

	assert.Equal("theother", h.get())
	assert.Equal("that", h.back())
	assert.Equal("this", h.back())
	assert.Equal(cwd, h.back())
	assert.Equal(cwd, h.back())
	assert.Equal("this", h.forward())

	h.add("notthat")
	assert.Equal("notthat", h.get())
	assert.Equal("notthat", h.forward())
}
