package json

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	case reflect.Map:
		if value.IsNil() {
			return []byte("{}"), nil
		}
		bf.WriteByte('{')

		lv := value.MapKeys()
		for i:= 0;i<len(lv);i++ {
			v := value.MapIndex(lv[i])

			bf.WriteString("\"")
			if lv[i].Kind() == reflect.Int {
				bf.WriteString(strconv.Itoa(lv[i].Interface().(int)))
			} else {
				bf.WriteString(lv[i].Interface().(string))
			}
			bf.WriteString("\"")
			bf.WriteString(":")

			if bs, err := Marshal(v.Interface()); err != nil {
				return nil, err
			} else {
				bf.Write(bs)
			}

			bf.WriteString(",")
		}

		bf.Truncate(len(bf.Bytes()) - 1)

		bf.WriteByte('}')
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

func Unmarshal(data []byte, v interface{}) error {
	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	// 需要对v进行修改，只能传指针
	if typ.Kind() != reflect.Ptr {
		return errors.New("v 必须是指针类型")
	}

	s := string(data)

	s = strings.TrimLeft(s, " ")
	s = strings.TrimRight(s, " ")

	if len(s) == 0 {
		return nil
	}

	// 解析指针
	typ = typ.Elem()
	value = value.Elem()

	switch typ.Kind() {
	case reflect.String:
		if s[0] == '"' && s[len(s) - 1] == '"' {
			value.SetString(s[1:len(s) - 1])
		} else {
			return fmt.Errorf("非法字符串")
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(s); err == nil {
			value.SetBool(b)
		} else {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			value.SetFloat(f)
		} else {
			return err
		}
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
			//strconv.ParseInt (s, 10, 64) 是将字符串s变成10进制
			// 如 "10" = 10
			// strconv.ParseInt (s, 2, 64) 是将字符串s变成2进制
			// 如 "11" = 3
		if f, err := strconv.ParseInt(s, 10, 64); err == nil {
			value.SetInt(f)
		} else {
			return err
		}
	case reflect.Uint,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if f, err := strconv.ParseUint(s, 10, 64); err == nil {
			value.SetUint(f)
		} else {
			return err
		}

	case reflect.Slice:
		if s[0] != '[' || s[len(s) - 1] != ']' {
			return fmt.Errorf("不是slice类型")
		}

		arr := SplitJson(s[1:len(s) - 1])

		if len(arr) > 0 {
			slice := value
			slice.Set(reflect.MakeSlice(typ, len(arr), len(arr))) // 创建一个slice

			for i:=0;i < len(arr) ;i++  {
				eleValue := slice.Index(i)
				eleType := eleValue.Type() // 转化成type
				if eleType.Kind() != reflect.Ptr {
					eleValue=eleValue.Addr() // 转化成指针
				}

				if err := Unmarshal([]byte(arr[i]), eleValue.Interface()); err != nil {
					return err
				}
			}
		}

	case reflect.Struct:
		if s[0] != '{' || s[len(s) - 1 ] != '}' {
			return fmt.Errorf("不是struct类型")
		}

		arr := SplitJson(s[1 : len(s) - 1])
		
		if len(arr) > 0 {
			// get fieldNum
			fieldCount := typ.NumField()
			
			// 建立映射关系tag
			tag2Field := make(map[string]string, fieldCount)

			for i:=0; i< fieldCount ;i++  {
				fieldType := typ.Field(i)
				
				name := fieldType.Name
				
				if len(fieldType.Tag.Get("json")) > 0 {
					name = fieldType.Tag.Get("json")
				}
				
				tag2Field[name] = fieldType.Name
			}

			for _, ele := range arr {
				// 切割：
				brr := strings.SplitN(ele, ":", 2)
				
				tag := strings.Trim(brr[0], " ")

				if tag[0] != '"' || tag[len(tag) - 1] != '"' {
					return fmt.Errorf("非法数据类型%s", tag)
				}

				tag = tag[1:len(tag)-1] // 获得fieldName

				if fieldName, ok := tag2Field[tag]; ok {
					fieldValue := value.FieldByName(fieldName)
					fieldType := fieldValue.Type()
					if fieldType.Kind() != reflect.Ptr {
						fieldValue = fieldValue.Addr()

						if err := Unmarshal([]byte(brr[1]), fieldValue.Interface()); err != nil {
							return err
						}
					} else {
						// 如果是指针的话，需要创建
						newValue := reflect.New(fieldType.Elem())

						if err := Unmarshal([]byte(brr[1]), newValue.Interface()); err != nil {
							return err
						}

						value.FieldByName(fieldName).Set(newValue)
					}
				} else {
					fmt.Printf("字段%s找不到\n", tag)
				}
			}
		}
		
		
	}
	return nil
}