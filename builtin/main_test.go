package builtin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "",
			want:  []string{""},
		},
		{
			input: "abc def ghij",
			want:  []string{"abc def ghij"},
		},
		{
			input: "abc def ghij kdjfkdj aks kajfkhe sead akkdj3i 3 dkahd .kj kekaldk kdi akkhdh eiz, kdai kdka cmwnandn eiqu a kaizl kkdk eiqon ieiaj di kaekkk kaqqkao zkcj keuaijdk kak zkp, jf euqhtyg akkdj qiiej gjakdjaiwj dk ekaijdka eksk",
			want: []string{
				"abc def ghij kdjfkdj aks kajfkhe sead akkdj3i 3 dkahd .kj kekaldk kdi akkhdh",
				"eiz, kdai kdka cmwnandn eiqu a kaizl kkdk eiqon ieiaj di kaekkk kaqqkao zkcj",
				"keuaijdk kak zkp, jf euqhtyg akkdj qiiej gjakdjaiwj dk ekaijdka eksk",
			},
		},
		{
			input: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want: []string{
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"aaaaaaaaaaaaaaaa",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			got := wrap(test.input)
			assert.Equal(got, test.want)
		})
	}
}
