package challenge_parser

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
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
	buffedReader := bufio.NewReader(reader)

	// 	buffedReader.ReadString('\n')
	typeOf := reflect.TypeOf(t)
	if typeOf.Kind() != reflect.Ptr {
		return fmt.Errorf("top structure must be a pointer")
	}
	if t == nil {
		return fmt.Errorf("input must not be nil")
	}

	if typeOf.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected pointed to struct")
	}
	parser, parseTypeErr := parseType(typeOf.Elem())
	if parseTypeErr != nil {
		return parseTypeErr
	}
	root := &visitor{
		value:      reflect.ValueOf(t).Elem(),
		parser:     parser,
		prev:       nil,
		position:   0,
		isLast:     false,
		isArrayEl:  false,
		isStructEl: false,
	}
	_, err := parseInput(root, buffedReader)
	return err
}

func parseInput (current *visitor, reader *bufio.Reader) (interface{}, error) {
	currentParser := current.parser
	v := current.value
	if kindInSlice(currentParser.kind, basic_types) {
		split := current.findDelimiter()
		text, readErr := reader.ReadString(split)
		text = strings.TrimSuffix(text, string(split))
		if readErr != nil {
			return nil, readErr
		}
		switch currentParser.kind {
		case reflect.String:
			v.SetString(text)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("could not parse %s for kind %s", text, currentParser.kind)
			}
			if v.OverflowInt(n) {
				return nil, fmt.Errorf("overflow for %s with kind %s", text, currentParser.kind)
			}
			v.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(text, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("overflow value with %s for kind %s", text, currentParser.kind)
			}
			if v.OverflowUint(n) {
				return nil, fmt.Errorf("overflow for %s with kind %s", text, currentParser.kind)
			}
			v.SetUint(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(text, v.Type().Bits())
			if err != nil {
				return nil, fmt.Errorf("error parsing float %s for %s of kind %s", err, text, currentParser.kind)
			}
			if v.OverflowFloat(n) {
				return nil, fmt.Errorf("overflow for \"%s\" with kind %s", text, currentParser.kind)
			}
			v.SetFloat(n)
		default:
			panic(fmt.Errorf("not implemented %s", currentParser.kind))
		}
		return nil, nil
	}
	if currentParser.kind == reflect.Slice {
		nv := reflect.MakeSlice(currentParser.ownType, current.nElems,current.nElems)
		for i := 0; i < current.nElems; i++ {
			value := nv.Index(i)
			next := &visitor{
				value:      value,
				parser:     currentParser.elem,
				prev:       current,
				position:   i,
				isLast:     i == current.nElems - 1,
				isArrayEl:  true,
				nElems:     0,
				isStructEl: false,
			}
			_, parseErr := parseInput(next, reader)
			if parseErr != nil {
				return nil, parseErr
			}
			/*_, readErr := reader.ReadString(current.findDelimiter())
			if readErr != nil {
				return nil, readErr
			}*/
		}
		v.Set(nv)
		return nil, nil
	}
	if currentParser.kind == reflect.Struct {
		for i, fieldParser := range (currentParser.fields) {
			var nElems int
			if fieldParser.kind == reflect.Slice {
				nElems = int(current.value.Field(fieldParser.field.indexedBy).Int())
			}
			next := &visitor{
				value:      current.value.Field(i),
				parser:     fieldParser,
				prev:       current,
				position:   i,
				isLast:     i == len(currentParser.fields)-1,
				isArrayEl:  false,
				nElems:     nElems,
				isStructEl: true,
			}
			_, parseErr := parseInput(next, reader)
			if parseErr != nil {
				return nil, parseErr
			}
			/*_, readErr := reader.ReadString(current.findDelimiter())
			if readErr != nil {
				return nil, readErr
			}*/
		}
	}
	return nil, nil
}


func (v * visitor) findDelimiter () byte {
	if v.prev == nil { // last delimiter
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
	value reflect.Value
	parser *typeParser
	prev * visitor
	position int
	isLast bool
	isArrayEl bool
	nElems int
	isStructEl bool
}


func parseType (typeOf reflect.Type) (*typeParser, error) {
	kind := typeOf.Kind()
	if kindInSlice(kind, basic_types) {
		return &typeParser{
			kind: kind,
			ownType: typeOf,
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
			ownType: typeOf,
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
		ownType: typeOf,
		kind: kind,
		fields: make([]*typeParser, nFields),
		nFields: nFields,
	}
	nameToField := map[string]int{}
	for i := 0; i < nFields; i++ {
		structField := typeOf.Field(i)
		name := structField.Name
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
	ownType reflect.Type
	elem *typeParser // in case of array
	fields []*typeParser // in case of slice
	nFields int // in case of slice
	field field // in case of element of struct
}


type field struct {
	name string
	index int
	value reflect.Value
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




