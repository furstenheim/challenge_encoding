package challenge_parser

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
)
const SPACE_DELIMITER = "space"

func Parse (t interface{}, reader io.Reader) {
	buffedReader := bufio.NewReader(reader)
	typeOf := reflect.TypeOf(t)
	for i := 0; i < typeOf.NumField(); i++ {
		structField := typeOf.Field(i)
		fmt.Println(structField)
	}
	buffedReader.ReadString('\n')
}

type field struct {
	indexedBy string
	index int
	delimiter string
}
