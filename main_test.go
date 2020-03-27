package challenge_parser_test

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"challenge_parser"
)

func TestParse(t *testing.T) {
	testCases := []struct{
		name string
		input string
		parser interface{}
	}{
		{
			"First test",
			`4
				6 10
				8 21
				5 10
				1 5
				6
				2 0 10 5 9 69`,
				firstTestFile{},
		},
	}
	for _, tc := range(testCases) {
		input := regexp.MustCompile(`\n\s+`).ReplaceAll([]byte(tc.input), []byte("\n"))
		fmt.Println(string(input))
		fmt.Println(reflect.TypeOf(tc.parser))
		challenge_parser.Parse(tc.parser, bytes.NewReader(input))
	}
}

type firstTestFile struct {
	NCases int `0`
	Calls []firstTestCall `1 indexed:"NCases"`
	NQueries int `3`
	Queries []query `4 indexed:"NQueries"`

}
type query struct {
	time int `0 delimiter:"space"`
}
type firstTestCall struct {
	Start int `0 delimiter:"space"`
	End int `1`
}
