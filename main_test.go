package challenge_parser

import (
	"fmt"
	"regexp"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct{
		name string
		input string}{
		{
			"First test",
			`4
				6 10
				8 21
				5 10
				1 5
				6
				2 0 10 5 9 69`,
		},
	}
	for _, tc := range(testCases) {
		input := regexp.MustCompile(`\n\s+`).ReplaceAll([]byte(tc.input), []byte("\n"))
		fmt.Println(string(input))
	}
}
