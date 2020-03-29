package challenge_parser_test

import (
	"bytes"
	"challenge_parser"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"regexp"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct{
		name string
		input string
		parser interface{}
		expected interface{}
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
				firstTestFile{
					NCases:   4,
					Calls:    []firstTestCall{{6, 10}, {8, 21}, {5, 10}, {1, 5}},
					NQueries: 6,
					Queries:  []query{{2}, {0}, {10}, {5}, {9}, {69}},
				},
		},
		{
			"Tuenti challenge 9th question",
			`3
					七十二千 OPERATOR 二三百十二 = 四二千百十二
					五十百二 OPERATOR 五十六 = 千十万二百一五六
					千八百二 OPERATOR 七百四千十七 = 千五百三六十五
`,
				&tuentiChallengeQuestion9{},
				tuentiChallengeQuestion9{
					NCases: 3,
					Cases:  []tuentiChallengeQuestion9Case{
						{"七十二千", "OPERATOR", "二三百十二", "=", "四二千百十二"},
						{"五十百二", "OPERATOR", "五十六", "=", "千十万二百一五六"},
						{"千八百二", "OPERATOR", "七百四千十七", "=", "千五百三六十五"},
					},
				},
		},
		{
			"Tuenti challenge 11th question",
			`2
					2
					2.0 2.5
					0.0 3.14
					12.0 100.0
					4 5
					20
					6.0
					4
					0.3 0.4 0.5 0.6
					0.15 0.18 1.15 1.6
					28.8 216.0 27.0 432.0
					1532 770 1250 1630
					3330
					2.0
`,
				&tuentiChallengeQuestion10{},
				tuentiChallengeQuestion10{
					NCases: 2,
					Cases:  []tuentiChallengeQuestion10Case{
						{
							NMoons:    2,
							Distances: []float64{2.0, 2.5},
							Positions: []float64{0.0, 3.14},
							Periods:   []float64{12.0, 100.0},
							Weights:   []int{4, 5},
							Capacity:  20,
							Range:     6.0,
						},
						{
							NMoons:    4,
							Distances: []float64{0.3, 0.4, 0.5, 0.6},
							Positions: []float64{0.15, 0.18, 1.15, 1.6},
							Periods:   []float64{28.8, 216.0, 27.0, 432.0},
							Weights:   []int{1532, 770, 1250, 1630},
							Capacity:  3330,
							Range:     2.0,
						},
					},
				},
		},
	}
	for _, tc := range(testCases) {
		input := regexp.MustCompile(`\n\s+`).ReplaceAll([]byte(tc.input), []byte("\n"))

		err := challenge_parser.Parse(tc.parser, bytes.NewReader(input))
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, tc.expected, reflect.ValueOf(tc.parser).Elem().Interface())
	}
}

func TestParseFromFile (t *testing.T) {
	reader, err := os.Open("./resources/problem11Input.txt")
	if err != nil {
		t.Error(err)
		return
	}
	challenge := &tuentiChallengeQuestion10{}
	err = challenge_parser.Parse(challenge, reader)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 10, challenge.NCases)
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

type tuentiChallengeQuestion9 struct {
	NCases int `index:"0"`
	Cases []tuentiChallengeQuestion9Case `index:"1" indexed:"NCases"`
}

type tuentiChallengeQuestion9Case struct {
	Lhs1 string `index:"0" delimiter:"space"`
	Operator string `index:"1" delimiter:"space"`
	Lhs2 string `index:"2" delimiter:"space"`
	Equal string `index:"3" delimiter:"space"`
	Rhs string `index:"4" delimiter:"space"`
}

type tuentiChallengeQuestion10 struct {
	NCases int `index:"0"`
	Cases []tuentiChallengeQuestion10Case `index:"1" indexed:"NCases"`
}

type tuentiChallengeQuestion10Case struct {
	NMoons int `index:"0"`
	Distances []float64 `index:"1" elem_delimiter:"space" indexed:"NMoons"`
	Positions []float64 `index:"2" elem_delimiter:"space" indexed:"NMoons"`
	Periods []float64 `index:"3" elem_delimiter:"space" indexed:"NMoons"`
	Weights []int `index:"4" elem_delimiter:"space" indexed:"NMoons"`
	Capacity int `index:"5"`
	Range float64 `index:"6"`
}

// TODO test aliases
// TODO test fixed size arrays
// TODO test array separated by spaces are not spaces
// TODO test uint, int64 int32
// TODO test struct -> struct -> elem
// TODO error test on [][]
