package ast

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	code  string
	file  *File
	error error
}

var testCases = []testCase{
	{
		code: `func main() {}`,
		file: &File{
			Functions: []*Function{
				{
					Name: "main",
				},
			},
		},
	},
	{
		code: `func main(x Int,) { }`,
		file: &File{
			Functions: []*Function{
				{
					Name: "main",
					Parameters: []*Parameter{{
						Name: "x",
						Type: "Int",
					}},
				},
			},
		},
	},
}

func TestParseString(t *testing.T) {
	for i, tc := range testCases {
		_, got, err := ParseString(tc.code)

		gotJson, _ := json.Marshal(got)
		wantJson, _ := json.Marshal(tc.file)

		if err != tc.error {
			t.Errorf("TestParseString [%d]\nGot:\n%v\nWant:\n%v", i, err, tc.error)
		}

		if !cmp.Equal(got, tc.file) {
			t.Errorf("TestParseString [%d]\nGot:\n%v\nWant:\n%v", i, string(gotJson), string(wantJson))
		}
	}
}
