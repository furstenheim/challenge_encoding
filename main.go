package challenge_parser

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

func Parse (t interface{}, reader io.Reader) (error) {
	// buffedReader := bufio.NewReader(reader)
	typeOf := reflect.TypeOf(t)
	nFields := typeOf.NumField()
	parser := typeParser{
		fields: make([]field, nFields),
		nFields: nFields,
	}
	nameToField := map[string]int{}
	for i := 0; i < nFields; i++ {
		structField := typeOf.Field(i)
		name := structField.Name
		index, parseIndexError := parser.parseIndex(structField)
		if parseIndexError != nil {
			return parseIndexError
		}
		delimiter, delimiterError := parser.parseDelimiter(structField)
		if delimiterError != nil {
			return delimiterError
		}
		indexedByString, indexedByError := parser.parseIndexedBy(structField)
		if indexedByError != nil {
			return indexedByError
		}
		nameToField[name] = index
		parser.fields[index] = field{
			name: name,
			index: int(index),
			delimiter: delimiter,
			indexedByString: indexedByString,
		}
	}
	for i, f := range(parser.fields) {
		if f.indexedByString != "" {
			index, ok := nameToField[f.indexedByString]
			if !ok {
				return fmt.Errorf("unknown property indexed %s for field %s", f.indexedByString, f.name)
			}
			if index >= i {
				return fmt.Errorf("slice has to be indexed by previous property %s", f.name)
			}
			f.indexedBy = index
			parser.fields[i] = f
		}
	}
	log.Println(parser)
	return nil
// 	buffedReader.ReadString('\n')
}
type typeParser struct {
	fields []field
	nFields int
}


type field struct {
	name string
	indexedBy int
	index int
	delimiter string
	indexedByString string
}

func (parser typeParser) parseDelimiter (structField reflect.StructField) (string, error) {
	tag := structField.Tag
	delimiterText, ok := tag.Lookup("delimiter")
	if !ok {
		return "\n", nil
	}
	if delimiterText == "space" {
		return " ", nil
	}
	return "", fmt.Errorf("unknown delimiter %s", delimiterText)
}

func (parser typeParser) parseIndexedBy (structField reflect.StructField) (string, error) {
	tag := structField.Tag
	delimitedBy, ok := tag.Lookup("indexed")
	kind := structField.Type.Kind()
	if ok && kind != reflect.Slice {
		return "", fmt.Errorf("non slice was indexed %s", structField.Name)
	}
	if  kind == reflect.Slice && (!ok || delimitedBy == "") {
		return "", fmt.Errorf("slice should've been indexed and it was not %s", structField.Name)
	}
	return delimitedBy, nil
}

func (parser typeParser) parseIndex (structField reflect.StructField) (int, error) {
	tag := structField.Tag
	indexString, okIndex := tag.Lookup("index")
	if !okIndex {
		return 0, fmt.Errorf("missing index for field %s", structField.Name)
	}
	index, parseError := strconv.ParseInt(indexString, 10, 64)
	if parseError != nil {
		return 0, fmt.Errorf("index property could not be parsed for field %s", structField.Name)
	}
	if int(index) >= parser.nFields {
		return 0, fmt.Errorf("received index too big for field %s value %d", structField.Name, index)
	}
	if index < 0 {
		return 0, fmt.Errorf("received negative index for field %s value %d", structField.Name, index)
	}
	if parser.fields[index].delimiter != "" {
		return 0, fmt.Errorf("received repeated index for field %s", structField.Name)
	}
	return int(index), nil
}




