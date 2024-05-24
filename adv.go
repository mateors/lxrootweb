package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/rs/xid"
)

type structFieldDetails struct {
	FiledName string
	TagName   string
	Type      string
}

type charInfo struct {
	Index int
	Char  rune
}

type mapStringScan struct {
	cp       []interface{}
	row      map[string]string
	colCount int
	colNames []string
}

func newMapStringScan(columnNames []string) *mapStringScan {
	lenCN := len(columnNames)
	s := &mapStringScan{
		cp:       make([]interface{}, lenCN),
		row:      make(map[string]string, lenCN),
		colCount: lenCN,
		colNames: columnNames,
	}
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}

func (s *mapStringScan) Update(rows *sql.Rows) error {

	if err := rows.Scan(s.cp...); err != nil {
		return err
	}

	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.row[s.colNames[i]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return fmt.Errorf("cannot convert index %d column %s to type *sql.RawBytes", i, s.colNames[i])
		}
	}
	return nil
}

func (s *mapStringScan) Get() map[string]string {
	return s.row
}

type stringStringScan struct {
	cp       []interface{}
	row      []string
	colCount int
	colNames []string
}

func newStringStringScan(columnNames []string) *stringStringScan {
	lenCN := len(columnNames)
	s := &stringStringScan{
		cp:       make([]interface{}, lenCN),
		row:      make([]string, lenCN*2),
		colCount: lenCN,
		colNames: columnNames,
	}
	j := 0
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
		s.row[j] = s.colNames[i]
		j = j + 2
	}
	return s
}
func (s *stringStringScan) Update(rows *sql.Rows) error {
	if err := rows.Scan(s.cp...); err != nil {
		return err
	}
	j := 0
	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.row[j+1] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return fmt.Errorf("Cannot convert index %d column %s to type *sql.RawBytes", i, s.colNames[i])
		}
		j = j + 2
	}
	return nil
}

func (s *stringStringScan) Get() []string {
	return s.row
}

func registerType(emptyStruct interface{}) {
	//fmt.Println(reflect.TypeOf(emptyStruct).String()) //main.Account
	//fmt.Println(reflect.TypeOf(emptyStruct).Name())   //Account
	typeRegistry[reflect.TypeOf(emptyStruct).Name()] = reflect.TypeOf(emptyStruct)
}
func makeInstance(name string) interface{} {
	return reflect.New(typeRegistry[name]).Elem().Interface()
}

// Part of Advance strategy
func strStructToFieldsType(structName string) (fieldList []*structFieldDetails) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("was panic, recovered value", r)
			fieldList = nil
		}
	}()

	sInstance := makeInstance(structName) //"main.Account"
	iVal := reflect.ValueOf(sInstance)
	iTypeOf := iVal.Type()
	//typeOfS := reflect.ValueOf(sInstance).Type()
	//fmt.Println(reflect.TypeOf(sInstance).Kind().String())

	for i := 0; i < iVal.NumField(); i++ {

		typeName := iTypeOf.Field(i).Type.String()
		fieldName := iTypeOf.Field(i).Name
		fieldTag := iTypeOf.Field(i).Tag.Get("json")

		var omitFound bool
		if strings.Contains(fieldTag, ",") {
			omitFound = true
		}
		if omitFound {
			commaFoundAt := strings.Index(fieldTag, ",")
			ntag := fieldTag[0:commaFoundAt]
			fieldList = append(fieldList, &structFieldDetails{fieldName, ntag, typeName})
		} else {
			fieldList = append(fieldList, &structFieldDetails{fieldName, fieldTag, typeName})
		}
	}
	return fieldList
}

func structValueProcess(structName string, form map[string]interface{}) map[string]interface{} {

	var rform = make(map[string]interface{})
	fslc := strStructToFieldsType(structName) //
	for _, fd := range fslc {

		//fmt.Println(i, fd.FiledName, fd.TagName, fd.Type)
		val := fmt.Sprintf("%v", form[fd.TagName])
		if fd.Type == "int" {
			kval, _ := strconv.Atoi(val)
			rform[fd.TagName] = kval

		} else if fd.Type == "int64" {
			kval, _ := strconv.ParseInt(val, 10, 64)
			rform[fd.TagName] = kval

		} else if fd.Type == "float64" {
			kval, _ := strconv.ParseFloat(val, 64)
			rform[fd.TagName] = kval

		} else if fd.Type == "string" {
			rform[fd.TagName] = val

		} else {
			rform[fd.TagName] = form[fd.TagName]
		}
	}
	return rform
}

