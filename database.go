package main

import (
	"database/sql"
	"fmt"
	"lxrootweb/lxql"
	"net/url"

	"github.com/rs/xid"
)

var (
	COMPANY_TABLE = "Company"
)

// func checkCollectionExist(collectionName string) bool {

// 	qs := `SELECT aid,count,name FROM %s USE KEYS["%s"];`
// 	sql := fmt.Sprintf(qs, tableToBucket("tableRowCount"), collectionName)
// 	pres := db.Query(sql)
// 	rows := pres.GetRows()
// 	return len(rows) != 0
// }

// func checkFileExist(filePath string) error {
// 	var err error
// 	if _, err = os.Stat(filePath); err == nil {
// 		return nil
// 	} else if errors.Is(err, os.ErrNotExist) {
// 		return fmt.Errorf("file %s does not exist!", filePath)
// 	}
// 	return err
// }

func createCollection(collectionName string, db *sql.DB) error {

	_, err = db.Exec("CREATE COLLECTION ?;", tableToBucket(collectionName))
	//sql := fmt.Sprintf("CREATE COLLECTION %s;", tableToBucket(collectionName))
	//pres := db.Query(sql)
	// if pres.Status != "success" {
	// 	var emsg string
	// 	for _, e := range pres.Errors {
	// 		emsg += fmt.Sprintf("%s-%s,", e.Message, e.Message)
	// 	}
	// 	log.Println("ERRsql:", sql)
	// 	return fmt.Errorf("ERROR %s", strings.TrimRight(emsg, ","))
	// }
	//time.Sleep(time.Millisecond * 200)
	return err
}

// func createSecondaryIndex(collectionName string) error {

// 	sql := fmt.Sprintf("CREATE INDEX `status%s` ON %s(status);", collectionName, tableToBucket(collectionName))
// 	//fmt.Println(sql)
// 	pres := db.Query(sql)
// 	if pres.Status != "success" {
// 		log.Println("ERRsql:", sql)
// 		return fmt.Errorf("ERROR %s", singleMsg(pres))
// 	}
// 	time.Sleep(time.Millisecond * 200)
// 	return nil
// }

// func createPrimaryIndex(collectionName string) error {
// 	sql := fmt.Sprintf("CREATE PRIMARY INDEX ON %s;", tableToBucket(collectionName))
// 	//fmt.Println(sql)
// 	pres := db.Query(sql)
// 	if pres.Status != "success" {
// 		log.Println("ERRsql:", sql)
// 		return fmt.Errorf("ERROR %s", singleMsg(pres))
// 	}
// 	time.Sleep(time.Millisecond * 200)
// 	return nil
// }

// func collectionInstance(collectionName string) error {

// 	var vmap = make(map[string]interface{})
// 	vmap["aid"] = collectionName
// 	vmap["count"] = 0
// 	vmap["name"] = collectionName
// 	vmap["type"] = "tableRowCount"
// 	insertStatm := `INSERT INTO %s (KEY,VALUE) VALUES("%s",%s);`
// 	bs, err := json.Marshal(&vmap)
// 	if err != nil {
// 		return err
// 	}
// 	sql := fmt.Sprintf(insertStatm, tableToBucket("tableRowCount"), collectionName, string(bs))
// 	pres := db.Query(sql)
// 	if pres.Status != "success" {
// 		return fmt.Errorf("ERROR %s", singleMsg(pres))
// 	}
// 	time.Sleep(time.Millisecond * 200)
// 	return nil
// }

func basicForm() {

	var form = make(map[string]interface{})
	form["id"] = "id::1"
	form["company_name"] = "Mostain"
	form["age"] = "40"
	form["lang"] = []string{"golang", "rust"}
	form["table"] = "Company"

	err = lxql.InsertUpdateMap(form, db)
	fmt.Println(err)
}

func nextSerial(tableName string) int {
	return lxql.CheckCount(tableName, fmt.Sprintf("type='%s'", tableName), db) + 1
}

func addCompany(companyName string) error {

	table := customTableName(COMPANY_TABLE)
	var form = make(map[string]interface{})
	id := xid.New().String()
	form["id"] = id
	form["company_name"] = companyName
	form["table"] = COMPANY_TABLE //model
	form["type"] = table
	form["serial"] = nextSerial(table)
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, db)
	return err
}

func modelAction(modelName string, form url.Values) error {

	var mForm = make(map[string]interface{})
	table := customTableName(modelName) //database table
	mForm["table"] = modelName
	mForm["id"] = xid.New().String()
	mForm["type"] = table
	mForm["serial"] = nextSerial(table)

	for key := range form {
		mForm[key] = form.Get(key)
	}
	err = lxql.InsertUpdateMap(mForm, db)
	return err
}
