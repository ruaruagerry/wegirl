package gamecfg

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func getFieldInfos(rType reflect.Type, parentIndexChain []int) []*fieldInfo {
	fieldsCount := rType.NumField()
	fieldsList := make([]*fieldInfo, 0, fieldsCount)

	for i := 0; i < fieldsCount; i++ {
		field := rType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		// if the field is an embedded struct, ignore the csv tag
		if field.Anonymous {
			continue
		}

		fieldTag := field.Tag.Get("csv")
		if len(fieldTag) < 1 {
			continue
		}

		var cpy = make([]int, len(parentIndexChain))
		copy(cpy, parentIndexChain)
		indexChain := append(cpy, i)
		fieldInfo := &fieldInfo{IndexChain: indexChain}

		fieldTags := strings.Split(fieldTag, TagSeparator)
		filteredTags := []string{}

		for _, fieldTagEntry := range fieldTags {
			if fieldTagEntry != "omitempty" {
				filteredTags = append(filteredTags, fieldTagEntry)
			} else {
				fieldInfo.omitEmpty = true
			}
		}

		if len(filteredTags) == 1 && filteredTags[0] == "-" {
			continue
		} else if len(filteredTags) > 0 && filteredTags[0] != "" {
			fieldInfo.keys = filteredTags
		} else {
			fieldInfo.keys = []string{field.Name}
		}

		fieldInfo.field = field

		fieldsList = append(fieldsList, fieldInfo)
	}

	return fieldsList
}

func (p *parser) getStructInfo(rType reflect.Type) *structInfo {
	stInfo, ok := p.structs[rType]
	if ok {
		return stInfo
	}

	fieldsList := getFieldInfos(rType, []int{})
	stInfo = &structInfo{}
	fieldsMap := make(map[string]*fieldInfo)
	for _, f := range fieldsList {
		fieldsMap[f.keys[0]] = f
	}
	stInfo.fieldsMap = fieldsMap
	stInfo.fieldsList = fieldsList

	p.structs[rType] = stInfo

	return stInfo
}

func getConcreteReflectValueAndType(in interface{}) (reflect.Value, reflect.Type) {
	value := reflect.ValueOf(in)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value, value.Type()
}

func toString(in interface{}) (string, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return inValue.String(), nil
	case reflect.Bool:
		b := inValue.Bool()
		if b {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%v", inValue.Uint()), nil
	case reflect.Float32:
		return strconv.FormatFloat(inValue.Float(), byte('f'), -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(inValue.Float(), byte('f'), -1, 64), nil
	}
	return "", fmt.Errorf("No known conversion from " + inValue.Type().String() + " to string")
}

func toBool(in interface{}) (bool, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := inValue.String()
		s = strings.TrimSpace(s)
		if strings.EqualFold(s, "yes") {
			return true, nil
		} else if strings.EqualFold(s, "no") || s == "" {
			return false, nil
		} else {
			return strconv.ParseBool(s)
		}
	case reflect.Bool:
		return inValue.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := inValue.Int()
		if i != 0 {
			return true, nil
		}
		return false, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := inValue.Uint()
		if i != 0 {
			return true, nil
		}
		return false, nil
	case reflect.Float32, reflect.Float64:
		f := inValue.Float()
		if f != 0 {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to bool")
}

func toInt(in interface{}) (int64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}
		return strconv.ParseInt(s, 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return inValue.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to int")
}

func toUint(in interface{}) (uint64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}

		// support the float input
		if strings.Contains(s, ".") {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, err
			}
			return uint64(f), nil
		}
		return strconv.ParseUint(s, 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return inValue.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to uint")
}

func toFloat(in interface{}) (float64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}
		return strconv.ParseFloat(s, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return inValue.Float(), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to float")
}

func setField(field reflect.Value, value string, omitEmpty bool) error {
	if field.Kind() == reflect.Ptr {
		if omitEmpty && value == "" {
			return nil
		}
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	switch field.Interface().(type) {
	case string:
		s, err := toString(value)
		if err != nil {
			return err
		}
		field.SetString(s)
	case bool:
		b, err := toBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case int, int8, int16, int32, int64:
		i, err := toInt(value)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case uint, uint8, uint16, uint32, uint64:
		ui, err := toUint(value)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case float32, float64:
		f, err := toFloat(value)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case UNumber:
		u, err := ParseUNumber(value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(u))
	default:
		return nil
	}
	return nil
}
