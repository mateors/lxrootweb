package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"strings"
	"time"
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

// func couchbaseConnTest() {

// 	//couchbase
// 	//"user:password@/dbname"
// 	n1ql, err := sql.Open("n1ql", "http://172.93.55.179:8093")
// 	if err != nil {
// 		fmt.Println("***", err)
// 		//log.Fatal(err)
// 		return
// 	}

// 	ac := []byte(`[{"user": "admin:Administrator", "pass": "Mostain321$"}]`)
// 	go_n1ql.SetQueryParams("creds", string(ac))

// 	//go_n1ql.SetQueryParams("timeout", "10s")
// 	//go_n1ql.SetQueryParams("scan_consistency", "request_plus")

// 	// err = n1ql.Ping()
// 	// if err != nil {
// 	// 	fmt.Println("###")
// 	// 	log.Fatal(err)
// 	// }

// 	// fmt.Println("ping success...")

// 	// Set query parameters
// 	//ac := []byte(`[{"user": "admin:Administrator", "pass": "asdasd"}]`)
// 	//go_n1ql.SetQueryParams("creds", string(ac))
// 	//go_n1ql.SetQueryParams("timeout", "10s")
// 	// go_n1ql.SetQueryParams("scan_consistency", "request_plus")

// 	rows, err := n1ql.Query("select id,name,age from lxroot;")
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// 	columnNames, err := rows.Columns()
// 	if err != nil {
// 		return
// 	}

// 	rc := newMapStringScan(columnNames)
// 	tableData := make([]map[string]interface{}, 0)

// 	for rows.Next() {

// 		//var id, name, age string
// 		//var age float64

// 		// if err := rows.Scan(&id, &name, &age); err != nil {
// 		// 	log.Fatal(err)
// 		// }
// 		//log.Printf("Row returned -> %s,%s,%s : \n", id, name, age)

// 		err = rc.Update(rows)
// 		if err != nil {
// 			break
// 		}
// 		cv := rc.Get()
// 		dd := make(map[string]interface{})
// 		for _, col := range columnNames {
// 			dd[col] = cv[col]
// 		}
// 		tableData = append(tableData, dd)
// 	}
// 	fmt.Printf("%v %T\n", tableData, tableData)

// 	// if err := rows.Err(); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// stt, err := n1ql.Prepare(`INSERT INTO lxroot (KEY, VALUE) VALUES ("doc2", {"name":"sanzida","age":32,"id":"doc2"}) RETURNING *`)
// 	// stt.Exec()
// 	// fmt.Println("insert:", err)

// 	// stt, err := n1ql.Prepare(`INSERT INTO lxroot (KEY, VALUE) VALUES ("doc3", {"name":"Wania","age":3,"id":"doc3"})`)
// 	// stt.Exec()
// 	// fmt.Println("insert:", err)

// 	// stt, err := n1ql.Prepare(`UPSERT INTO lxroot (KEY, VALUE) VALUES ("doc2", {"name":"Sanzida","age":33,"id":"doc2"})`)
// 	// stt.Exec()
// 	// fmt.Println("insert:", err)

// 	//SELECT id,name,age FROM lxroot USE KEYS "id::1"
// 	//res, err := n1ql.Exec("DELETE FROM lxroot USE KEYS ?", "id::1")
// 	//fmt.Println("delete:", res, err)

// }

func ParseDSN(dsn string) (err error) {

	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	// Find the last '/' (since the password or the net addr might contain a '/')
	var user, passwd, addr, net, dbname string
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			foundSlash = true
			var j, k int

			// left part is empty if i <= 0
			if i > 0 {
				// [username[:password]@][protocol[(address)]]
				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						// username[:password]
						// Find the first ':' in dsn[:j]
						for k = 0; k < j; k++ {
							if dsn[k] == ':' {
								passwd = dsn[k+1 : j]
								break
							}
						}
						user = dsn[:k]

						break
					}
				}

				// [protocol[(address)]]
				// Find the first '(' in dsn[j+1:i]
				for k = j + 1; k < i; k++ {
					if dsn[k] == '(' {
						// dsn[i-1] must be == ')' if an address is specified
						if dsn[i-1] != ')' {
							if strings.ContainsRune(dsn[k+1:i], ')') {
								return errors.New("invalid DSN: did you forget to escape a param value")
							}
							return errors.New("invalid DSN: network address not terminated (missing closing brace)")
						}
						addr = dsn[k+1 : i-1]
						break
					}
				}
				net = dsn[j+1 : k]
			}

			// dbname[?param1=value1&...&paramN=valueN]
			// Find the first '?' in dsn[i+1:]
			for j = i + 1; j < len(dsn); j++ {
				if dsn[j] == '?' {
					//if err = parseDSNParams(cfg, dsn[j+1:]); err != nil {
					//return
					//}
					break
				}
			}

			dbname = dsn[i+1 : j]
			// if cfg.DBName, err = url.PathUnescape(dbname); err != nil {
			// 	return fmt.Errorf("invalid dbname %q: %w", dbname, err)
			// }
			break
		}
	}

	if !foundSlash && len(dsn) > 0 {
		return errors.New("invalid DSN: missing the slash separating the database name")
	}

	// if err = cfg.normalize(); err != nil {
	// 	return  err
	// }
	fmt.Println(user, passwd, addr, net, dbname)
	// Qparams["user"] = user
	// Qparams["passwd"] = passwd
	// Qparams["address"] = addr
	// Qparams["protocol"] = net
	// Qparams["dbname"] = dbname
	return
}

