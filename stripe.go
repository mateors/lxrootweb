package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mateors/mtool"
)

// unixToDateTime 1642614110, "2006-01-02 15:04:05"
func unixToDateTime(sec int64) string {
	tm := time.Unix(sec, 0)
	return tm.Format("2006-01-02 15:04:05") //dateTimeFormat
}

func requestBalance(secretKey string) (map[string]interface{}, error) {

	url := "https://api.stripe.com/v1/balance"
	method := "GET"

	//payload := strings.NewReader("email=sunzida%40gmail.com&name=Dr%20Sanzida&phone=01765110255&source=tok_1IU9EmEJdAixg8N3xAwmFVBP")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	//req.Header.Add("Stripe-Account", "acct_1Es5BxEJdAixg8N3")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	// fmt.Println("output")
	// for key, val := range result {
	// 	fmt.Println(key, val)
	// }
	return result, nil
}

func listAllPrices(secretKey string) (map[string]interface{}, error) {

	url := "https://api.stripe.com/v1/prices"
	method := "GET"

	//payload := strings.NewReader("email=sunzida%40gmail.com&name=Dr%20Sanzida&phone=01765110255&source=tok_1IU9EmEJdAixg8N3xAwmFVBP")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	//req.Header.Add("Stripe-Account", "acct_1Es5BxEJdAixg8N3")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	// fmt.Println("output")
	// for key, val := range result {
	// 	fmt.Println(key, val)
	// }
	return result, nil
}

func apiGetRequest(rurl, secretKey string) (map[string]interface{}, error) {

	req, err := http.NewRequest(http.MethodGet, rurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var rmap = make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&rmap)

	//fmt.Println("res.StatusCode:", res.StatusCode, res.StatusCode != http.StatusOK)
	//error handling
	if res.StatusCode != http.StatusOK {
		//fmt.Printf("%v %T\n", rmap, rmap)
		erow, isOk := rmap["error"].(map[string]interface{})
		if isOk {
			etype, isOk := erow["type"].(string)
			if isOk {
				errTypes := []string{"invalid_request_error", "idempotency_error", "card_error", "api_error"}
				if mtool.ArrayValueExist(errTypes, etype) {
					message, _ := erow["message"].(string)
					return nil, errors.New(message)
				}
			}
		}
	}
	return rmap, nil
}

func apiPostRequest(rurl, secretKey string, fmap map[string]string) (map[string]interface{}, error) {

	var fv = make(url.Values)
	for key, val := range fmap {
		fv.Set(key, val)
	}
	payload := strings.NewReader(fv.Encode())
	req, err := http.NewRequest("POST", rurl, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var rmap = make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&rmap)

	//error handling
	if res.StatusCode != http.StatusOK {
		etype, isOk := rmap["type"].(string)
		if isOk {
			errTypes := []string{"invalid_request_error", "idempotency_error", "card_error", "api_error"}
			if mtool.ArrayValueExist(errTypes, etype) {
				message, _ := rmap["message"].(string)
				return nil, errors.New(message)
			}
		}
	}

	return rmap, nil
}

// when we create price it automatically added under a product
func createSubscriptionPrice(stripeSecretKey, productName, priceInCents string) (map[string]interface{}, error) {

	//type = recurring
	//stripeKey := "sk_test_51OjqyFJFUQv2NTJsitgDUhNX3CPbns3eE3IyxSdTc8yEhI5p24SDyn9lyEI4AqaMSRghw6V25XoStkYa8Zl7zEOg006vuF1cTQ"
	var fmap = make(map[string]string)
	fmap["currency"] = "usd"
	fmap["unit_amount"] = priceInCents
	fmap["recurring[interval]"] = "month"
	fmap["product_data[name]"] = productName //product.name
	fmap["nickname"] = productName
	return apiPostRequest("https://api.stripe.com/v1/prices", stripeSecretKey, fmap)
}

// stripe checkout session
func createSession(stripeSecretKey, docNumber, customerEmail, priceId, qty string) (map[string]interface{}, error) {

	//type = recurring
	//stripeKey := "sk_test_51OjqyFJFUQv2NTJsitgDUhNX3CPbns3eE3IyxSdTc8yEhI5p24SDyn9lyEI4AqaMSRghw6V25XoStkYa8Zl7zEOg006vuF1cTQ"
	successUrl := "https://lxroot.com/complete"
	var fmap = make(map[string]string)
	fmap["client_reference_id"] = docNumber

	fmap["line_items[0][price]"] = priceId
	fmap["line_items[0][quantity]"] = qty
	fmap["line_items[1][price]"] = "price_1PPddlJFUQv2NTJs4sxm013J"
	fmap["line_items[1][quantity]"] = "2"

	fmap["success_url"] = successUrl
	fmap["customer_email"] = customerEmail
	fmap["mode"] = "subscription" //payment,setup,subscription
	//fmap["customer_creation"]="always" //payment
	//set payment_intent_data.setup_future_usage to have Checkout automatically
	return apiPostRequest("https://api.stripe.com/v1/checkout/sessions", stripeSecretKey, fmap)
}

