package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  *Program
	}{
		{
			input: "cd",
			want: &Program{
				Command: &Command{Name: "cd"},
			},
		},
		{
			input: `     cd
			`,
			want: &Program{
				Command: &Command{Name: "cd"},
			},
		},
		{
			input: ` cd /123`,
			want: &Program{
				Command: &Command{
					Name:      "cd",
					Arguments: []string{"/123"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)
			got, err := Parse("shell", []byte(test.input))
			assert.NoError(err)
			assert.Equal(got, test.want)
		})
	}
}
