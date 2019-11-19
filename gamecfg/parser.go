package gamecfg

import (
	"reflect"
	"strings"
)

const (
	// TagSeparator seperator
	TagSeparator = ","
)

// fieldInfo field
type fieldInfo struct {
	omitEmpty  bool
	IndexChain []int
	keys       []string

	field reflect.StructField
}

// structInfo strct
type structInfo struct {
	fieldsMap  map[string]*fieldInfo
	fieldsList []*fieldInfo
}

// parser ctx
type parser struct {
	// all struct type encounter
	structs map[reflect.Type]*structInfo
}

func newParser() *parser {
	ctx := parser{}
	ctx.structs = make(map[reflect.Type]*structInfo)

	return &ctx
}

func (p *parser) unmarshalRow(header []string, row []string, out interface{}) {
	// construct out's fields info
	outValue, outType := getConcreteReflectValueAndType(out)
	structInfo := p.getStructInfo(outType)

	// parse every column of csv row
	hIdx := 0
	for ; hIdx < len(header); hIdx++ {
		v := row[hIdx]
		if v == "" {
			// empty value, just discard
			continue
		}

		h := header[hIdx]

		if strings.HasPrefix(h, "arr_") {
			// if encounter 'arr_', then extract array name,
			// find array field in 'out', begin to parse array,
			// there are two type of array: struct type, or basic type.
			// arr_awardsParams_13
			h1 := h[4:]
			index := strings.Index(h1, "_")
			h2 := h1[:index]

			fieldInfo, ok := structInfo.fieldsMap[h2]
			if !ok {
				// not found
				continue
			}

			hIdx = p.parseArrayElement(h2, fieldInfo, row, hIdx, outValue)
		} else {
			fieldInfo, ok := structInfo.fieldsMap[h]
			if !ok {
				// not found
				continue
			}

			field := outValue.FieldByIndex(fieldInfo.IndexChain)
			setField(field, v, fieldInfo.omitEmpty)
		}
	}
}

func (p *parser) parseArrayElement(arrayName string,
	fieldInfo *fieldInfo,
	row []string,
	hIdx int,
	outValue reflect.Value) int {

	// arrayElementPtrType: e.g. *AwardsParams
	arrayElementPtrType := fieldInfo.field.Type.Elem()
	var arrayElementType reflect.Type
	var isStructType = false

	// if it is not pointer type, then it is normal built-in type
	if arrayElementPtrType.Kind() == reflect.Ptr {
		arrayElementType = arrayElementPtrType.Elem()
	} else {
		arrayElementType = arrayElementPtrType
	}

	// struct type that exclude UNumber
	if arrayElementType.Kind() == reflect.Struct {
		// need to exclude UNumber type
		// UNumber is struct type, but we don't use it as []*UNumber
		// indeed we use it as a value type, e.g. []UNumber
		var un UNumber
		if arrayElementPtrType != reflect.TypeOf(un) {
			isStructType = true
		}
	}

	array := outValue.FieldByIndex(fieldInfo.IndexChain)
	if array.IsNil() {
		// create a new array
		array = reflect.MakeSlice(reflect.SliceOf(arrayElementPtrType), 0, 32)
	}

	// construct new element and set each field value, now the element is ptr type
	element := reflect.New(arrayElementType)

	if isStructType {
		// get struct info
		elementStructInfo := p.getStructInfo(arrayElementType)
		fieldsList := elementStructInfo.fieldsList

		// now we set each field's value
		n := 0
		for ; hIdx < len(row); hIdx++ {
			v := row[hIdx]
			fieldInfo := fieldsList[n]
			if v != "" {
				// element is ptr type, use Elem() to get value and then
				// we can access FieldByIndex method, exception occurs without Elem()
				field := element.Elem().FieldByIndex(fieldInfo.IndexChain)
				setField(field, v, fieldInfo.omitEmpty)
			}

			n++
			if n == len(fieldsList) {
				// all fields have been set
				break
			}
		}
	} else {
		// set value directly
		v := row[hIdx]
		setField(element, v, true)
	}

	// logrus.Printf("parseArrayElement, %+v", element)
	if isStructType {
		// struct type array is always: []*structType
		// and element is ptr type, thus we can append element directly
		array = reflect.Append(array, element)
	} else {
		// non-struct type array is: []type
		// we need element.Elem() to get value of the ptr, then append to array
		array = reflect.Append(array, element.Elem())
	}

	// reset the new result array to outValue relative field
	outValue.FieldByIndex(fieldInfo.IndexChain).Set(array)

	return hIdx
}
