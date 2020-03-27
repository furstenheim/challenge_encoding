package challenge_parser

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)
const SPACE_DELIMITER = "space"

func Parse (t interface{}, reader io.Reader) (error) {
	// buffedReader := bufio.NewReader(reader)
	typeOf := reflect.TypeOf(t)
	nFields := typeOf.NumField()
	fields := make([]*field, nFields)
	for i := 0; i < nFields; i++ {
		structField := typeOf.Field(i)
		tag := structField.Tag
		fmt.Println(structField)
		fmt.Println(tag)
		indexString, ok := tag.Lookup("index")
		if !ok {
			return fmt.Errorf(fmt.Sprintf("Missing index for field %s", structField.Name))
		}
		index, parseError := strconv.ParseInt(indexString, 10, 64)
		if parseError != nil {
			return fmt.Errorf(fmt.Sprintf("Index property could not be parsed for field %s", structField.Name))
		}
		if int(index) >= nFields {
			return fmt.Errorf(fmt.Sprintf("Received index too big for field %s value %d", structField.Name, index))
		}
		if index < 0 {
			return fmt.Errorf(fmt.Sprintf("Received negative index for field %s value %d", structField.Name, index))
		}
		if fields[index] != nil {
			return fmt.Errorf(fmt.Sprintf("Received repeated index for field %s", structField.Name))
		}
		fields[index] = &field{
			index: int(index),
		}
	}
	log.Println(fields)
	return nil
// 	buffedReader.ReadString('\n')
}

type field struct {
	indexedBy string
	index int
	delimiter string
}
