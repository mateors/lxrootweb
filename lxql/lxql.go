package lxql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
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

	_, err := db.Exec(sql)
	// stmt, err := db.Prepare(sql)
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()

	// r, err := stmt.Exec()
	// if err != nil {
	// 	return err
	// }

	// n, err := r.RowsAffected()
	// if err != nil {
	// 	return err
	// }

	// if n > 0 {
	// 	return nil
	// }
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

func structFieldValMap(anyStructToPointer interface{}) (map[string]interface{}, error) {

	empv := reflect.ValueOf(anyStructToPointer) //must be a pointer
	empt := reflect.TypeOf(anyStructToPointer)  //
	if empv.Kind() != reflect.Pointer {
		return nil, errors.New("anyStructToPointer must be a pointer")
	}
	var row = make(map[string]interface{})
	for i := 0; i < empv.Elem().NumField(); i++ {

		var key string
		kfield := empt.Elem().Field(i)
		key = kfield.Name
		tag := kfield.Tag.Get("json")
		if tag != "" {
			key = tag
		}
		row[key] = empv.Elem().Field(i).Interface()
	}
	return row, nil
}

// objPtr is StructPropertyWithValue
func InsertUpdateObject(tableName, docID string, objPtr interface{}, db *sql.DB) error {

	row, err := structFieldValMap(objPtr)
	if err != nil {
		return err
	}
	jsonTxt := vMapToJsonStr(row)
	query := upsertQueryBuilder(tableToBucket(tableName), docID, jsonTxt)
	_, err = db.Exec(query)
	return err
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

func cleanText(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func getColsFromSql(sql string) (cols []string) {

	csql := strings.ToLower(cleanText(sql))
	si := strings.Index(csql, "select")
	fi := strings.Index(csql, "from")
	commaSeparatedFields := csql[si+7 : fi]
	slc := strings.Split(commaSeparatedFields, ",")
	for _, col := range slc {
		colname := sqlColumnPattern(strings.TrimSpace(col))
		cols = append(cols, colname)
	}
	return
}

func regExFindMatch(pattern, data string) (match []string) {

	var myExp = regexp.MustCompile(pattern)
	match = myExp.FindStringSubmatch(data)
	return
}

func sqlColumnPattern(colname string) string {

	var fcolname string
	patterns := []string{`(.*)\s+as\s+([a-z_]*)`, `(.*)\.\s+([a-z_]*)`, `(.*)\.([a-z_]*)`}

	for _, pattern := range patterns {
		slc := regExFindMatch(pattern, colname)
		if len(slc) == 3 {
			fcolname = slc[2]
			break
		}
	}
	if len(fcolname) == 0 {
		fcolname = colname
	}
	return fcolname
}

// func SingleRow(sql string, db *sql.DB) (map[string]interface{}, error) {

// 	var orow = make(map[string]interface{})

// 	cols := getColsFromSql(sql) //GetColumnNamesFromQuery(sql)
// 	count := len(cols)
// 	values := make([]interface{}, count)
// 	valuePtrs := make([]interface{}, count)

// 	fmt.Println(cols)

// 	var isStarFound bool
// 	var colCount int

// 	for i := range cols {
// 		valuePtrs[i] = &values[i]
// 	}

// 	srow := db.QueryRow(sql)
// 	srow.Scan(valuePtrs...)

// 	//var orow = make(map[string]interface{})
// 	for i, col := range cols {
// 		colCount++
// 		if col == "*" {
// 			isStarFound = true
// 		}
// 		//val := values[i]
// 		orow[col] = values[i]
// 		fmt.Printf("%d %s %v %T\n", i, col, bytesToStr(values[i].([]byte)), values[i])
// 	}

// 	if isStarFound {

// 		for _, col := range cols {
// 			//fmt.Printf("%v = %v %T\n", col, bytesToStr(row[col].([]uint8)), row[col])
// 			var vmap = make(map[string]interface{})
// 			json.Unmarshal(orow[col].([]uint8), &vmap)
// 			for key := range vmap {
// 				vrow, isOk := vmap[key].(map[string]interface{})
// 				if isOk {
// 					//nrows = append(nrows, vrow)
// 					orow = vrow
// 				} else {
// 					fmt.Printf("%v %T\n", vmap[key], vmap[key])
// 				}
// 			}
// 		}

// 	} else if colCount == 1 {

// 		for _, val := range orow {
// 			json.Unmarshal(val.([]uint8), &orow)
// 		}

// 	} else if colCount > 1 {

// 		for _, col := range cols {

// 			val, isOk := orow[col].([]uint8)
// 			if isOk {
// 				orow[col] = bytesToStr(val)
// 			}
// 			//fmt.Println(">", col, orow[col])
// 		}

// 	}
// 	//orow = colsToRowMap(cols, orow)
// 	return orow, nil
// }

// func SingleRow(sql string, db *sql.DB) (map[string]interface{}, error) {

// 	var orow = make(map[string]interface{})

// 	cols := getColsFromSql(sql) //GetColumnNamesFromQuery(sql)
// 	row := make([][]byte, len(cols))
// 	rowPtr := make([]any, len(cols))

// 	for i := range row {
// 		rowPtr[i] = &row[i]
// 	}

// 	srow := db.QueryRow(sql)
// 	err := srow.Scan(rowPtr...)
// 	if err != nil {
// 		log.Println(">>", err)
// 		return nil, err
// 	}

// 	for i, rowp := range rowPtr {

// 		switch val := rowp.(type) {

// 		case *[]uint8:
// 			orow[cols[i]] = *val

// 		default:
// 			log.Println("SingleRow() Type is unknown!")
// 		}
// 	}
// 	orow = colsToRowMap(cols, orow)
// 	return orow, nil
// }

func GetRows(sql string, db *sql.DB) ([]map[string]interface{}, error) {

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	var isStarFound bool
	var colCount int

	var orows = make([]map[string]interface{}, 0)

	for rows.Next() {

		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		var orow = make(map[string]interface{})
		for i, col := range columns {
			colCount++
			if col == "*" {
				isStarFound = true
			}
			val := values[i]
			orow[col] = val
		}
		orows = append(orows, orow)
	} //

	//process
	var nrows = make([]map[string]interface{}, 0)

	if isStarFound {

		for _, row := range orows {

			for _, col := range columns {
				var vmap = make(map[string]interface{})
				json.Unmarshal(row[col].([]uint8), &vmap)
				for key := range vmap {
					vrow, isOk := vmap[key].(map[string]interface{})
					if isOk {
						nrows = append(nrows, vrow)
					} else {
						fmt.Printf("%v %T\n", vmap[key], vmap[key])
					}
				}
			}
		}

	} else if colCount == 1 {

		for _, row := range orows {
			for _, val := range row {
				json.Unmarshal(val.([]uint8), &row)
			}
			nrows = append(nrows, row)
		}

	} else if colCount > 1 {

		for _, row := range orows {
			var srow = make(map[string]interface{})
			for key, val := range row {

				vbs, isOk := val.([]uint8)
				if isOk {
					srow[key] = bytesToStr(vbs) //?
				} else {
					srow[key] = val
				}
				//fmt.Printf("%v => %v | %T\n", key, val, val)
			}
			nrows = append(nrows, srow)
		}

	} else {
		//fmt.Println("else:>", orows, len(orows))
	}
	return nrows, nil
}

// GetAllRowsByQuery Get all table rows using raw sql query
// func GetRows(sql string, db *sql.DB) ([]map[string]interface{}, error) {

// 	rows, err := db.Query(sql)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cols, err := rows.Columns()
// 	if err != nil {
// 		return nil, err
// 	}
// 	row := make([][]byte, len(cols))
// 	rowPtr := make([]any, len(cols))

// 	for i := range row {
// 		rowPtr[i] = &row[i]
// 	}
// 	defer rows.Close()
// 	tableData := make([]map[string]interface{}, 0)

// 	for rows.Next() {

// 		if err := rows.Scan(rowPtr...); err != nil {
// 			return nil, err
// 		}
// 		var orow = make(map[string]interface{})

// 		for i, rowp := range rowPtr {

// 			switch val := rowp.(type) {

// 			case *[]uint8:
// 				orow[cols[i]] = *val

// 			default:
// 				fmt.Println("GetRows() Type is unknown!")
// 			}
// 		}
// 		tableData = append(tableData, colsToRowMap(cols, orow))
// 	}
// 	return tableData, nil
// }

// FieldByValue Get one field_value using where clause
func FieldByValue(table, fieldName, where string, db *sql.DB) string {

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s;", fieldName, tableToBucket(table), where)
	rows := db.QueryRow(sql)
	var vfield string
	var jsonBytes []uint8
	err := rows.Scan(&jsonBytes)
	if err != nil {
		//log.Println("FieldByValue:", err.Error())
		return ""
	}
	cmap := scanMap(jsonBytes)
	if len(cmap) > 0 {
		vfield = fmt.Sprint(cmap[fieldName])
	}
	return vfield
}
