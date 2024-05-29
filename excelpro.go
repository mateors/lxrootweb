package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func dateParser(myDateString string) (string, error) {

	//11-Dec-21
	//01-11-21
	//31/12/2021 23:59:41
	//5/12/2021  3:45:55 PM
	//1/12/21 08:04
	//5-Dec-21
	//20211228074514UTC =>20060102150405MST
	//12/22/2021 => mm/dd/yyyy > 01/02/2006
	//12-22-21
	//2021-11-16 14:00:00.0000000 +06:00
	//12/9/2021  > 	1/11/2021
	//2021-12-27 10:45:22.6245698 +06:00 => 2006-01-02 15:04:05.9999999 Z07:00
	// 	1/11/2021
	// 20-11-2021
	var layout string = "2006-01-02"
	if myDate, err := time.Parse("02-01-2006", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("02/01/2006 15:04:05", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("02/01/2006 3:04:05 PM", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("20060102150405MST", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("01/02/2006", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("2006-01-02 15:04:05.9999999 Z07:00", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("1/02/2006", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("1/2/2006", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("2/01/06 3:04", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("2/01/06 15:04", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("02/01/06 15:04", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("02/01/2006", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("02/Jan/2006", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("02-Jan-06", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("2-Jan-06", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("02-01-06", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("01-02-06", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("1-2-06", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	if myDate, err := time.Parse("02-01-2006", myDateString); err == nil {
		return myDate.Format(layout), nil
	}
	return "", fmt.Errorf("error on date parsing *%s*", myDateString)
}

// bulk excel reader with column formatting
// excelReader("file.xlsx", "Sheet1", nil)
func excelReader(filePath, sheetName string, dateTimeFormat map[string]string) ([]map[string]interface{}, error) {

	sRows := make([]map[string]interface{}, 0)
	keys := make([]string, 0) //first row columns

	xlf, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	if dateTimeFormat != nil {

		if dateCols, isFound := dateTimeFormat["date"]; isFound {
			dates := strings.Split(dateCols, ",")
			//https://xuri.me/excelize/en/style.html#number_format

			styleID, _ := xlf.NewStyle(&excelize.Style{NumFmt: 15}) //14 == `{"number_format":15}`
			for _, column := range dates {
				//xlf.SetCellStyle(sheetName, "K", "K", styleID)
				err = xlf.SetColStyle(sheetName, column, styleID) //column must be in capital
				if err != nil {
					log.Println("ErrDateFormat", err)
				}
			}
		}

		if dateTimeCols, isFound := dateTimeFormat["datetime"]; isFound {
			dates := strings.Split(dateTimeCols, ",")
			styleID, _ := xlf.NewStyle(&excelize.Style{NumFmt: 22}) //`{"number_format":22}`
			for _, column := range dates {
				err = xlf.SetColStyle(sheetName, column, styleID) //column must be in capital
				if err != nil {
					log.Println("ErrDateFormat", err)
				}
			}
		}

	}

	var i int
	erows, err := xlf.Rows(sheetName)
	if err != nil {
		return nil, err
	}

	for erows.Next() {

		cols, err := erows.Columns()
		if err != nil {
			return nil, err
		}
		if len(cols) == 0 {
			continue
		}

		i++
		if i == 1 {
			keys = cols
			continue
		}

		vmap := make(map[string]interface{})
		for j := 0; j < len(keys); j++ {
			key := strings.TrimSpace(keys[j]) //trime
			if len(key) > 0 {
				var val string
				if j < len(cols) {
					val = strings.TrimSpace(cols[j]) //trime
				}
				vmap[key] = val
				//fmt.Println(key, "==>", val)
			}
		}
		if len(vmap) > 0 {
			sRows = append(sRows, vmap)
		}
	}
	return sRows, nil
}

// /"Portwallet (Raw)"
func excelReaderF(filePath, sheetName string, dMap, colMap map[string]string) ([]map[string]interface{}, error) {

	//filePath := "data/inputs/revenue_december_2021/Portwallet_DONE.xlsx"
	var rows = make([]map[string]interface{}, 0)
	// dMap := map[string]string{
	// 	"date": "K", //K,B
	// }
	// colMap := map[string]string{
	// 	"K": "Server Time",
	// }
	erows, err := excelReader(filePath, sheetName, dMap)
	if err != nil {
		return nil, err
	}

	for _, row := range erows {

		var nrow = make(map[string]interface{})
		for _, val := range colMap {
			stime := row[val].(string)
			delete(row, val)
			idate, _ := dateParser(stime)
			//fmt.Printf("%s %s\n", val, idate)
			nrow[val] = idate
		}
		for key, val := range row {
			//fmt.Printf("%d %s ==> %s\n", i, key, val)
			nrow[key] = val
		}
		rows = append(rows, nrow)
	}
	return rows, nil
}

func colID(num int) string {

	var alphabets string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if num <= 0 {
		return ""
	}
	if num < 26 {
		return fmt.Sprintf(`%c`, alphabets[num-1])
	}
	rem := num % 26 //remainder
	fld := num / 26 //floor-division
	if rem == 0 {
		rem = 26
		if fld == 1 {
			return fmt.Sprintf(`%c`, alphabets[rem-1])
		}
		return fmt.Sprintf(`%s%c`, colID(fld-1), alphabets[rem-1])
	}
	return fmt.Sprintf(`%s%c`, colID(fld), alphabets[rem-1])
}

func excelWriter(nameWithoutExt string, nrows, outMaps []map[string]interface{}) error {

	f := excelize.NewFile()
	outputSheetName := "Sheet1"
	sheet1, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}
	for i, row := range outMaps {
		colname := fmt.Sprintf("%v", row["colname"])
		columnID := fmt.Sprintf("%s1", colID(i+1))
		f.SetCellValue(outputSheetName, columnID, colname)
	}

	//fmt.Println("nrows:", len(nrows), outMaps)
	for i, row := range nrows {
		for j, orow := range outMaps {
			columnID := fmt.Sprintf("%s%d", colID(j+1), i+2)
			colname := fmt.Sprintf("%v", orow["colname"])
			f.SetCellValue(outputSheetName, columnID, row[colname])
		}
	}

	f.SetActiveSheet(sheet1)
	if err := f.SaveAs(nameWithoutExt); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// func countryExportAsExcel(outputFile string) error {

// 	sql := "SELECT * FROM country;"
// 	rows, err := msql.GetAllRowsByQuery(sql, database.DB)
// 	if err != nil {
// 		return err
// 	}
// 	outMaps := make([]map[string]interface{}, 0)
// 	if len(rows) > 0 {
// 		row := rows[0]
// 		for key := range row {
// 			cmap := make(map[string]interface{})
// 			cmap["colname"] = key
// 			outMaps = append(outMaps, cmap)
// 		}
// 	}
// 	return excelWriter(outputFile, rows, outMaps) //country.xlsx
// }
