package ast

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	name  string
	code  string
	file  *File
	error error
}

var testCases = []testCase{
	{
		name: "Empty function",
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
		name: "Function with one parameter",
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
	{
		name: "Function with two parameters",
		code: `func main(int Int, str String,) { }`,
		file: &File{
			Functions: []*Function{
				{
					Name: "main",
					Parameters: []*Parameter{{
						Name: "int",
						Type: "Int",
					}, {Name: "str", Type: "String"}},
				},
			},
		},
	},
}

func TestParseString(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, got, err := ParseString(tc.code)

			gotJson, _ := json.Marshal(got)
			wantJson, _ := json.Marshal(tc.file)

			if err != tc.error {
				t.Errorf("TestParseString \nGot:\n%v\nWant:\n%v", err, tc.error)
			}

			if !cmp.Equal(got, tc.file) {
				t.Errorf("TestParseString \nGot:\n%v\nWant:\n%v", string(gotJson), string(wantJson))
			}
		})
	}
}
