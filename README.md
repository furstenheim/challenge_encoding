## Challenge_encoding

    go get github.com/furstenheim/challenge_encoding
    
Challenge_encoding implements an encoder for code competitions such as [Code jam](https://codingcompetitions.withgoogle.com/codejam) or [Tuenti Challenge](https://contest.tuenti.net/). This avoids having to write a parser for each problem and just focus on the problem.

### Example

A [standard problem](https://contest.tuenti.net/resources/2019/Question_11.html) for a competition would look like 

    1
    2
    2.0 2.5
    0.0 3.14
    12.0 100.0
    4 5
    20
    6.0
    

In this case the topic is about space travel. "1" the number of cases. Then we have "2" as the number of moons of the first case. For each of those an array with properties. And finally some properties on our ship. Using this encoding it can be summarized as following:

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
	> 1
	fmt.Println(len(output.Cases[0].Distances)
	> 2
	fmt.Println(output.Cases[0].Capacity)
	> 20


### Supported annotations

#### index
Index describes in what position of the input the property will be received. For example in:

    type example struct {
        IComeSecond int `index:"1"`
        IComeFirst int `index:"0"`
    }
Index is required in all exported properties.

#### delimiter
How a property finishes. By default it is assumed to be a newline. Possible value is "space"

    type spaceDelimited struct {
        First int `index:"0" delimiter:"space"`
        Second int `index:"1"
        Third int
    }
    input := `1 2
    3`
    parsed := spaceDelimited{First: 1, Second: 2, Third: 3}
    
#### indexed
All variable size slices are assumed to be indexed by another property. That is, there is another field that specifies the length of the slice.

    type sliceExample struct {
        LengthOfSlice int `index:"0"`
        Slice []int `index:"1" indexed:"LengthOfSlice"` 
    }
    input := `2
    1
    3`
    parsed := sliceExample{LengthOfSlice: 2, Slice: []int{1, 3}}
    
#### elem_delimiter
Given a slice we need to know how to elements of it are delimited. This is defined by this property. Default value is new line. It can be "space"

    type sliceExample struct {
        LengthOfSlice int `index:"0"`
        NewLineSlice []int `index:"1" indexed:"LengthOfSlice"`
        SpaceLineSlice []int `index:"1" indexed:"LengthOfSlice" elem_delimiter:"space"` 
    }
    input := `2
    1
    3
    4 5`
    parsed := sliceExample{LengthOfSlice: 2, NewLineSlice: []int{1, 3}, SpaceLineSlice: []int{4, 5}}

#### Not currently supported
* Fixed size arrays. Instead one can use struct with the given number of properties
* Not indexed slices. Most of the cases in the competitions include a field with the length of variable size arrays, so this is not supported.
* Nested slices. Like in:

        type example struct {
            Matrix [][]int
        }
        
For that we would need to have two fields indexing the different lengths. Something like elem_delimiter_2
