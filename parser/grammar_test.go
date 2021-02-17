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
				Command: &Redirection{
					Truncate,
					"output.txt",
					&Exec{
						Name: "pwd",
					},
				},
			},
		},
		{
			input: `echo "this is the thing"    >    output.txt`,
			want: &Program{
				Command: &Redirection{
					Truncate,
					"output.txt",
					&Exec{
						Name:      "echo",
						Arguments: []string{"this is the thing"},
					},
				},
			},
		},
		{
			input: `echo "this is the thing"    >>    output.txt`,
			want: &Program{
				Command: &Redirection{
					Append,
					"output.txt",
					&Exec{
						Name:      "echo",
						Arguments: []string{"this is the thing"},
					},
				},
			},
		},
		{
			input: `echo "this is the thing" >& output.txt`,
			want: &Program{
				Command: &Redirection{
					TruncateAll,
					"output.txt",
					&Exec{
						Name:      "echo",
						Arguments: []string{"this is the thing"},
					},
				},
			},
		},
		{
			input: `echo "this is the thing" | grep "this"`,
			want: &Program{
				Command: &Pipe{
					&Exec{
						Name:      "echo",
						Arguments: []string{"this is the thing"},
					},
					&Exec{
						Name:      "grep",
						Arguments: []string{"this"},
					},
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