func lxqlCon() {

	//go_n1ql.OpenN1QLConnection()
	//lxql.OpenConn()
	//lxql.ParseDSN("username:password@tcp(localhost:8309)/lxrootdb")
	//n1ql2, err := sql.Open("n1ql", "lxrtestusr:Test54321$@(172.93.55.179:8309)/lxrootdb")
	n1ql2, err := sql.Open("n1ql", "http://lxrtestusr:Test54321$@172.93.55.179:8093")
	//n1ql2, err := sql.Open("n1ql", "http://172.93.55.179:8093")
	fmt.Println(err)

	//ac := []byte(`[{"user": "admin:lxrtestusr", "pass": "Test54321$"}]`)
	//lxql.SetQueryParams("creds", string(ac))

	err = n1ql2.Ping()
	fmt.Println("ping..", err)

	/*
		stt, err := n1ql2.Prepare("SELECT id,name,age FROM lxroot WHERE name=? AND age=?")
		if err != nil {
			log.Println(err)
			return
		}

		rows, err := stt.Query("Sanzida", 33)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {

			var id, name, age string
			if err := rows.Scan(&id, &name, &age); err != nil {
				log.Fatal(err)
			}
			log.Printf("-> %s,%s,%s : \n", id, name, age)
		}
	*/

	/*
		stmt, err := n1ql2.Prepare("Upsert INTO lxroot values (?,?)")
		if err != nil {
			log.Fatal(err)
		}

		// Map Values need to be marshaled
		value, _ := json.Marshal(map[string]interface{}{"name": "irish", "type": "contact"})
		result, err := stmt.Exec("irish4", value)
		if err != nil {
			log.Fatal(err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Rows affected %d", rowsAffected)

		for i := 0; i < 20; i++ {

			key := fmt.Sprintf("irish%d", i)
			result, err = stmt.Exec(key, value)
			if err != nil {
				log.Fatal(err)
			}

			ra, err := result.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			rowsAffected += ra
		}
		log.Printf("Total Rows Affected %d", rowsAffected)
		stmt.Close()
	*/

	// lid, err := res.LastInsertId()
	// fmt.Println("lastInsert:", lid, err)

	//res, err := n1ql2.Exec("DELETE FROM lxroot USE KEYS ?", "doc4")
	//fmt.Println("delete:", res, err)

	//os.Exit(1)
}

func IsIPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() != nil
}

func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() == nil
}

func getCountryRegionFromIp(ipAddress string) (map[string]string, error) {

	var row = make(map[string]string)
	rurl := fmt.Sprintf("http://ip-api.com/json/%s?fields=country,regionName", ipAddress)
	resp, err := http.Get(rurl)
	if err != nil {
		return nil, err
	}

	remaining := resp.Header.Get("X-Rl")
	ttl := resp.Header.Get("X-Ttl")
	err = json.NewDecoder(resp.Body).Decode(&row)
	if err != nil {
		return nil, err
	}
	row["remaining"] = remaining //
	row["ttl"] = ttl             // check if it is 0
	row["status"] = resp.Status  // check if 429
	//fmt.Println(row)
	return row, nil
}

func getIpToCountry(ipv4Address string) (string, error) {

	if !IsIPv4(ipv4Address) {
		return "", errors.New("ipv4 allowed only")
	}
	url := fmt.Sprintf("http://ip2c.org/?ip=%s", ipv4Address)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	text := string(bs)
	//slc := strings.Split(text, ";")
	//fmt.Println(slc, len(slc))
	return text, nil
}

func getLocation(ipAddress string) (string, error) {

	var country, region string
	row, err := getCountryRegionFromIp(ipAddress)
	if err != nil {
		return "", err
	}
	country = row["country"]
	region = row["regionName"]
	status := row["status"]
	region = strings.TrimSpace(strings.Replace(region, "Division", "", -1))
	//fmt.Println(row["ttl"], row["remaining"])

	if strings.Contains(status, "429") {
		if IsIPv4(ipAddress) {
			txt, err := getIpToCountry(ipAddress)
			logError("", err)
			slc := strings.Split(txt, ";")
			if len(slc) == 4 {
				country = slc[3]
			}
		} else {
			region = "IP"
			country = ipAddress
			log.Println(status, ipAddress, "ipv6 address has no provider yet")
		}
	}
	//time.Sleep(time.Second * 1)
	return fmt.Sprintf("%s,%s", region, country), nil
}

func getLocationWithinSec(ipAddress string) string {

	var location string
	// type output struct {
	// 	out string
	// 	err error
	// }
	ch := make(chan string)
	go func() {
		res, err := getLocation(ipAddress)
		logError("getLocation", err)
		ch <- res
	}()

	select {

	case <-time.After(1 * time.Second):
		location = "" //"timed out"
		log.Println("getLocation() time out")

	case x := <-ch:
		location = x
	}
	return location
}

func performTask(ctx context.Context, output chan<- string) {

	time.Sleep(time.Second * 3)
	output <- "mostain"
}

func ctxReq() {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan string)
	go performTask(ctx, ch)

	select {
	case <-ctx.Done():
		fmt.Println("Task timed out")

	case x := <-ch:
		fmt.Println(x)
	}
}
