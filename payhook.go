package main

import (
	"encoding/json"
	"fmt"
	"log"
	"lxrootweb/database"
	"lxrootweb/lxql"
	"net/http"
	"time"

	"github.com/mateors/mtool"
	uuid "github.com/satori/go.uuid"
)

const (
	CUSTOMER_CREATED              = "customer.created"              //1
	CUSTOMER_SUBSCRIPTION_CREATED = "customer.subscription.created" //2
	INVOICE_CREATED               = "invoice.created"               //2 ok
	CHARGE_SUCCEEDED              = "charge.succeeded"              //4 ok
	INVOICE_UPDATED               = "invoice.updated"               //6 ok
	INVOICE_PAID                  = "invoice.paid"                  //7 ok
	INVOICE_PAYMENT_SUCCEEDED     = "invoice.payment_succeeded"     //8 ok
	CHECKOUT_SESSION_COMPLETED    = "checkout.session.completed"    //

)

func paymentHook(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		ipAddress := cleanIp(r.RemoteAddr)
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

		if evt.Type == CUSTOMER_CREATED { //1

			pcus, err := customerParser(evt.Data)
			fmt.Println("1-->", CUSTOMER_CREATED, err, pcus.Name, pcus.Email, pcus.Address.PostalCode, pcus.Address.Country)

		} else if evt.Type == CUSTOMER_SUBSCRIPTION_CREATED { //3

			pSubs, err := subscriptionParser(evt.Data)
			fmt.Println("2-->", CUSTOMER_SUBSCRIPTION_CREATED, err, pSubs.CurrentPeriodStart, pSubs.CurrentPeriodEnd)

		} else if evt.Type == INVOICE_CREATED { //2
			fmt.Println("3-->", INVOICE_CREATED)

		} else if evt.Type == CHARGE_SUCCEEDED { //4

			pCharge, err := chargeParser(evt.Data)
			fmt.Println("4-->", CHARGE_SUCCEEDED, err, evt.ID)
			if err == nil {
				//checkout_session
				fmt.Println("4-->", pCharge.ReceiptUrl)
				durl, err := stripeReceiptToPdfUrl(pCharge.ReceiptUrl)
				if err == nil {
					filename, err := DownloadFile("data/invoice", durl)
					logError("unableToDownloadInv", err)
					if err == nil {
						addFileStore("doc_keeper", pCharge.Invoice, "pdf", filename, evt.ID)
					}
				}
			}

		} else if evt.Type == INVOICE_UPDATED { //6
			fmt.Println("5-->", INVOICE_UPDATED)

		} else if evt.Type == INVOICE_PAID { //7

			inv, err := invoiceParser(evt.Data)
			fmt.Println("6-->", INVOICE_PAID, err, inv.Customer, inv.AmountPaid, inv.Subscription2, inv.PaymentIntent, inv.CustomerEmail, inv.CustomerName, inv.Number, inv.CustomerPhone)

		} else if evt.Type == INVOICE_PAYMENT_SUCCEEDED { //8
			fmt.Println("7-->", INVOICE_PAYMENT_SUCCEEDED, evt.ID)

		} else if evt.Type == CHECKOUT_SESSION_COMPLETED {

			fmt.Println("8.1-->", CHECKOUT_SESSION_COMPLETED, evt.ID)
			pSession, err := checkoutSessionParser(evt.Data)
			if err == nil {

				fmt.Println("8.2->", pSession.Invoice, pSession.CancelUrl, pSession.Subscription2)
				time.Sleep(time.Millisecond * 1500)
				iurl := stripeInvoiceReceiptUrl(pSession.Invoice)
				sql := fmt.Sprintf("UPDATE %s SET receipt_url=%q, payment_status=%q, doc_status=%q,update_date=%q WHERE doc_number=%q;", tableToBucket("doc_keeper"), iurl, pSession.PaymentStatus, pSession.Status, mtool.TimeNow(), pSession.ClentReferenceId)
				_, err = database.DB.Exec(sql)
				logError("checkoutSessionSuccessERR", err)
				if err != nil {
					log.Println(sql)
				}

				sql = fmt.Sprintf("SELECT login_id,account_id,total_payable FROM %s WHERE doc_number=%q;", tableToBucket("doc_keeper"), pSession.ClentReferenceId)
				row, _ := singleRow(sql)
				loginId, _ := row["login_id"].(string)
				accountId, _ := row["account_id"].(string)
				totalPayable, _ := row["total_payable"].(string)

				licenseKey := uuid.NewV1().String()
				subscriptionStart, subscriptionEnd := subscriptionStartEnd()
				addSubscription(accountId, pSession.Customer, licenseKey, "monthly", totalPayable, pSession.PaymentStatus, subscriptionStart, subscriptionEnd, "")

				docNumber := stripeInvoiceToNumber(pSession.Invoice)
				addDocKeeper("invoice", "sales", pSession.ClentReferenceId, docNumber, "", pSession.Status, "", "", totalPayable, loginId, accountId, ipAddress)

				sql = fmt.Sprintf("UPDATE %s SET reference=%q WHERE owner_table='doc_keeper' AND reference=%q;", tableToBucket("file_store"), docNumber, pSession.Invoice)
				database.DB.Exec(sql)

				//pSession.ClentReferenceId == doc_number
				invoice := pSession.ClentReferenceId
				dmap := docNumberToAccountInfo(pSession.ClentReferenceId)
				name, _ := dmap["account_name"].(string)
				receiptUrl, _ := dmap["receipt_url"].(string)
				amount, _ := dmap["total_payable"].(string)
				email := pSession.CustomerEmail
				err = paymentConfirmationEmail(email, name, amount, invoice, receiptUrl)
				logError("paymentConfirmationEmail", err)
				err = salesEmail(pSession.CustomerEmail, name, licenseKey)
				logError("salesEmail", err)
			}

		} else {
			log.Println("--> Unknown event.Type")
		}

		err = lxql.InsertUpdateObject("event", evt.ID, &evt, database.DB)
		logError("eventInsertERR", err)
		w.WriteHeader(http.StatusOK)

	} else {
		fmt.Fprintln(w, r.RemoteAddr)
	}
}
