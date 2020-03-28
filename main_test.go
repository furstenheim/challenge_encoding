package challenge_parser_test

import (
	"bytes"
	"challenge_parser"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"testing"
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
				2 0 10 5 9 69
`,
				&firstTestFile{},
		},
	}
	for _, tc := range(testCases) {
		input := regexp.MustCompile(`\n\s+`).ReplaceAll([]byte(tc.input), []byte("\n"))
		// fmt.Println(string(input))
		fmt.Println(reflect.TypeOf(tc.parser))
		fmt.Println(reflect.TypeOf(tc.parser).Elem().Field(1).Type.Elem())
		fmt.Println(reflect.ValueOf(reflect.TypeOf(tc.parser).Elem().Field(1).Type.Elem()))
		//reflect.ValueOf(tc.parser).Elem().Field(1).SetCap(3)
		newv := reflect.MakeSlice(reflect.ValueOf(tc.parser).Elem().Field(1).Type(), 3, 3)
		reflect.Copy(newv, reflect.ValueOf(tc.parser).Elem().Field(1))
		reflect.ValueOf(tc.parser).Elem().Field(1).Set(newv)
		fmt.Println(reflect.ValueOf(tc.parser).Elem().Field(1))
		fmt.Println(reflect.ValueOf(tc.parser).Elem().Field(1).Index(0).Field(0))

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
	Queries []query `index:"3" indexed:"NQueries" elem_delimiter:"space"`
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
// TODO test struct -> struct -> elem
// TODO error test on [][]
