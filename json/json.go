package json

import (
	"bytes"
	"container/list"
	"fmt"
	"reflect"
)

func SplitJson(json string) []string  {
	rect := make([]string, 0, 10)
	stack := list.New()

	beginIndex := 0

	for i, r := range json {
		if r == rune('{') || r == rune('[') {
			stack.PushBack(struct {

			}{}) // 随便插入，只是做个栈标示位
		} else if r == rune('}') || r == rune(']') {
			ele := stack.Back() // 获取栈最后的元素然后推出

			if ele != nil {
				stack.Remove(ele)
			}
		} else if r == rune(',') {
			// 如果栈空，则返回
			if stack.Len() == 0 {
				rect = append(rect, json[beginIndex:i])
				beginIndex = i + 1
			}
		}
	}

	rect = append(rect, json[beginIndex:])

	return rect
}

func Marshal(v interface{}) ([]byte, error) {
	value := reflect.ValueOf(v)
	typ := value.Type() // value直接转化成TypeOf
	if typ.Kind() == reflect.Ptr {
		if value.IsNil() {
			return []byte("null"), nil
		} else {
			value = value.Elem()
			typ = typ.Elem()
		}
	}

	bf := bytes.Buffer{}

	switch typ.Kind() {
	case reflect.String:
		return []byte(fmt.Sprintf("\"%s\"", value.String())), nil
	case reflect.Bool:
		return []byte(fmt.Sprintf("%t", value.Bool())), nil
	case reflect.Float32,reflect.Float64:
		return []byte(fmt.Sprintf("%f", value.Float())), nil
	case reflect.Uint,
		 reflect.Uint16,
		 reflect.Uint32,
		 reflect.Uint64,
		 reflect.Int,
		 reflect.Int8,
		 reflect.Int16,
		 reflect.Int32,
		 reflect.Int64:
		 	return []byte(fmt.Sprintf("%v", value.Interface())), nil
	case reflect.Slice:
		if value.IsNil() {
			return []byte("null"), nil
		}

		bf.WriteByte('[')

		if value.Len()>0{
			for i := 0; i<value.Len() ;i++  {

				if bs, err := Marshal(value.Index(i).Interface()); err != nil{
					return nil, err
				} else {
					bf.Write(bs)
					bf.WriteByte(',')
				}
			}

			bf.Truncate(len(bf.Bytes()) - 1)
		}

		bf.WriteByte(']')
		return bf.Bytes(), nil
	case reflect.Struct:
		bf.WriteByte('{')
		if value.NumField() > 0 {
			for i:=0;i<value.NumField() ;i++  {
				filedValue := value.Field(i)
				filedType := typ.Field(i)

				name := filedType.Name
				if len(filedType.Tag.Get("json")) > 0 {
					name = filedType.Tag.Get("json")
				}

				bf.WriteString("\"")
				bf.WriteString(name)
				bf.WriteString("\"")
				bf.WriteString(":")

				if bs, err := Marshal(filedValue.Interface()); err != nil {
					return nil, err
				} else {
					bf.Write(bs)
				}

				bf.WriteString(",")
			}

			bf.Truncate(len(bf.Bytes()) - 1)
		}
		bf.WriteByte('}')

		return bf.Bytes(), nil

	default:
		return []byte("null"), nil
	}
}