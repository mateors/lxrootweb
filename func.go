package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
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

func couchbaseConnTest() {

	//couchbase
	//"user:password@/dbname"
	n1ql, err := sql.Open("n1ql", "http://172.93.55.179:8093")
	if err != nil {
		fmt.Println("***", err)
		//log.Fatal(err)
		return
	}

	//ac := []byte(`[{"user": "admin:Administrator", "pass": "asdasd"}]`)
	//go_n1ql.SetQueryParams("creds", string(ac))
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

	rows, err := n1ql.Query("select element (select * from lxroot)")

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var contacts string
		if err := rows.Scan(&contacts); err != nil {
			log.Fatal(err)
		}
		log.Printf("Row returned %s : \n", contacts)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows.Close()

}
