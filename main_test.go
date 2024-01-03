package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

func TestPineapl(t *testing.T) {
	testDir := "./test"
	entries, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(entries))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".pin") {
			fmt.Println(e.Name())

			t.Run(e.Name(), func(t *testing.T) {
				fullInputFileName := path.Join(testDir, e.Name())
				got, err := compile(fullInputFileName)
				if err != nil {
					t.Fatal(err)
				}

				wantFile := strings.Replace(e.Name(), ".pin", ".ir", 1)

				wantBytes, err := os.ReadFile(path.Join(testDir, wantFile))
				if err != nil {
					t.Fatal(err)
				}

				want := string(wantBytes)

				if got != want {
					t.Errorf("TestPineapl \nGot:\n%v\nWant:\n%v", string(got), string(want))
				}
			})
		}
	}
}
