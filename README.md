## Challenge_encoding

    go get github.com/furstenheim/challenge_encoding
    
Challenge_encoding implements an encoder for code competitions such as [Code jam](https://codingcompetitions.withgoogle.com/codejam) or [Tuenti Challenge](https://contest.tuenti.net/). This avoids having to write a parser for each problem and just focus on the problem.

### Example

A [standard problem](https://contest.tuenti.net/resources/2019/Question_11.html) for a competition would look like 

    2
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

In this case the topic is about space travel. 2 the number of cases. Then we have 2 as the number of moons. For each of those an array with properties. And finally some properties on our ship. Using this encoding it can be summarized as following:

    type tuentiChallengeQuestion11 struct {
        NCases int `index:"0"`
        Cases []tuentiChallengeQuestion11Case `index:"1" indexed:"NCases"`
    }
    
    type tuentiChallengeQuestion11Case struct {
        NMoons int `index:"0"`
        Distances []float64 `index:"1" elem_delimiter:"space" indexed:"NMoons"`
        Positions []float64 `index:"2" elem_delimiter:"space" indexed:"NMoons"`
        Periods []float64 `index:"3" elem_delimiter:"space" indexed:"NMoons"`
        Weights []int `index:"4" elem_delimiter:"space" indexed:"NMoons"`
        Capacity int `index:"5"`
        Range float64 `index:"6"`
    }
    
 In order to parse the input, we only need:
 
    output := &firstTestFile{}
	err := challenge_encoding.Unmarshall(output, bytes.NewReader([]byte(input)))
	

And that would be all.

    fmt.Println(output.NCases)
	> 2
	fmt.Println(len(output.Cases[0].Distances)
	> 2
	fmt.Println(output.Cases[0].Capacity)
	> 20