func vMapToJsonStr(vMap map[string]interface{}) string {

	bs, err := json.Marshal(&vMap)
	if err != nil {
		return ""
	}
	return string(bs)
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

	//text := "ActivityLogTable"
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

func tableToBucket(table string) string {
	return fmt.Sprintf(`%s.%s.%s`, BUCKET_NAME, SCOPE_NAME, table)
}

// func findAailabledocID(table string) (string, int64) {

// 	var docID string
// 	//var scount int64 = int64(maxDoc("cbi.Bucket", table, db))
// 	var scount int64 = tableCount(table)
// 	i := scount + 1
// 	bucket := tableToBucket(table)

// 	for {
// 		count := i
// 		docID = fmt.Sprintf(`%s::%v`, table, count)
// 		//sql := fmt.Sprintf(`SELECT COUNT(*)as cnt FROM %s WHERE META().id="%s";`, bucketName, docID)
// 		sql := fmt.Sprintf(`SELECT META(d).id, true AS docexist FROM %s AS d USE KEYS ["%s"];`, bucket, docID)
// 		//fmt.Println(sql)
// 		pRes := db.Query(sql)
// 		rows := pRes.GetRows()

// 		if len(rows) == 0 {
// 			return docID, count
// 		}
// 		exist := rows[0]["docexist"].(bool)
// 		fmt.Println("ALREADY EXIST", i, exist, sql)
// 		i++
// 	}
// }

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

func upsertQueryBuilder(bucketName, docID, bytesTxt string) (nqlStatement string) {

	qs := `UPSERT INTO %s (KEY, VALUE)
	VALUES ("%s", %s)
	RETURNING *`
	nqlStatement = fmt.Sprintf(qs, bucketName, docID, bytesTxt)
	return
}

func InsertUpdateMap(form map[string]interface{}, db *sql.DB) error {

	//Struct to its fields
	//tableName := form["table"].(string)
	//json.Marshal(&form)
	//tableToBucket(modelName)
	modelName, isOk := form["table"].(string) //collection
	if !isOk {
		return fmt.Errorf("table missing")
	}
	tableName, isOk := form["type"].(string)
	if !isOk {
		tableName = customTableName(modelName)
	}
	docID, isOk := form["id"].(string)
	if !isOk {
		return fmt.Errorf("table missing")
	}

	form2 := structValueProcess(modelName, form) //n1ql
	jsonTxt := vMapToJsonStr(form2)
	//fmt.Println(jsonTxt)
	query := upsertQueryBuilder(tableToBucket(tableName), docID, jsonTxt)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		return err
	}
	lcount, _ := res.RowsAffected()
	fmt.Println("insert:", lcount)
	return nil
}

// CheckCount Get row count using where condition
func CheckCount(table, where string, db *sql.DB) (count int64) {

	sql := fmt.Sprintf("SELECT count(*)as cnt FROM %v WHERE %v;", tableToBucket(table), where)
	rows := db.QueryRow(sql)
	var cmap = make(map[string]interface{}, 0)
	var jsonBytes []uint8
	err = rows.Scan(&jsonBytes)
	if err != nil {
		log.Println("CheckCount:", err.Error())
	}
	json.Unmarshal(jsonBytes, &cmap)
	if len(cmap) > 0 {
		count, _ = strconv.ParseInt(fmt.Sprint(cmap["cnt"]), 10, 64) //float64
	}
	return
}

func addCompany(companyName string) error {

	modelName := "Company"
	var form = make(map[string]interface{})
	id := xid.New().String()
	form["id"] = id
	form["company_name"] = companyName
	form["table"] = modelName
	form["type"] = customTableName(modelName)
	form["serial"] = 1
	form["status"] = 1

	err = InsertUpdateMap(form, db)
	return err
}
