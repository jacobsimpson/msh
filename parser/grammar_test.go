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
				Command: &Exec{Name: "cd"},
			},
		},
		{
			input: `     cd
			`,
			want: &Program{
				Command: &Exec{Name: "cd"},
			},
		},
		{
			input: ` cd /123`,
			want: &Program{
				Command: &Exec{
					Name:      "cd",
					Arguments: []string{"/123"},
				},
			},
		},
		{
			input: `ls -l -a`,
			want: &Program{
				Command: &Exec{
					Name:      "ls",
					Arguments: []string{"-l", "-a"},
				},
			},
		},
		{
			input: `echo " a b c d "`,
			want: &Program{
				Command: &Exec{
					Name:      "echo",
					Arguments: []string{" a b c d "},
				},
			},
		},
		{
			input: ` echo ' a b c d '`,
			want: &Program{
				Command: &Exec{
					Name:      "echo",
					Arguments: []string{" a b c d "},
				},
			},
		},
		{
			input: `pwd>output.txt`,
			want: &Program{
				Command: &Exec{
					Name:        "pwd",
					Redirection: &Redirection{Truncate, "output.txt"},
				},
			},
		},
		{
			input: `echo "this is the thing"    >    output.txt`,
			want: &Program{
				Command: &Exec{
					Name:        "echo",
					Arguments:   []string{"this is the thing"},
					Redirection: &Redirection{Truncate, "output.txt"},
				},
			},
		},
		{
			input: `echo "this is the thing"    >>    output.txt`,
			want: &Program{
				Command: &Exec{
					Name:        "echo",
					Arguments:   []string{"this is the thing"},
					Redirection: &Redirection{Append, "output.txt"},
				},
			},
		},
		{
			input: `echo "this is the thing" >& output.txt`,
			want: &Program{
				Command: &Exec{
					Name:        "echo",
					Arguments:   []string{"this is the thing"},
					Redirection: &Redirection{TruncateAll, "output.txt"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)
			got, err := Parse("shell", []byte(test.input))
			assert.NoError(err)
			assert.Equal(test.want, got)
		})
	}
}
