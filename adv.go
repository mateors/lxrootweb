package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type structFieldDetails struct {
	FiledName string
	TagName   string
	Type      string
	OmiyEmpty bool
}
type charInfo struct {
	Index int
	Char  rune
}

func structName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func vMapToJsonStr(vMap map[string]interface{}) string {

	bs, err := json.Marshal(&vMap)
	if err != nil {
		return ""
	}
	return string(bs)
}

func tableToBucket(table string) string {
	return fmt.Sprintf(`%s.%s.%s`, BUCKET_NAME, SCOPE_NAME, table)
}

// helper func
func upperCount(text string) []*charInfo {

	var list []*charInfo
	for i, ch := range text {
		if unicode.IsUpper(ch) {
			list = append(list, &charInfo{i, ch})
		}
	}
	return list
}

// helper func
func splitByUpperCase(text string) []string {

	var list []string
	ci := upperCount(text)
	splitIndex := 0
	for i, c := range ci {
		if i > 0 {
			//fmt.Println(i, text[splitIndex:c.Index]) //c.Index, c.Char, fmt.Sprintf(`%c`, c.Char)
			list = append(list, text[splitIndex:c.Index])
			splitIndex = c.Index
		}
	}
	if splitIndex < len(text) {
		//fmt.Println(name[splitIndex:])
		list = append(list, text[splitIndex:])
	}
	return list
}

func customTableName(structName string) string {

	list := splitByUpperCase(structName)
	for i, part := range list {
		list[i] = strings.ToLower(part)
	}
	return strings.Join(list, "_")
}
