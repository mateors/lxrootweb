package main

import (
	"encoding/json"
	"fmt"
	"log"
	"lxrootweb/database"
	"lxrootweb/lxql"
	"net/http"
)

const (
	PAYMENT_INENT_CREATED     = "payment_intent.created"    //1 ok
	INVOICE_CREATED           = "invoice.created"           //2 ok
	INVOICE_FINALIZED         = "invoice.finalized"         //3 ok
	CHARGE_SUCCEEDED          = "charge.succeeded"          //4 ok
	PAYMENT_INTENT_SUCCEEDED  = "payment_intent.succeeded"  //5 ok
	INVOICE_UPDATED           = "invoice.updated"           //6 ok
	INVOICE_PAID              = "invoice.paid"              //7 ok
	INVOICE_PAYMENT_SUCCEEDED = "invoice.payment_succeeded" //8 ok
)

func paymentHook(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		//var eMap = make(map[string]interface{})
		//err := json.NewDecoder(r.Body).Decode(&eMap)
		var evt Event //StripeEvent
		err := json.NewDecoder(r.Body).Decode(&evt)
		if err != nil {
			log.Println(err)
			return
		}
		// evt.Type, isOk := eMap["type"].(string)
		// if !isOk {
		// 	log.Println("unable to parse evt.Type")
		// 	return
		// }
		// created, _ := eMap["created"].(string)
		// createdInt64, _ := strconv.ParseInt(created, 10, 64)
		// fmt.Println("created:", unixToDateTime(createdInt64))

		if evt.Type == PAYMENT_INENT_CREATED { //1
			fmt.Println("1-->", PAYMENT_INENT_CREATED)

		} else if evt.Type == INVOICE_CREATED { //2
			fmt.Println("2-->", INVOICE_CREATED)

		} else if evt.Type == INVOICE_FINALIZED { //3
			fmt.Println("3-->", INVOICE_FINALIZED)

		} else if evt.Type == CHARGE_SUCCEEDED { //4
			fmt.Println("4-->", CHARGE_SUCCEEDED)

		} else if evt.Type == PAYMENT_INTENT_SUCCEEDED { //5
			fmt.Println("5-->", PAYMENT_INTENT_SUCCEEDED)

		} else if evt.Type == INVOICE_UPDATED { //6
			fmt.Println("6-->", INVOICE_UPDATED)

		} else if evt.Type == INVOICE_PAID { //7
			fmt.Println("7-->", INVOICE_PAID)

		} else if evt.Type == INVOICE_PAYMENT_SUCCEEDED { //8
			fmt.Println("8-->", INVOICE_PAYMENT_SUCCEEDED)

		} else {
			fmt.Println("--> Unknown event.Type")
		}

		//_, err = addEvent() //lxql.InsertUpdateMap(eMap, database.DB)
		err = lxql.InsertUpdateObject("event", evt.ID, &evt, database.DB)
		logError("eventInsertERR", err)
		w.WriteHeader(http.StatusOK)

	} else {
		fmt.Fprintln(w, r.RemoteAddr)
	}
}
