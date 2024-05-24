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

// Part of Advance strategy
// func strStructToFieldsType(structName string) (fieldList []*structFieldDetails) {

// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println("was panic, recovered value", r)
// 			fieldList = nil
// 		}
// 	}()

// 	sInstance := makeInstance(structName) //"main.Account"
// 	iVal := reflect.ValueOf(sInstance)
// 	iTypeOf := iVal.Type()

// 	for i := 0; i < iVal.NumField(); i++ {

// 		typeName := iTypeOf.Field(i).Type.String()
// 		fieldName := iTypeOf.Field(i).Name
// 		fieldTag := iTypeOf.Field(i).Tag.Get("json")

// 		var omitFound bool
// 		if strings.Contains(fieldTag, ",") {
// 			omitFound = true
// 		}
// 		if omitFound {
// 			commaFoundAt := strings.Index(fieldTag, ",")
// 			ntag := fieldTag[0:commaFoundAt]
// 			fieldList = append(fieldList, &structFieldDetails{fieldName, ntag, typeName, omitFound})
// 		} else {
// 			fieldList = append(fieldList, &structFieldDetails{fieldName, fieldTag, typeName, omitFound})
// 		}
// 	}
// 	return fieldList
// }

// func strStructToFields(structName string) []string {

// 	fieldList := []string{}
// 	sInstance := makeInstance(structName) //"main.Account"
// 	if sInstance == nil {
// 		return nil
// 	}
// 	iVal := reflect.ValueOf(sInstance)
// 	iTypeOf := iVal.Type()
// 	for i := 0; i < iVal.NumField(); i++ {

// 		fieldName := iTypeOf.Field(i).Name
// 		fieldTag := iTypeOf.Field(i).Tag.Get("json")
// 		if fieldTag == "" {
// 			fieldList = append(fieldList, fieldName)
// 		} else if strings.Contains(fieldTag, ",") {
// 			slc := strings.Split(fieldTag, ",")
// 			fieldList = append(fieldList, slc[0])
// 		} else {
// 			fieldList = append(fieldList, fieldTag)
// 		}
// 	}
// 	return fieldList
// }

// func structValueProcess(structName string, form map[string]interface{}) map[string]interface{} {

// 	var rform = make(map[string]interface{})
// 	fslc := strStructToFieldsType(structName) //
// 	for _, fd := range fslc {

// 		//fmt.Println(i, fd.FiledName, fd.TagName, fd.Type)
// 		val := fmt.Sprintf("%v", form[fd.TagName])
// 		if fd.Type == "int" {
// 			kval, _ := strconv.Atoi(val)
// 			rform[fd.TagName] = kval

// 		} else if fd.Type == "int64" {
// 			kval, _ := strconv.ParseInt(val, 10, 64)
// 			rform[fd.TagName] = kval

// 		} else if fd.Type == "float64" {
// 			kval, _ := strconv.ParseFloat(val, 64)
// 			rform[fd.TagName] = kval

// 		} else if fd.Type == "string" {
// 			rform[fd.TagName] = val

// 		} else {
// 			rform[fd.TagName] = form[fd.TagName]
// 		}
// 	}
// 	return rform
// }

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

func upperCount(text string) []*charInfo {

	var list []*charInfo
	for i, ch := range text {
		if unicode.IsUpper(ch) {
			list = append(list, &charInfo{i, ch})
		}
	}
	return list
}

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

// func readSructColumnsType(i interface{}) []string {

// 	cols := make([]string, 0)
// 	iVal := reflect.ValueOf(i).Elem()
// 	for i := 0; i < iVal.NumField(); i++ {
// 		f := iVal.Field(i)
// 		vtype := f.Kind().String()
// 		cols = append(cols, vtype)
// 	}
// 	return cols
// }
