package main

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/couchbase/go_n1ql"
	//_ "github.com/couchbase/go_n1ql"
)

func featureBySlug(slugName string, rows []map[string]interface{}) map[string]interface{} {

	for _, row := range rows {
		slug := row["slug"].(string)
		if slug == slugName {
			return row
		}
	}
	return nil
}

func findFeaturesById(featureId int, fRows []map[string]interface{}) []map[string]interface{} {

	var rows = make([]map[string]interface{}, 0)
	for _, row := range fRows {
		fid := row["feature_id"].(int)
		if fid == featureId {
			rows = append(rows, row)
		}
	}
	return rows
}

// Round float to 2 decimal places
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func calculateAppPrice(numApps int) float64 {

	costPerAppRanges := []struct {
		min  int
		max  int
		cost float64
	}{
		{1, 5000, 0.15},
		{5001, 25000, 0.10},
		{25001, 100000, 0.075},
	}

	var totalCost float64
	remainingApps := numApps

	for _, rangeItem := range costPerAppRanges {

		appsInRange := int(math.Min(float64(remainingApps), float64(rangeItem.max-rangeItem.min+1)))
		totalCost += float64(appsInRange) * rangeItem.cost
		remainingApps -= appsInRange

		if remainingApps <= 0 {
			break
		}
	}

	// Ensure minimum billing threshold
	//fmt.Sprintf("%.2f", number)
	totalCost = math.Max(totalCost, 15)
	return totalCost
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

func couchbaseConnTest() {

	//couchbase
	//"user:password@/dbname"
	n1ql, err := sql.Open("n1ql", "http://172.93.55.179:8093")
	if err != nil {
		fmt.Println("***", err)
		//log.Fatal(err)
		return
	}

	ac := []byte(`[{"user": "admin:Administrator", "pass": "Mostain321$"}]`)
	go_n1ql.SetQueryParams("creds", string(ac))

	//go_n1ql.SetQueryParams("timeout", "10s")
	//go_n1ql.SetQueryParams("scan_consistency", "request_plus")

	// err = n1ql.Ping()
	// if err != nil {
	// 	fmt.Println("###")
	// 	log.Fatal(err)
	// }

	// fmt.Println("ping success...")

	// Set query parameters
	//ac := []byte(`[{"user": "admin:Administrator", "pass": "asdasd"}]`)
	//go_n1ql.SetQueryParams("creds", string(ac))
	//go_n1ql.SetQueryParams("timeout", "10s")
	// go_n1ql.SetQueryParams("scan_consistency", "request_plus")

	rows, err := n1ql.Query("select id,name,age from lxroot;")
	if err != nil {
		return
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return
	}

	rc := newMapStringScan(columnNames)
	tableData := make([]map[string]interface{}, 0)

	for rows.Next() {

		//var id, name, age string
		//var age float64

		// if err := rows.Scan(&id, &name, &age); err != nil {
		// 	log.Fatal(err)
		// }
		//log.Printf("Row returned -> %s,%s,%s : \n", id, name, age)

		err = rc.Update(rows)
		if err != nil {
			break
		}
		cv := rc.Get()
		dd := make(map[string]interface{})
		for _, col := range columnNames {
			dd[col] = cv[col]
		}
		tableData = append(tableData, dd)
	}
	fmt.Printf("%v %T\n", tableData, tableData)

	// if err := rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }

}
