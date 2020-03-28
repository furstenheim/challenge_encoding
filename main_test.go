package challenge_parser_test

import (
	"bytes"
	"fmt"
	"log"
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
		err := challenge_parser.Parse(tc.parser, bytes.NewReader(input))
		if err != nil {
			log.Println(err)
		}
	}
}

type firstTestFile struct {
	NCases int `index:"0"`
	Calls []firstTestCall `index:"1" indexed:"NCases"`
	NQueries int `index:"2"`
	Queries []query `index:"3" indexed:"NQueries"`

}
type query struct {
	Time int `index:"0" delimiter:"space"`
}
type firstTestCall struct {
	Start int `index:"0" delimiter:"space"`
	End int `index:"1"`
}

// TODO test aliases
// TODO test fixed size arrays
// TODO test array separated by spaces are not spaces
// TODO test uint, int64 int32
