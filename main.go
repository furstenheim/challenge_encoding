package challenge_parser

import (
	"bufio"
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
	if typeOf.Kind() != reflect.Struct {
		return fmt.Errorf("top structure must be a struct")
	}
	// TODO throw if not struct
	_, err := parseType(typeOf)
	return err
}

func parseInput (current *visitor, reader bufio.Reader) (interface{}, error) {
	currentParser := current.parser
	if kindInSlice(currentParser.kind, basic_types) {
		split := current.findDelimiter()
		reader.ReadString(split)
	}
	return nil, nil
}

func (v * visitor) findDelimiter () byte {
	if v == nil { // last delimiter
		return '\n'
	}
	if v.isStructEl && !v.isLast {
		return v.parser.field.delimiter
	}
	if v.isStructEl && v.isLast {
		return v.prev.findDelimiter()
	}
	if v.isArrayEl && !v.isLast {
		return v.prev.parser.field.elemDelimiter
	}
	if v.isArrayEl && v.isLast {
		return v.prev.findDelimiter()
	}
	panic("Unexpected")
}

type visitor struct {
	parser *typeParser
	prev * visitor
	position int
	isLast bool
	isArrayEl bool
	isStructEl bool
}


func parseType (typeOf reflect.Type) (*typeParser, error) {
	kind := typeOf.Kind()
	if kindInSlice(kind, basic_types) {
		return &typeParser{
			kind: kind,
		}, nil
	}
	if kind == reflect.Slice {
		elem := typeOf.Elem()
		if elem.Kind() == reflect.Slice {
			return nil, fmt.Errorf("doubly nested array not supported since it is not indexed")
		}
		elemParser, elemParseErr := parseType(elem)
		if elemParseErr != nil {
			return &typeParser{}, elemParseErr
		}
		parser := typeParser{
			elem: elemParser,
			kind: kind,
		}
		return &parser, nil
	}
	if kind != reflect.Struct {
		return nil, fmt.Errorf("kind not supported %s", kind)
	}
	nFields := typeOf.NumField()
	parser := typeParser{
		kind: kind,
		fields: make([]*typeParser, nFields),
		nFields: nFields,
	}
	nameToField := map[string]int{}
	for i := 0; i < nFields; i++ {
		structField := typeOf.Field(i)
		name := structField.Name
		fmt.Println(name, structField.Type.Kind())
		index, parseIndexError := parser.parseIndex(structField)
		if parseIndexError != nil {
			return nil, parseIndexError
		}
		delimiter, delimiterError := parser.parseDelimiter(structField)
		if delimiterError != nil {
			return nil, delimiterError
		}
		indexedByString, indexedByError := parser.parseIndexedBy(structField)
		if indexedByError != nil {
			return nil, indexedByError
		}
		elemDelimiter, elemDelimiterError := parser.parseElemDelimiter(structField)
		if elemDelimiterError != nil {
			return nil, elemDelimiterError
		}
		fieldTypeParser, subParserError := parseType(structField.Type)
		if subParserError != nil {
			return nil, subParserError
		}
		nameToField[name] = index
		fieldTypeParser.field = field{
			name: name,
			// parser: &fieldTypeParser,
			index: int(index),
			elemDelimiter: elemDelimiter,
			delimiter: delimiter,
			indexedByString: indexedByString,
		}

		parser.fields[index] = fieldTypeParser
	}
	for i, f := range(parser.fields) {
		if f.field.indexedByString != "" {
			index, ok := nameToField[f.field.indexedByString]
			if !ok {
				return nil, fmt.Errorf("unknown property indexed %s for field %s", f.field.indexedByString, f.field.name)
			}
			if index >= i {
				return nil, fmt.Errorf("slice has to be indexed by previous property %s", f.field.name)
			}
			f.field.indexedBy = index
			parser.fields[i] = f
		}
	}
	log.Println(parser)
	return &parser, nil
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
	elem *typeParser // in case of array
	fields []*typeParser // in case of slice
	nFields int // in case of slice
	field field // in case of element of struct
}


type field struct {
	name string
	index int
	// parser *typeParser
	delimiter byte
	elemDelimiter byte // Used for arrays
	indexedByString string
	indexedBy int
}

func (parser typeParser) parseDelimiter (structField reflect.StructField) (byte, error) {
	tag := structField.Tag
	delimiterText, ok := tag.Lookup(HEADER_DELIMITER)
	if !ok {
		return '\n', nil
	}
	if delimiterText == DELIMITER_SPACE {
		return ' ', nil
	}
	return ' ', fmt.Errorf("unknown delimiter %s", delimiterText)
}

func (parser typeParser) parseIndexedBy (structField reflect.StructField) (string, error) {
	tag := structField.Tag
	indexedBy, ok := tag.Lookup(HEADER_INDEXED)
	kind := structField.Type.Kind()
	if ok && kind != reflect.Slice {
		return "", fmt.Errorf("non slice was indexed %s", structField.Name)
	}
	if  kind == reflect.Slice && (!ok || indexedBy == "") {
		return "", fmt.Errorf("slice should've been indexed and it was not %s", structField.Name)
	}
	return indexedBy, nil
}

func (parser typeParser) parseElemDelimiter (structField reflect.StructField) (byte, error) {
	tag := structField.Tag
	delimitedBy, ok := tag.Lookup(HEADER_ELEM_DELIMITER)
	kind := structField.Type.Kind()
	var delimitedByByte byte
	if ok && kind != reflect.Slice {
		return ' ', fmt.Errorf("non slice was elem delimited %s", structField.Name)
	}
	if kind == reflect.Slice && (!ok || delimitedBy == "") {
		delimitedByByte = '\n'
	}
	if  kind == reflect.Slice && (delimitedBy == DELIMITER_SPACE) {
		delimitedByByte = ' '
	}
	return delimitedByByte, nil
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
	if parser.fields[index] != nil {
		return 0, fmt.Errorf("received repeated index for field %s", structField.Name)
	}
	return int(index), nil
}




