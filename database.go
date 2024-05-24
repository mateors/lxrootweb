package main

import (
	"database/sql"
	"fmt"
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

	err = InsertUpdateMap(form, db)
	fmt.Println(err)
}