func customerParser(evtData map[string]interface{}) (*Customer, error) {

	obj, isOk := evtData["object"]
	if !isOk {
		return nil, errors.New("wrong object")
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var cus = &Customer{}
	err = json.Unmarshal(jsonBytes, cus)
	return cus, err
}

func subscriptionParser(evtData map[string]interface{}) (*Subscription, error) {

	obj, isOk := evtData["object"]
	if !isOk {
		return nil, errors.New("wrong object")
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var subs = &Subscription{}
	err = json.Unmarshal(jsonBytes, subs)
	return subs, err
}

func invoiceParser(evtData map[string]interface{}) (*Invoice, error) {

	obj, isOk := evtData["object"]
	if !isOk {
		return nil, errors.New("wrong object")
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var inv = &Invoice{}
	err = json.Unmarshal(jsonBytes, inv)
	return inv, err
}

func paymentIntentParser(evtData map[string]interface{}) (*PaymentIntent, error) {

	obj, isOk := evtData["object"]
	if !isOk {
		return nil, errors.New("wrong object")
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var payint = &PaymentIntent{}
	err = json.Unmarshal(jsonBytes, payint)
	return payint, err
}

func chargeParser2(obj interface{}) (*Charge, error) {

	// obj, isOk := evtData["object"]
	// if !isOk {
	// 	return nil, errors.New("wrong object")
	// }
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var chrg = &Charge{}
	err = json.Unmarshal(jsonBytes, chrg)
	return chrg, err
}

func sessionParser(evtData map[string]interface{}) (*Session, error) {

	obj, isOk := evtData["object"]
	if !isOk {
		return nil, errors.New("wrong object")
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var sess = &Session{}
	err = json.Unmarshal(jsonBytes, sess)
	return sess, err
}

func chargeParser(evtData map[string]interface{}) (*Charge, error) {

	obj, isOk := evtData["object"]
	if !isOk {
		return nil, errors.New("wrong object")
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var chrg = &Charge{}
	err = json.Unmarshal(jsonBytes, chrg)
	return chrg, err
}

func stripeWebHook(w http.ResponseWriter, r *http.Request) {

	//r.ParseForm()
	//fmt.Println(r.Method)
	if r.Method == "POST" {

		var cid string = COMPANY_ID
		//code snippet
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		//print result
		//bodyString := string(bodyBytes)
		//fmt.Println(bodyString)
		var evt Event
		err = json.Unmarshal(bodyBytes, &evt)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("---------------------------------------------------", cid)
		fmt.Println("Time:", mtool.TimeNow())
		//fmt.Println("ApiVersion", evt.ApiVersion)
		fmt.Println("created", evt.Created, unixToDateTime(evt.Created))
		//fmt.Println()
		//fmt.Println("data", evt.Data)
		//fmt.Println()
		fmt.Println("id", evt.ID)
		//fmt.Println("livemode", evt.Livemode)
		//fmt.Println("object", evt.Object)
		fmt.Println("pending_webhooks", evt.PendingWebhooks)
		//fmt.Println("request", evt.Request)
		fmt.Println("type", evt.Type)
		//fmt.Println()

		if evt.Type == "customer.created" {

			cus, err := customerParser(evt.Data)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(err, cus)
			//customerEmail := sess.Customer_email
			// acc, err := emailToAccountDetails(cus.Email)
			// if err != nil {
			// 	log.Println(err)
			// }

			// sql := fmt.Sprintf(`UPDATE %s SET account_name="%s",country="%s",customid="%s" WHERE META().id="%s";`,
			// 	tableToBucket(AccountTable), cus.Name, cus.Address.Country, cus.ID, acc.ID)
			// //fmt.Println(sql)
			// pres := db.Query(sql)
			// if pres.Status != "success" {
			// 	log.Println(pres.Errors)
			// }

			// form := make(url.Values)
			// form.Set("model", "SubscriberStatus")
			// form.Set("account", acc.ID)
			// form.Set("email", cus.Email)
			// form.Set("cid", cid)
			// form.Set("customer", cus.ID)
			// form.Set("total_received", "0")
			// form.Set("remarks", "")
			// form.Set("status", "0")
			// sDocID := dynamicInsert(form)
			// fmt.Println(sDocID)

		} else if evt.Type == "customer.subscription.created" {

			subs, err := subscriptionParser(evt.Data)
			fmt.Println(err, subs)
			// startDate := unixToDateTime(subs.Current_period_start)
			// endDate := unixToDateTime(subs.Current_period_end)

			// sql := fmt.Sprintf(`UPDATE %s SET subscriber="%s",subscription_start="%s",subscription_end="%s",latest_invoice="%s",update_date="%s" WHERE customer="%s";`,
			// 	tableToBucket(SubscriberStatusTable), subs.ID, startDate, endDate, subs.Latest_invoice, mtool.TimeNow(), subs.Customer)
			// pres := db.Query(sql)
			// if pres.Status != "success" {
			// 	log.Println(pres.Errors)
			// 	fmt.Println(sql)
			// }

		} else if evt.Type == "invoice.created" {

			inv, err := invoiceParser(evt.Data)
			//fmt.Println(err, inv)
			fmt.Println(err, inv.AmountPaid)
			fmt.Println(err, inv.AmountDue)
			fmt.Println(err, inv.AmountRemaining)

		} else if evt.Type == "invoice.finalized" {

			fmt.Println("invoice.finalized yet to HANDLE")

		} else if evt.Type == "payment_intent.created" {

			payint, err := paymentIntentParser(evt.Data)
			fmt.Println(err, payint)

		} else if evt.Type == "payment_intent.succeeded" {

			fmt.Println("payment_intent.succeeded HANDLE")
			payint, err := paymentIntentParser(evt.Data)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(payint)
			// err = dInsert("PaymentIntent", &payint)
			// if err != nil {
			// 	log.Println(err)
			// }
			// if err == nil {

			// 	charges := payint.Charges
			// 	charge := charges.Data[0]
			// 	chr, err := chargeParser2(&charge)
			// 	if err != nil {
			// 		log.Println("ERR", err)
			// 		return
			// 	}
			// 	fmt.Println("Amount", chr.Amount)
			// 	//fmt.Println("chr", chr.Amount_captured)
			// 	fmt.Println("ID", chr.ID)
			// 	fmt.Println("Invoice", chr.Invoice)
			// 	fmt.Println("PIntent", chr.Payment_intent)
			// 	fmt.Println("Receipt", chr.Receipt_url)
			// 	fmt.Println("Status", chr.Status)
			// 	//fmt.Println("chr.Billing_details", chr.Billing_details)
			// 	fmt.Println("Email", chr.Billing_details.Email)
			// 	fmt.Println("Name", chr.Billing_details.Name)
			// 	fmt.Println("Country", chr.Billing_details.Address.Country)
			// 	//fmt.Println("PaymentMethodDetials", chr.PaymentMethodDetials)
			// 	cardInfo := chr.PaymentMethodDetials["card"].(map[string]interface{})
			// 	for key, val := range cardInfo {
			// 		fmt.Println(key, val)
			// 	}

			// 	acc, err := emailToAccountDetails(chr.Billing_details.Email)
			// 	if err != nil {
			// 		log.Println(err)
			// 	}
			// 	sql := fmt.Sprintf(`UPDATE %s SET account="%s",remarks="%s",payment_intent="%s",amount=%v,total_received=total_received+%d,update_date="%s",status=1 WHERE customer="%s";`,
			// 		tableToBucket("subscriber_status"),
			// 		acc.ID,
			// 		chr.Status,
			// 		chr.Payment_intent,
			// 		chr.Amount,
			// 		chr.Amount,
			// 		mtool.TimeNow(),
			// 		chr.Customer)
			// 	db.Query(sql)
			// 	fmt.Println(sql)

			// 	//send email
			// 	var eInfo emailInfo
			// 	eInfo.From = "support@complimention.com" //system owner
			// 	eInfo.Name = chr.Billing_details.Name
			// 	eInfo.Type = acc.AccountType
			// 	//eInfo.Token = sID.String()
			// 	emailSubject, emailBody := signupEmailPrepare(&eInfo, cid)
			// 	sendMail(emailSubject, emailBody, []string{chr.Billing_details.Email})
			// 	//fmt.Println(isSent, emailSubject, emailBody)
			// }

		} else if evt.Type == "charge.succeeded" {

			chrg, err := chargeParser(evt.Data)
			fmt.Println(err, chrg)
			fmt.Println(chrg.Billing_details)
			fmt.Println(chrg.Billing_details.Email)
			fmt.Println(chrg.Paid)

		} else if evt.Type == "invoice.payment_succeeded" {

			fmt.Println("invoice.payment_succeeded yet to HANDLE")

		} else if evt.Type == "invoice.paid" {

			inv, err := invoiceParser(evt.Data)
			fmt.Println(err, inv)

		} else if evt.Type == "invoice.updated" {

			fmt.Println("invoice.updated yet to HANDLE")

		} else if evt.Type == "checkout.session.completed" {

			sess, err := sessionParser(evt.Data)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(err, sess)
			// fmt.Println(sess.CustomerDetails.Email)
			// fmt.Println(sess.Payment_status) //paid
			// fmt.Println(sess.Subscription)
			// fmt.Println(sess.Status) //complete
			// //sess.Customer
			// // total_received
			// sql := fmt.Sprintf(`UPDATE %s SET amount=%v,checkout_session_status="%s",payment_status="%s",update_date="%s" WHERE customer="%s";`,
			// 	tableToBucket(SubscriberStatusTable),
			// 	sess.Amount_total,
			// 	sess.Status,
			// 	sess.Payment_status,
			// 	mtool.TimeNow(),
			// 	sess.Customer)

			// pres := db.Query(sql)
			// //fmt.Println(sql)
			// if pres.Status != "success" {
			// 	log.Println(pres.Errors)
			// 	fmt.Println(sql)
			// }

		} else {

			fmt.Println(evt.Type, "NOT PROCESSED YET")
		}

		//fmt.Fprintln(w, "OK")
		//dInsert("Event", &evt)
		w.WriteHeader(http.StatusOK)

	} else {
		fmt.Fprintln(w, r.RemoteAddr)
	}
}
