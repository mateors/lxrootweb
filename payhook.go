package main

import (
	"encoding/json"
	"fmt"
	"log"
	"lxrootweb/database"
	"lxrootweb/lxql"
	"net/http"

	"github.com/mateors/mtool"
	uuid "github.com/satori/go.uuid"
)

const (
	PAYMENT_INENT_CREATED      = "payment_intent.created"     //1 ok
	INVOICE_CREATED            = "invoice.created"            //2 ok
	INVOICE_FINALIZED          = "invoice.finalized"          //3 ok
	CHARGE_SUCCEEDED           = "charge.succeeded"           //4 ok
	PAYMENT_INTENT_SUCCEEDED   = "payment_intent.succeeded"   //5 ok
	INVOICE_UPDATED            = "invoice.updated"            //6 ok
	INVOICE_PAID               = "invoice.paid"               //7 ok
	INVOICE_PAYMENT_SUCCEEDED  = "invoice.payment_succeeded"  //8 ok
	CHECKOUT_SESSION_COMPLETED = "checkout.session.completed" //
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

			pCharge, err := chargeParser(evt.Data)
			fmt.Println("4-->", CHARGE_SUCCEEDED, err)
			if err == nil {

				//pCharge.BillingDetails.Address.Country
				//checkout_session
				durl, err := stripeReceiptToPdfUrl(pCharge.ReceiptUrl)
				if err == nil {
					//fmt.Println(err, durl)
					filename, err := DownloadFile("data/invoice", durl)
					logError("unableToDownloadInv", err)
					if err == nil {
						//pCharge.ReceiptUrl
						//fmt.Println(pCharge.BillingDetails.Email, filename)
						docNumber, err := emailToDocNumber(pCharge.BillingDetails.Email)
						if err == nil {
							sql := fmt.Sprintf("UPDATE %s SET receipt_url=%q, update_date=%q WHERE doc_number=%q;", tableToBucket("doc_keeper"), pCharge.ReceiptUrl, mtool.TimeNow(), docNumber)
							lxql.RawSQL(sql, database.DB)

							remarks := evt.ID
							addFileStore("doc_keeper", docNumber, "pdf", filename, remarks)
						}
					}
				}
			}

		} else if evt.Type == PAYMENT_INTENT_SUCCEEDED { //5
			fmt.Println("5-->", PAYMENT_INTENT_SUCCEEDED)

		} else if evt.Type == INVOICE_UPDATED { //6
			fmt.Println("6-->", INVOICE_UPDATED)

		} else if evt.Type == INVOICE_PAID { //7

			//fmt.Println("7-->", INVOICE_PAID)
			//inv, err := invoiceParser(evt.Data)
			//fmt.Println(err, inv.Customer, inv.AmountPaid, inv.Subscription2, inv.PaymentIntent, inv.CustomerEmail, inv.CustomerName, inv.Number, inv.CustomerPhone)
			//fmt.Println(inv.InvoicePdf, inv.HostedInvoiceUrl)
			//inv.PeriodStart

		} else if evt.Type == INVOICE_PAYMENT_SUCCEEDED { //8
			fmt.Println("8-->", INVOICE_PAYMENT_SUCCEEDED)

		} else if evt.Type == CHECKOUT_SESSION_COMPLETED {

			fmt.Println("9-->", CHECKOUT_SESSION_COMPLETED)
			pSession, err := checkoutSessionParser(evt.Data)
			if err == nil {

				sql := fmt.Sprintf("UPDATE %s SET payment_status=%q, doc_status=%q WHERE doc_number=%q;", tableToBucket("doc_keeper"), pSession.PaymentStatus, pSession.Status, pSession.ClentReferenceId)
				err = lxql.RawSQL(sql, database.DB)
				logError("checkoutSessionSuccessERR", err)

				sql = fmt.Sprintf("SELECT login_id,account_id,total_payable FROM %s WHERE doc_number=%q;", tableToBucket("doc_keeper"), pSession.ClentReferenceId)
				row, _ := singleRow(sql)
				loginId, _ := row["login_id"].(string)
				accountId, _ := row["account_id"].(string)
				totalPayable, _ := row["total_payable"].(string)

				licenseKey := uuid.NewV1().String()
				subscriptionStart := ""
				subscriptionEnd := ""
				addSubscription(accountId, pSession.Customer, licenseKey, "monthly", totalPayable, pSession.PaymentStatus, subscriptionStart, subscriptionEnd, "")

				docNumber := pSession.Invoice
				addDocKeeper("invoice", "sales", pSession.ClentReferenceId, docNumber, "", pSession.Status, "", "", totalPayable, loginId, accountId)

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
