package lxql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var (
	BUCKET, SCOPE string
)
var typeRegistry = make(map[string]reflect.Type)

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

func RegisterModel(emptyStruct interface{}) {
	typeRegistry[reflect.TypeOf(emptyStruct).Name()] = reflect.TypeOf(emptyStruct)
}

func makeInstance(name string) interface{} {
	rval, isFound := typeRegistry[name]
	if !isFound {
		return nil
	}
	return reflect.New(rval).Elem().Interface()
}

func structNameToFields(structName string) (fieldList []*structFieldDetails) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("was panic, recovered value", r)
			fieldList = nil
		}
	}()

	sInstance := makeInstance(structName) //
	iVal := reflect.ValueOf(sInstance)
	iTypeOf := iVal.Type()

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
			fieldList = append(fieldList, &structFieldDetails{fieldName, ntag, typeName, omitFound})
		} else {
			fieldList = append(fieldList, &structFieldDetails{fieldName, fieldTag, typeName, omitFound})
		}
	}
	return fieldList
}

func strStructToFields(structName string) []string {

	fieldList := []string{}
	sInstance := makeInstance(structName) //
	if sInstance == nil {
		return nil
	}
	iVal := reflect.ValueOf(sInstance)
	iTypeOf := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {

		fieldName := iTypeOf.Field(i).Name
		fieldTag := iTypeOf.Field(i).Tag.Get("json")
		if fieldTag == "" {
			fieldList = append(fieldList, fieldName)
		} else if strings.Contains(fieldTag, ",") {
			slc := strings.Split(fieldTag, ",")
			fieldList = append(fieldList, slc[0])
		} else {
			fieldList = append(fieldList, fieldTag)
		}
	}
	return fieldList
}

func structValueProcess(structName string, form map[string]interface{}) map[string]interface{} {

	var rform = make(map[string]interface{})
	fslc := structNameToFields(structName) //

	for _, fd := range fslc {

		val, isParsed := form[fd.TagName]
		if !isParsed {
			val = ""
		}
		if fd.Type == "int" {
			kval, _ := strconv.Atoi(fmt.Sprint(val))
			rform[fd.TagName] = kval

		} else if fd.Type == "int64" {
			kval, _ := strconv.ParseInt(fmt.Sprint(val), 10, 64)
			rform[fd.TagName] = kval

		} else if fd.Type == "float64" {
			kval, _ := strconv.ParseFloat(fmt.Sprint(val), 64)
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

	var list []string
	ci := upperCount(text)
	splitIndex := 0
	for i, c := range ci {
		if i > 0 {
			list = append(list, text[splitIndex:c.Index])
			splitIndex = c.Index
		}
	}
	if splitIndex < len(text) {
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
	return fmt.Sprintf(`%s.%s.%s`, BUCKET, SCOPE, table)
}

func upsertQueryBuilder(bucketName, docID, bytesTxt string) (nqlStatement string) {

	qs := `UPSERT INTO %s (KEY, VALUE)
	VALUES ("%s", %s)
	RETURNING *`
	nqlStatement = fmt.Sprintf(qs, bucketName, docID, bytesTxt)
	return
}

func RawSQL(sql string, db *sql.DB) error {

	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec()
	if err != nil {
		return err
	}

	n, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}
	return err
}

func InsertUpdateMap(form map[string]interface{}, db *sql.DB) error {

	modelName, isFound := form["table"].(string) //collection
	if !isFound {
		return fmt.Errorf("collection name missing")
	}
	docID, isFound := form["id"].(string)
	if !isFound {
		return fmt.Errorf("id missing")
	}
	tableName, isFound := form["type"].(string)
	if !isFound {
		tableName = customTableName(modelName)
		form["type"] = tableName
	}

	form2 := structValueProcess(modelName, form) //n1ql
	jsonTxt := vMapToJsonStr(form2)
	query := upsertQueryBuilder(tableToBucket(tableName), docID, jsonTxt)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func scanMap(jsonBytes []uint8) map[string]interface{} {

	var cmap = make(map[string]interface{}, 0)
	err := json.Unmarshal(jsonBytes, &cmap)
	if err != nil {
		return nil
	}
	return cmap
}

// CheckCount Get row count using where condition
func CheckCount(table, where string, db *sql.DB) (count int) {

	sql := fmt.Sprintf("SELECT count(*)as cnt FROM %s WHERE %s;", tableToBucket(table), where)
	rows := db.QueryRow(sql)
	var jsonBytes []uint8
	err := rows.Scan(&jsonBytes)
	if err != nil {
		log.Println("CheckCount:", err.Error())
		return 0
	}
	cmap := scanMap(jsonBytes)
	if len(cmap) > 0 {
		count, _ = strconv.Atoi(fmt.Sprint(cmap["cnt"]))
	}
	return
}

// ReadTable2Columns Get table all columns as a slice of string
func ReadTable2Columns(table string, db *sql.DB) (cols []string, err error) {
	cols = strStructToFields(table)
	return
}

func bytesToStr(slc []uint8) string {

	var str string
	for _, c := range slc {
		if c != 34 { //remove "
			str += fmt.Sprintf("%c", c)
		}
	}
	return str
}

func colsToRowMap(cols []string, orow map[string]interface{}) map[string]interface{} {

	if len(cols) > 1 {
		for key := range orow {
			slc := orow[key].([]uint8)
			orow[key] = bytesToStr(slc)
		}

	} else if len(cols) == 1 {

		colname := cols[0]
		if colname == "*" {
			var row = make(map[string]map[string]interface{})
			json.Unmarshal(orow[colname].([]uint8), &row)
			for key := range row {
				orow = row[key]
			}

		} else {
			var row = make(map[string]interface{})
			json.Unmarshal(orow[colname].([]uint8), &row)
			orow = row
		}
	}
	return orow
}

// GetAllRowsByQuery Get all table rows using raw sql query
func GetRows(sql string, db *sql.DB) ([]map[string]interface{}, error) {

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	row := make([][]byte, len(cols))
	rowPtr := make([]any, len(cols))

	for i := range row {
		rowPtr[i] = &row[i]
	}
	defer rows.Close()
	tableData := make([]map[string]interface{}, 0)

	for rows.Next() {

		if err := rows.Scan(rowPtr...); err != nil {
			return nil, err
		}
		var orow = make(map[string]interface{})

		for i, rowp := range rowPtr {

			switch val := rowp.(type) {

			case *[]uint8:
				orow[cols[i]] = *val

			default:
				fmt.Println("GetRows() Type is unknown!")
			}
		}
		tableData = append(tableData, colsToRowMap(cols, orow))
	}
	return tableData, nil
}

// FieldByValue Get one field_value using where clause
func FieldByValue(table, fieldName, where string, db *sql.DB) string {

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s;", fieldName, table, where)
	rows := db.QueryRow(sql)
	var vfield string
	var jsonBytes []uint8
	err := rows.Scan(&jsonBytes)
	if err != nil {
		log.Println("CheckCount:", err.Error())
		return ""
	}
	cmap := scanMap(jsonBytes)
	if len(cmap) > 0 {
		vfield = fmt.Sprint(cmap[fieldName])
	}
	return vfield
}
