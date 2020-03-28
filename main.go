package challenge_parser

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)
const HEADER_DELIMITER = "delimiter"
const HEADER_INDEX = "index"
const HEADER_INDEXED = "indexed"
const HEADER_ELEM_DELIMITER = "elem_delimiter"
const DELIMITER_SPACE = "space"

var basic_types = []reflect.Kind{
	reflect.Bool,
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Uintptr,
	reflect.Float32,
	reflect.Float64,
	reflect.String,
}

func Parse (t interface{}, reader io.Reader) (error) {
	// buffedReader := bufio.NewReader(reader)

	// 	buffedReader.ReadString('\n')
	typeOf := reflect.TypeOf(t)
	_, err := parseType(typeOf)
	return err
}
func parseType (typeOf reflect.Type) (typeParser, error) {
	kind := typeOf.Kind()
	if kindInSlice(kind, basic_types) {
		return typeParser{
			kind: kind,
		}, nil
	}
	if kind == reflect.Slice {
		elem := typeOf.Elem()
		elemParser, elemParseErr := parseType(elem)
		if elemParseErr != nil {
			return typeParser{}, elemParseErr
		}
		parser := typeParser{

			elem: &elemParser,
			kind: kind,
		}
		return parser, nil
	}
	nFields := typeOf.NumField()
	parser := typeParser{
		kind: kind,
		fields: make([]field, nFields),
		nFields: nFields,
	}
	nameToField := map[string]int{}
	for i := 0; i < nFields; i++ {
		structField := typeOf.Field(i)
		name := structField.Name
		fmt.Println(name, structField.Type.Kind())
		index, parseIndexError := parser.parseIndex(structField)
		if parseIndexError != nil {
			return typeParser{}, parseIndexError
		}
		delimiter, delimiterError := parser.parseDelimiter(structField)
		if delimiterError != nil {
			return typeParser{}, delimiterError
		}
		indexedByString, indexedByError := parser.parseIndexedBy(structField)
		if indexedByError != nil {
			return typeParser{}, indexedByError
		}
		elemDelimiter, elemDelimiterError := parser.parseElemDelimiter(structField)
		if elemDelimiterError != nil {
			return typeParser{}, elemDelimiterError
		}
		fieldTypeParser, subParserError := parseType(structField.Type)
		if subParserError != nil {
			return typeParser{}, subParserError
		}
		nameToField[name] = index
		parser.fields[index] = field{
			name: name,
			parser: &fieldTypeParser,
			index: int(index),
			elemDelimiter: elemDelimiter,
			delimiter: delimiter,
			indexedByString: indexedByString,
		}
	}
	for i, f := range(parser.fields) {
		if f.indexedByString != "" {
			index, ok := nameToField[f.indexedByString]
			if !ok {
				return typeParser{}, fmt.Errorf("unknown property indexed %s for field %s", f.indexedByString, f.name)
			}
			if index >= i {
				return typeParser{}, fmt.Errorf("slice has to be indexed by previous property %s", f.name)
			}
			f.indexedBy = index
			parser.fields[i] = f
		}
	}
	log.Println(parser)
	return parser, nil
}

func kindInSlice(a reflect.Kind, list []reflect.Kind) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type typeParser struct {
	kind reflect.Kind
	elem *typeParser
	fields []field
	nFields int
}


type field struct {
	name string
	index int
	parser *typeParser
	delimiter string
	elemDelimiter string // Used for arrays
	indexedByString string
	indexedBy int
}

func (parser typeParser) parseDelimiter (structField reflect.StructField) (string, error) {
	tag := structField.Tag
	delimiterText, ok := tag.Lookup(HEADER_DELIMITER)
	if !ok {
		return "\n", nil
	}
	if delimiterText == DELIMITER_SPACE {
		return " ", nil
	}
	return "", fmt.Errorf("unknown delimiter %s", delimiterText)
}

func (parser typeParser) parseIndexedBy (structField reflect.StructField) (string, error) {
	tag := structField.Tag
	delimitedBy, ok := tag.Lookup(HEADER_INDEXED)
	kind := structField.Type.Kind()
	if ok && kind != reflect.Slice {
		return "", fmt.Errorf("non slice was indexed %s", structField.Name)
	}
	if  kind == reflect.Slice && (!ok || delimitedBy == "") {
		return "", fmt.Errorf("slice should've been indexed and it was not %s", structField.Name)
	}
	return delimitedBy, nil
}

func (parser typeParser) parseElemDelimiter (structField reflect.StructField) (string, error) {
	tag := structField.Tag
	delimitedBy, ok := tag.Lookup(HEADER_ELEM_DELIMITER)
	kind := structField.Type.Kind()
	if ok && kind != reflect.Slice {
		return "", fmt.Errorf("non slice was elem delimited %s", structField.Name)
	}
	if kind == reflect.Slice && (!ok || delimitedBy == "") {
		delimitedBy = "\n"
	}
	if  kind == reflect.Slice && (delimitedBy == DELIMITER_SPACE) {
		delimitedBy = " "
	}
	return delimitedBy, nil
}

func (parser typeParser) parseIndex (structField reflect.StructField) (int, error) {
	tag := structField.Tag
	indexString, okIndex := tag.Lookup(HEADER_INDEX)
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




