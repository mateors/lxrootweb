package main

import (
	"encoding/json"
	"fmt"
	"log"
	"lxrootweb/database"
	"lxrootweb/lxql"
	"lxrootweb/utility"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/mateors/mtool"
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
)

func GetBaseURL(r *http.Request) string {

	var baseurl, proto string
	fproto := r.Header.Get("X-Forwarded-Proto")
	proto = "http"
	if fproto == "https" {
		proto = "https"
	} else if r.TLS != nil {
		proto = "https"
	}
	baseurl = fmt.Sprintf("%s://%s", proto, r.Host)
	return baseurl
}

func homePage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/section_getstarted.gohtml",
			"templates/footer.gohtml",
			"wpages/home.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Application Hosting Platform | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func support(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/support.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Support | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func features(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		slug := chi.URLParam(r, "slug")

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/features.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		features := []map[string]interface{}{
			{"id": 1, "name": "Server management", "icon": "devices", "slug": "server-management"},
			{"id": 2, "name": "Security", "icon": "verified_user", "slug": "security"},
			{"id": 3, "name": "Web server", "icon": "vertical_split", "slug": "web-server"},
			{"id": 4, "name": "Backups", "icon": "compost", "slug": "backups"},
			{"id": 5, "name": "Databases", "icon": "layers", "slug": "databases"},
			{"id": 6, "name": "DNS", "icon": "dns", "slug": "dns"},

			{"id": 7, "name": "Runtimes", "icon": "terminal", "slug": "runtimes"},
			{"id": 8, "name": "Applications", "icon": "apps", "slug": "applications"},
			{"id": 16, "name": "Developer tools", "icon": "construction", "slug": "developer-tools"},

			{"id": 9, "name": "Website management", "icon": "web", "slug": "website-management"},
			{"id": 10, "name": "Email", "icon": "email", "slug": "email"},
			{"id": 11, "name": "Domains", "icon": "public", "slug": "domains"},
			{"id": 12, "name": "Client management", "icon": "manage_accounts", "slug": "client-management"},
			{"id": 13, "name": "User interface", "icon": "grid_view", "slug": "user-interface"},
			{"id": 14, "name": "Integrations", "icon": "integration_instructions", "slug": "integrations"},
			{"id": 15, "name": "Resellers", "icon": "groups_3", "slug": "resellers"},
		}

		fRows := []map[string]interface{}{
			{"feature_id": 1, "title": "Clustering support", "desc": "Deploy 1 - 10,000 servers to your cluster and manage them all from a single User Interface."},
			{"feature_id": 1, "title": "Disaster recovery", "desc": "Restore application to another server within your cluster from the last available backup."},
			{"feature_id": 1, "title": "Simple server deployment", "desc": "Provision new servers to your cluster with a single command."},
			{"feature_id": 1, "title": "Zero config scaling", "desc": "Each server automatically inherits your global settings - make changes per-server where needed."},
			{"feature_id": 1, "title": "WordPress toolkit", "desc": "Install themes and plugins, configure settings, update and manage."},

			//{"feature_id": 2, "title": "Role containerisation", "desc": "Every role is containerised, including all components of the email system. Containers have no access to website files, even if the application role is installed on the same server."},
			{"feature_id": 2, "title": "Change tracker", "desc": "Track your source code changes and restore them if needed."},
			{"feature_id": 2, "title": "App containerisation", "desc": "Each stand alone application is containerised to protect other apps in your cluster from any vulnerabilities."},
			{"feature_id": 2, "title": "Brute force protection", "desc": "Configure rate limits and block/allow lists for email addresses and IPs to help protect your cluster."},
			{"feature_id": 2, "title": "ModSecurity (OWASP)", "desc": "Enable or disable ModSecurity for each server, customise default OWASP settings, add custom rules using an update the running OWASP version."},
			{"feature_id": 2, "title": "SSL / TLS", "desc": "All applications and services on LxRoot fully support SSL/TLS."},

			{"feature_id": 3, "title": "NGINX", "desc": "The Nginx web server is a light weight open-source web server. Nginx web server allows for cache exclusion, purge cache, URL rewriting and FastCGI cache to be configured on a per domain basis."},

			{"feature_id": 4, "title": "Incremental backups", "desc": "A powerful incremental application backup system utilising btrfs is built into the LxRoot platform."},
			{"feature_id": 4, "title": "S3 compatible backups", "desc": "Set up an external backup system. It can work independently or alongside LxRoot’s built-in system."},
			{"feature_id": 4, "title": "Granular backup restores", "desc": "Restore individual databases, email accounts or whole applications."},
			{"feature_id": 4, "title": "On-demand backups", "desc": "Create on demand full application and email backups in the click of a button."},
			{"feature_id": 4, "title": "Schedule backups", "desc": "Configure automatic cluster wide application backups specifying frequency, time and more."},

			{"feature_id": 5, "title": "MySQL", "desc": "MySQL is an open-source relational database management system (RDBMS)"},
			{"feature_id": 5, "title": "PostgreSQL", "desc": "PostgreSQL, also known as Postgres, is a free and open-source relational database management system emphasizing extensibility and SQL compliance."},
			{"feature_id": 5, "title": "Couchbase", "desc": "Couchbase is a distributed NoSQL cloud database that delivers versatility, performance, scalability, and financial value for all of your cloud, mobile, AI, and edge computing"},

			{"feature_id": 6, "title": "PowerDNS", "desc": "Install the DNS service and have everything work automatically without any additional configuration."},
			{"feature_id": 6, "title": "Unlimited DNS servers", "desc": "Run a decentralised DNS cluster of any size with all zones automatically synched between them."},
			{"feature_id": 6, "title": "DNS management", "desc": "Utilise simple but powerful DNS management tools, powered by the open source PowerDNS"},
			{"feature_id": 6, "title": "DNS templating", "desc": "Create DNS records which are automatically added to every zone created within your hosting platform."},
			{"feature_id": 6, "title": "Domain re-mapping tool", "desc": "Map or move domains to point to specific web roots within your hosting space."},
			{"feature_id": 6, "title": "Cloudflare", "desc": "Synchronise your domain's DNS and set your proxy status with Cloudflare."},
			{"feature_id": 6, "title": "DNSSEC", "desc": "Enable/disable DNSSEC on a per domain basis to authenticate existing DNS records."},
			{"feature_id": 6, "title": "Custom nameservers", "desc": "Add custom nameservers such as ns1.your-domain.com and ns2.your-domain.com"},

			{"feature_id": 7, "title": "HTML5", "desc": "Provide hosting for all versions of HTML and HTML5 files."},
			{"feature_id": 7, "title": "PHP", "desc": "Run different php versions on a per website basis, with 5.6, 7.4, 8.0, 8.1, 8.2 and 8.3 supported."},
			{"feature_id": 7, "title": "Golang", "desc": "Golang environment always uses the latest stable Go 1.x version. You can’t set a different version unless you deploy a docker image."},
			{"feature_id": 7, "title": "Rust", "desc": "Always uses the latest stable rust version."},
			{"feature_id": 7, "title": "NodeJS", "desc": "Always uses the latest stable nodejs version."},
			{"feature_id": 7, "title": "Java", "desc": "Always uses the latest stable version."},
			{"feature_id": 7, "title": "Python3", "desc": "Always uses the latest stable version."},
			{"feature_id": 7, "title": "Ruby", "desc": "Always uses the latest stable version."},
			{"feature_id": 7, "title": "Docker", "desc": "Deploy a Docker image on LxRoot, it can use virtually any programming language and framework."},

			{"feature_id": 8, "title": "Lxpkg", "desc": "Linux package manager, assists you in managing packages, listing services, and handling process and network services."},
			{"feature_id": 8, "title": "Cron Job", "desc": "Tme-based job scheduler to automate common task."},
			{"feature_id": 8, "title": "Backup", "desc": "Helps you to take automatic and on-demand data backups."},
			{"feature_id": 8, "title": "Cert", "desc": "Manage third party TLS certificates on LxRoot."},
			{"feature_id": 8, "title": "Webhooks", "desc": "Deploy your application directly from github using webhooks."},
			{"feature_id": 8, "title": "Cloner", "desc": "Make a carbon copy of one of your existing application."},
			{"feature_id": 8, "title": "Build", "desc": "Install, build, run, execute and set linux environment variables."},
			{"feature_id": 8, "title": "WordPress", "desc": "Install wordpress with just a single click wordpress installer."},
			{"feature_id": 8, "title": "Ticket", "desc": "Built-in customer support ticketing system ensures seamless customer support."},

			{"feature_id": 9, "title": "File manager", "desc": "Our inline file manager is fully-featured and the interface easily compatible with tablets, phones, as well as desktops and laptops. Drag and drop files, easily perform mass operations and so on."},
			{"feature_id": 9, "title": "Clone a website", "desc": "Clone an existing site or add a site with an application pre-installed. Cloning automatically duplicates the database and updates website configuration files and paths."},
			{"feature_id": 9, "title": "SFTP accounts", "desc": "Add, remove and manage SFTP or FTP accounts."},
			{"feature_id": 9, "title": "Cron jobs", "desc": "Schedule common tasks to automatically run at a specified date and time."},
			{"feature_id": 9, "title": "SSH access", "desc": "Connect with SSH keys or SSH password authentication."},
			{"feature_id": 9, "title": "Transfer website ownership", "desc": "Seamlessly move website from your account to a customer's."},
			{"feature_id": 9, "title": "IP manager", "desc": "Allow or block individual IP addresses from accessing your website"},
			{"feature_id": 9, "title": "Bulk management", "desc": "Quickly and easily perform common upkeep and management tasks across multiple websites and accounts."},
			{"feature_id": 9, "title": "Vhost configuration", "desc": "Enter a custom configuration for each domain on a website."},
			{"feature_id": 9, "title": "Opcode caching", "desc": "Activate opcode caching native to LxRoot. Simply toggle on/off for individual application."},
			{"feature_id": 9, "title": "cPanel importer", "desc": "Quickly and easily move websites, emails and all passwords from cPanel using an account backup or SCP."},
			{"feature_id": 9, "title": "Plesk importer", "desc": "Import websites and email account from Plesk using a site backup."},

			{"feature_id": 10, "title": "Webmail (SnappyMail)", "desc": "Powered by Rainloop. Supporting password changes, spam filtering and other advanced features through the webmail interface."},
			{"feature_id": 10, "title": "Inbound spam filtering", "desc": "Industry standard Rspam based filtering. Granular control over scoring requirements and whitelisting/blacklisting"},
			{"feature_id": 10, "title": "SPF and DKIM authentication", "desc": "Quickly and easily configure and monitor the status of SPF and DKIM records."},
			{"feature_id": 10, "title": "System generated emails", "desc": "Send vital password reset and user invite emails using LxRoot mail transfer agent of configure your own custom SMTP."},
			{"feature_id": 10, "title": "Allow / Block lists", "desc": "Automatically accept or block inbound emails from specific email addresses."},
			{"feature_id": 10, "title": "DMARC authentication", "desc": "Automatically add a dns record for Domain-based Message Authentication and Reporting."},

			{"feature_id": 11, "title": "DNS editor", "desc": "Create, edit, and delete Domain Name System (DNS) zone records."},
			{"feature_id": 11, "title": "Alias, addon and sub domains", "desc": "Add and manage alias, addon and subdomains."},
			{"feature_id": 11, "title": "301 / 302 Redirects", "desc": "Add, edit and manage 301 and 302 redirects."},
			{"feature_id": 11, "title": "Change primary domain", "desc": "Modify the primary domain name for an existing website in just a few clicks."},
			{"feature_id": 11, "title": "Cloudflare", "desc": "Synchronise your domain's DNS and set your proxy status with Cloudflare."},
			{"feature_id": 11, "title": "DNSSEC", "desc": "Enable/disable DNSSEC on a per domain basis to authenticate existing DNS records."},

			{"feature_id": 12, "title": "Impersonation", "desc": "Access a customer's account and see exactly what they see without having to leave your account."},
			{"feature_id": 12, "title": "Upgrade / Downgrade package", "desc": "Seamlessly upgrade or downgrade a customer's package to another of your hosting packages."},
			{"feature_id": 12, "title": "Soft / Permanently delete", "desc": "Soft delete applications and customer accounts. Set a time frame for permanent removable from your cluster and restore on a whim."},
			{"feature_id": 12, "title": "Client Role", "desc": "We support three different role Admin, Reseller and User."},

			{"feature_id": 13, "title": "Single unified panel", "desc": "Manage all servers, customers and applications all in one place. Admins and users login to the same panel and see only the tools that are relevant to them."},
			{"feature_id": 13, "title": "Fully responsive", "desc": "Access LxRoot on desktop, tablet and mobile without any feature limitations."},
			{"feature_id": 13, "title": "Multi-account support", "desc": "Manage multiple distinct accounts from a single login."},
			{"feature_id": 13, "title": "Multi-language support", "desc": "We support English, Bangla, Spanish, Hindi and more."},
			{"feature_id": 13, "title": "Global search", "desc": "Quickly and easily find what you are looking for with a powerful estate wide search functionality."},

			{"feature_id": 14, "title": "Cloudflare", "desc": "Synchronise your domain's DNS and set your proxy status with Cloudflare."},
			{"feature_id": 14, "title": "GraphQL API", "desc": "Our powerful GraphQL API allows you to integrate with 3rd party systems you already use."},
			{"feature_id": 14, "title": "Slack notifications", "desc": "Receive notifications direct to a Slack channel of your choice."},

			{"feature_id": 15, "title": "Custom nameservers", "desc": "Create branded nameservers based on your company domain."},
			{"feature_id": 15, "title": "Custom packages", "desc": "Create packages with unique quotas and tool sets."},
			{"feature_id": 15, "title": "Whitelabel", "desc": "Customise the look and feel of the panel with logos, colours, button styling and more."},
			{"feature_id": 15, "title": "Custom domains", "desc": "Ability to set their own branded control panel, phpMyAdmin and webmail client."},
			{"feature_id": 15, "title": "Multi-tiered reselling", "desc": "With unlimited levels of reselling straight out the box, you can effortlessly reach new audiences."},

			{"feature_id": 16, "title": "php.ini editor", "desc": "Add and manage php directives on a platform and per website basis."},
			{"feature_id": 16, "title": "Document root", "desc": "Update the directory which stores the application files."},
			{"feature_id": 16, "title": "Nginx editor", "desc": "Add and manage nginx directives on a platform and per application basis."},
		}

		//DMARC record enables domain owners to protect their domains from unauthorized access and usage
		fmap := featureBySlug(slug, features)
		var rows = make([]map[string]interface{}, 0)
		if len(fmap) > 0 {
			fid := fmap["id"].(int)
			rows = findFeaturesById(fid, fRows)
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			FeatureName  string
			FeatureIcon  string
			FeatureRows  []map[string]interface{}
			Flinks       []map[string]interface{}
		}{
			Title:        "Features | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
			FeatureName:  fmt.Sprint(fmap["name"]),
			FeatureIcon:  fmt.Sprint(fmap["icon"]),
			FeatureRows:  rows,
			Flinks:       features,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func technology(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/section_getstarted.gohtml",
			"templates/footer.gohtml",
			"wpages/technology.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Technology | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func apphosting(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/section_getstarted.gohtml",
			"templates/footer.gohtml",
			"wpages/apphosting.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Application Hosting | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func roadmap(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/roadmap.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Product Roadmap | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func pricing(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/pricing.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Pricing | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func terms(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/terms.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Terms & Conditions | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func privacy(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/privacy.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Privacy Policy | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func shop(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		smap, _ := getSessionInfo(r)
		var yourname string = "Sign In"
		if len(smap) > 1 {
			yourname, _ = smap["account_name"].(string)
		}

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/shop.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		docNumber, err := getCookie("docid", r)
		if err == nil {
			docStatus := lxql.FieldByValue("doc_keeper", "doc_status", fmt.Sprintf("doc_number=%q", docNumber), database.DB)
			if docStatus != "pending" {
				delCookie("docid", r, w)
			}
		}

		var count int = docToCartCount(docNumber)
		base := GetBaseURL(r)

		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			CartCount    int
			Yourname     string
		}{
			Title:        "LxRoot Shop",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
			CsrfToken:    ctoken,
			CartCount:    count,
			Yourname:     yourname,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		r.ParseForm()
		var errNo int = 1
		var errMsg string
		commonDataSet(r)
		//fmt.Println("##", r.Form)

		funcsMap := map[string]interface{}{
			"validCSRF": validCSRF,
		}
		rmap := make(map[string]interface{})
		for key := range r.Form {
			rmap[key] = r.FormValue(key)
		}
		response := CheckMultipleConditionTrue(rmap, funcsMap)

		if response == "OKAY" {

			//
			itemId := r.FormValue("item") //item.item_code
			docNumber, err := getCookie("docid", r)
			if err == nil {
				fmt.Println("update doc_keeper table", docNumber)
			}

			qty := "1"
			docRef := visitorInfo(r, w)
			docId, err := addToCart(itemId, qty, docRef, docNumber, "", "")
			if err == nil {
				errNo = 0
				errMsg = "OK"
				setCookie("docid", docId, 24*86400, w)

			} else if err != nil {
				errNo = 1
				errMsg = err.Error()
			}

		} else {

			errNo = 1
			errMsg = response
		}

		var row = make(map[string]interface{})
		row["error"] = errNo
		row["message"] = errMsg
		bs, err := json.Marshal(row)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(bs))
	}
}

func complete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/complete.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			Yourname     string
		}{
			Title:        "LxRoot order complete",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
			Yourname:     "Sign In", //need to change
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func product(w http.ResponseWriter, r *http.Request) {

	smap, _ := getSessionInfo(r)

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/product.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		docNumber, err := getCookie("docid", r)
		if err == nil {
			docStatus := lxql.FieldByValue("doc_keeper", "doc_status", fmt.Sprintf("doc_number=%q", docNumber), database.DB)
			if docStatus != "pending" {
				delCookie("docid", r, w)
			}
		}

		sql := fmt.Sprintf(`SELECT t.id,t.item_id,t.quantity,t.price,t.tax_amount,t.payable_amount,t.trx_type,i.item_name,i.tags,d.doc_number,d.total_payable,d.total_discount 
							FROM lxroot._default.transaction_record t
							LEFT JOIN lxroot._default.doc_keeper d ON d.doc_number=t.doc_number
							LEFT JOIN lxroot._default.item i ON i.id=t.item_id
							WHERE t.doc_number="%s" AND d.doc_status="pending";`, docNumber)
		sql = cleanText(sql)
		rows, err := lxql.GetRows(sql, database.DB)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		var yourname string = "Sign In"
		if len(smap) > 1 {
			yourname, _ = smap["account_name"].(string)
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			Yourname     string
			CartCount    int
		}{
			Title:        "Product Details | LxRoot",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
			CsrfToken:    ctoken,
			Yourname:     yourname,
			CartCount:    len(rows),
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	} else if r.Method == http.MethodPost {

		itemId := r.FormValue("item") //item.item_code
		docNumber, err := getCookie("docid", r)
		if err == nil {
			fmt.Println("update doc_keeper table", docNumber)
		}

		qty := "1"
		docRef := visitorInfo(r, w)
		docId, err := addToCart(itemId, qty, docRef, docNumber, "", "")
		if err == nil {
			setCookie("docid", docId, 24*86400, w)
			http.Redirect(w, r, "/checkout", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/product", http.StatusSeeOther)
		return
	}
}

func checkout(w http.ResponseWriter, r *http.Request) {

	var loginRequired bool = true
	smap, err := getSessionInfo(r)
	if err == nil {
		loginRequired = false
	}

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(FuncMap).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/checkout.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		docNumber, err := getCookie("docid", r)
		if err == nil {
			docStatus := lxql.FieldByValue("doc_keeper", "doc_status", fmt.Sprintf("doc_number=%q", docNumber), database.DB)
			if docStatus != "pending" {
				delCookie("docid", r, w)
			}
		}

		sql := fmt.Sprintf(`SELECT t.id,t.item_id,t.quantity,t.price,t.tax_amount,t.payable_amount,t.trx_type,i.item_name,i.tags,d.doc_number,d.total_payable,d.total_discount 
							FROM lxroot._default.transaction_record t
							LEFT JOIN lxroot._default.doc_keeper d ON d.doc_number=t.doc_number
							LEFT JOIN lxroot._default.item i ON i.id=t.item_id
							WHERE t.doc_number="%s" AND d.doc_status="pending";`, docNumber)
		sql = cleanText(sql)
		rows, err := lxql.GetRows(sql, database.DB)
		if err != nil {
			log.Println(err)
			return
		}

		var totalPayable, totalDiscount string
		if len(rows) > 0 {
			totalPayable, _ = rows[0]["total_payable"].(string)
			totalDiscount, _ = rows[0]["total_discount"].(string)
			if totalDiscount == "0" {
				totalDiscount = ""
			}
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		var yourname string = "Sign In"
		if len(smap) > 1 {
			yourname, _ = smap["account_name"].(string)
		}

		base := GetBaseURL(r)
		data := struct {
			Title         string
			Base          string
			BodyClass     string
			MainDivClass  string
			CsrfToken     string
			Rows          []map[string]interface{} //cart items
			TotalPayable  string
			TotalDiscount string
			CartCount     int
			DocNumber     string
			LoginRequired bool
			Yourname      string
		}{
			Title:         "LxRoot Checkout",
			Base:          base,
			BodyClass:     "bg-white text-slate-700",
			MainDivClass:  "main min-h-[calc(100vh-52px)]",
			CsrfToken:     ctoken,
			Rows:          rows,
			TotalPayable:  totalPayable,
			TotalDiscount: totalDiscount,
			CartCount:     len(rows),
			DocNumber:     docNumber,
			LoginRequired: loginRequired,
			Yourname:      yourname,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		r.ParseForm()
		id := r.FormValue("id")
		todo := r.FormValue("todo")

		if strings.ToUpper(todo) == "DELETE" {

			//update docKepper table
			sql := fmt.Sprintf("SELECT doc_number,payable_amount FROM %s WHERE id=%q;", tableToBucket("transaction_record"), id)
			row, err := singleRow(sql)
			if err == nil {
				docNumber, _ := row["doc_number"].(string)
				payableAmount, _ := row["payable_amount"].(string)
				sql = fmt.Sprintf("UPDATE %s SET total_payable=total_payable-%s WHERE doc_number=%q;", tableToBucket("doc_keeper"), payableAmount, docNumber)
				lxql.RawSQL(sql, database.DB)

				sql := fmt.Sprintf("DELETE FROM %s WHERE id=%q;", tableToBucket("transaction_record"), id)
				err = lxql.RawSQL(sql, database.DB)
				logError("delCartItem", err)
			}

		} else if strings.ToUpper(todo) == "COUPON" {

			// code := r.FormValue("code") //discount code
			// docNumber, err := getCookie("docid", r)
			// if err == nil {
			// 	sql := fmt.Sprintf("SELECT * FROM %s WHERE doc_type='cart' AND status=1 AND doc_number=%q;", tableToBucket("doc_keeper"), docNumber)
			// 	row, err := singleRow(sql)
			// 	fmt.Println(err, row)
			// }

		} else if strings.ToUpper(todo) == "CHECKOUT" {

			if len(smap) == 0 {
				log.Println("No login session found to redirecting to checkout page...")
				http.Redirect(w, r, "/checkout", http.StatusSeeOther)
				return
			}

			loginId, _ := smap["id"].(string)
			accountId, _ := smap["account_id"].(string)
			customerEmail, _ := smap["username"].(string)
			docId := r.FormValue("docid")
			docNumber, err := getCookie("docid", r)
			logError("checkoutGetCookieERR", err)
			fmt.Println(loginId, accountId)

			if err == nil {

				row, err := createSession(utility.STRIPE_SECRETKEY, docNumber, customerEmail)
				if err != nil {
					log.Println("createSessionERROR:", err)
				}
				if err == nil {
					sessionId, _ := row["id"].(string) //checkout.session.id
					sql := fmt.Sprintf("UPDATE %s SET login_id=%q, account_id=%q, doc_status=%q, doc_description=%q WHERE doc_number=%q;", tableToBucket("doc_keeper"), loginId, accountId, "checkout_session", sessionId, docId)
					lxql.RawSQL(sql, database.DB)
					rurl, isOk := row["url"].(string)
					if isOk {
						fmt.Println("checkoutSQL>", sql)
						delCookie("docid", r, w)
						http.Redirect(w, r, rurl, http.StatusSeeOther)
					}
				}
			}
		}
		http.Redirect(w, r, "/checkout", http.StatusSeeOther)
		return
	}
}

func faqs(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/faqs.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "FAQs | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func about(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/about.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "About | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func contact(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/contact.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Contact | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func joinWaitlist(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/join-waitlist.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
		}{
			Title:        "Join-waitlist | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
			CsrfToken:    ctoken,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		var errNo int = 1
		var errMsg string = "OK"
		r.ParseForm()
		commonDataSet(r)

		ctoken, err := getCookie("ctoken", r) //cross check with form token
		logError("ctoken", err)
		r.Form.Set("ctoken", ctoken)

		funcsMap := map[string]interface{}{
			"validCSRF": validCSRF,
		}
		rmap := make(map[string]interface{})
		for key := range r.Form {
			rmap[key] = r.FormValue(key)
		}
		response := CheckMultipleConditionTrue(rmap, funcsMap)

		if response == "OKAY" {

			errNo = 0
			modelName := structName(WaitList{})
			err := modelUpsert(modelName, r.Form)
			if err == nil {

				email := r.FormValue("email")
				errMsg = settingsValue("waitlist_confirmation") //to display after form submission

				subject := "Thank You for Joining the LxRoot Waitlist!"
				emailTemplate := settingsValue("waitlist_email")

				dmap := make(map[string]interface{})
				dmap["first_name"] = r.FormValue("first_name")
				dmap["last_name"] = r.FormValue("last_name")
				emailBody, _ := templatePrepare(emailTemplate, dmap)
				err = SendEmail([]string{email}, subject, emailBody)
				logError("waitListSendEmail", err)
			}
		}

		var row = make(map[string]interface{})
		row["error"] = errNo
		row["message"] = errMsg
		bs, err := json.Marshal(row)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(bs))
	}
}

func signup(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/signup.gohtml", //signupm.gohtml
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		//setCookie("signup_success", "You have successfully registered with LxRoot.", 300, w)
		successMsg, _ := getCookie("signup_success", r)
		//fmt.Println(">>", err, successMsg)

		base := GetBaseURL(r)
		data := struct {
			Title          string
			Base           string
			BodyClass      string
			MainDivClass   string
			CsrfToken      string
			SuccessMessage string
			Yourname       string
		}{
			Title:          "Signup | LxRoot",
			Base:           base,
			BodyClass:      "",
			MainDivClass:   "main min-h-[calc(100vh-52px)]",
			CsrfToken:      ctoken,
			SuccessMessage: successMsg,
			Yourname:       "Sign In",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		r.ParseForm()
		var errNo int = 1
		var errMsg string
		commonDataSet(r)

		funcsMap := map[string]interface{}{
			"validCSRF":        validCSRF,
			"validSignupField": validSignupField,
			"validEmail":       validEmail,
		}
		rmap := make(map[string]interface{})
		for key := range r.Form {
			rmap[key] = r.FormValue(key)
		}

		response := CheckMultipleConditionTrue(rmap, funcsMap)
		if response == "ERROR email already exist" {
			errNo = 1
			errMsg = response

		} else if response == "ERROR password is required" {
			errNo = 2
			errMsg = response

		} else if response == "OKAY" {

			email := strings.ToLower(r.FormValue("email"))
			firstName := r.FormValue("first_name")
			lastName := r.FormValue("last_name")
			passwd := r.FormValue("passwd")
			ipAddress := cleanIp(mtool.ReadUserIP(r))

			accessName := "client"
			accessId := accessIdByName(accessName)
			username := email

			parentId := ""
			accountType := accessName
			accountName := fmt.Sprintf("%s %s", firstName, lastName)
			accountId, err := addAccount(parentId, accountType, email, accountName, firstName, lastName)
			if err != nil {
				errNo = 3
				errMsg = err.Error()
			}
			if err == nil {

				addAddress(accountId, "billing", "", "", "", "", "", "")
				loginId, _ := addLogin(accountId, accessId, accessName, username, passwd)
				code := uuid.NewV4().String()
				_, err = addVerification(username, "signup", code, "")
				logError("addVerification", err)

				//get the location usin thirdparty api
				location := getLocationWithinSec(ipAddress)
				if location == "" {
					location = fmt.Sprintf("IP: %s", ipAddress)
				}
				verfyUrl := fmt.Sprintf("https://lxroot.com/verify?email=%s&token=%s", email, code)
				err = signupEmail(email, accountName, location, verfyUrl)
				logError("signupEmail", err)

				logDetails := fmt.Sprintf("username %s signup", email)
				addActiviyLog(loginId, "INSERT", "account", "", logDetails, ipAddress)

				errNo = 0
				errMsg = "Congratulations! You have successfully registered with LxRoot. 🎉"
				setCookie("signup_success", "You have successfully registered with LxRoot.", 30, w)
			}

		} else {
			errNo = 9
			errMsg = "Unknown error! please report"
			log.Println(response)
		}

		var row = make(map[string]interface{})
		row["error"] = errNo
		row["message"] = errMsg
		bs, err := json.Marshal(row)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(bs))
	}
}

func verify(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		var message, errMsg string
		email := r.FormValue("email")
		token := r.FormValue("token")

		err = verifySignup(email, token)
		if err != nil {
			errMsg = err.Error()
		}

		if err == nil {

			//fmt.Println(ainfo, len(ainfo))
			ainfo, err := usernameToAccounInfo(email)
			if err == nil {

				accountId := ainfo["account_id"].(string)
				sql := fmt.Sprintf("UPDATE %s SET status=1,update_date=%q,remarks='verified' WHERE id=%q;", tableToBucket("account"), mtool.TimeNow(), accountId)
				database.DB.Exec(sql)

				loginId := ainfo["login_id"].(string)
				sql = fmt.Sprintf("UPDATE %s SET status=1 WHERE id=%q;", tableToBucket("login"), loginId)
				database.DB.Exec(sql)

				ipAddress := cleanIp(mtool.ReadUserIP(r))
				logDetails := fmt.Sprintf("username %s verified", email)
				addActiviyLog(loginId, "UPDATE", "account", "", logDetails, ipAddress)
				message = "You have been verified."
			}
		}

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/verify_signup.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title          string
			Base           string
			BodyClass      string
			MainDivClass   string
			ErrorMessage   string
			SuccessMessage string
			Yourname       string
		}{
			Title:          "LxRoot | Reset password",
			Base:           base,
			BodyClass:      "",
			MainDivClass:   "main min-h-[calc(100vh-52px)]",
			ErrorMessage:   errMsg,
			SuccessMessage: message,
			Yourname:       "Sign In",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func signin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		smap, err := getSessionInfo(r)
		if err == nil {
			if access, isOk := smap["access_name"].(string); isOk {
				vacc := []string{"superadmin", "admin", "client"}
				if mtool.ArrayValueExist(vacc, access) {
					http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				}
			}
		}

		sessionCode, err := getCookie("visitor_session", r)
		if err != nil {
			visitorInfo(r, w)
			log.Println("Visitor sessionCode generated", sessionCode)
		}

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/signin.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)
		errMessage := r.FormValue("error")

		var referer string
		purl, err := url.Parse(r.Referer())
		if err == nil {
			referer = purl.Path
			setCookie("redirect", referer, 120, w)
		}
		fmt.Println("referer:", referer)

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			ErrorMessage string
			Yourname     string
		}{
			Title:        "Signin | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
			CsrfToken:    ctoken,
			ErrorMessage: errMessage,
			Yourname:     "Sign In",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		//sTime := time.Now()
		r.ParseForm()
		var rurl string
		commonDataSet(r)

		funcsMap := map[string]interface{}{
			"validCSRF":    validCSRF,
			"validUsernme": validUsernme,
		}
		rmap := make(map[string]interface{})
		for key := range r.Form {
			rmap[key] = r.FormValue(key)
		}
		response := CheckMultipleConditionTrue(rmap, funcsMap)
		//fmt.Println("checkMultipleConditionTakes:", time.Since(sTime).Seconds(), "sec") //0.379260886 sec

		if response == "OKAY" {

			//sTime = time.Now()
			username := r.FormValue("email")
			txtpass := r.FormValue("passwd")
			ipAddress := cleanIp(mtool.ReadUserIP(r))
			userAgent := r.UserAgent()
			location := getLocationWithinSec(ipAddress)

			visitorSessionID, err := getCookie("visitor_session", r)
			if err != nil {
				visitorSessionID = visitorInfo(r, w)
				log.Println("visitorSession must required", err, visitorSessionID)
			}
			//fmt.Println("getLocationTakes:", time.Since(sTime).Seconds(), "sec") //0.405114986 sec

			//sTime = time.Now()
			sql := fmt.Sprintf("SELECT id,cid,account_id,access_name,username,passw,tfa_status,tfa_medium,tfa_setupkey FROM %s WHERE username='%s' AND status IN[1,6];", tableToBucket("login"), username)
			rows, err := lxql.GetRows(sql, database.DB)
			if err != nil {
				log.Println(err)
				return
			}
			//fmt.Println("loginQueryTakes:", time.Since(sTime).Seconds(), "sec") //0.365913886 sec
			if len(rows) > 0 {

				//sTime = time.Now()
				hashpass := rows[0]["passw"].(string) //
				if mtool.HashCompare(txtpass, hashpass) {

					//fmt.Println("hashCompareTakes:", time.Since(sTime).Seconds(), "sec") //1.601463995 sec
					//sTime = time.Now()
					loginId := rows[0]["id"].(string)
					accessName := rows[0]["access_name"].(string)
					accountId := rows[0]["access_name"].(string)
					token, err := vAuthToken(loginId, accountId, username, accessName, ipAddress) //takes 0.8 seconds to process
					logError("vAuthToken", err)
					//fmt.Println("vAuthTokenTakes:", time.Since(sTime).Seconds(), "sec")

					row := rows[0]
					tfaStatus := row["tfa_status"].(string)
					tfaMedium := row["tfa_medium"].(string)

					//TODO
					if tfaMedium == "email" {
						//generate code and send via email
						fmt.Println("email code")
					}
					if tfaStatus == "1" {

						//errNo = 0
						rurl = "/tfauth"
						//errMsg = "redirecting to TFA"
						setCookie("tfa", username, 300, w)

					} else {
						//sTime = time.Now()
						delete(row, "passw") //tfa_status,tfa_medium,tfa_setupkey
						delete(row, "tfa_status")
						delete(row, "tfa_medium")
						delete(row, "tfa_setupkey")
						delete(row, "account_type")
						row["session_code"] = visitorSessionID
						arow, _ := loginToAccountRow(loginId)
						for key, val := range arow {
							row[key] = val
						}
						row["ip"] = ipAddress
						jwtstr, err := utility.JWTEncode(row, utility.JWTSECRET) //takes 3 seconds to process
						if err != nil {
							log.Println(err)
							return
						}
						//fmt.Println("JWTEncodeTokenTakes:", time.Since(sTime).Seconds(), "sec")
						//sTime = time.Now()
						setCookie("login_session", jwtstr, 86400*30, w)
						setCookie("token", token, 86400*30, w) //30 days

						sql := fmt.Sprintf("UPDATE %s SET last_login='%s' WHERE id='%s';", tableToBucket("login"), mtool.TimeNow(), loginId)
						err = lxql.RawSQL(sql, database.DB)
						logError("updateLastLogin", err)

						var city, country string
						slc := strings.Split(location, ",")
						if len(slc) == 2 {
							city = slc[0]
							country = slc[1]
						}
						_, err = addLoginSession(loginId, visitorSessionID, ipAddress, city, country, userAgent)
						logError("addLoginSession", err)

						rurl = "/dashboard"
						redirect, _ := getCookie("redirect", r)
						if redirect == "/checkout" {
							rurl = redirect
							delCookie("redirect", r, w)
						}

						//fmt.Println("addUpdateLoginSessionTakes:", time.Since(sTime).Seconds(), "sec")
						//send email alert for login
						//sTime = time.Now()
						//_, _, browser := userAgntDetails(userAgent)
						//loginNotificationEmail(username, ipAddress, browser) //takes 6 seconds to process
						//fmt.Println("loginEmailTakes:", time.Since(sTime).Seconds(), "sec")
					}

				} else {
					rurl = "/signin?error=Invalid username or password"
				}
			}
			if len(rows) == 0 {
				rurl = "/signin?error=invalid username or password"
			}

		} else {
			rurl = "/signin?error=invalid username or password."
			log.Println(response)
		}

		http.Redirect(w, r, rurl, http.StatusSeeOther)
		return
	}
}

func resetpass(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		// smap, err := getSessionInfo(r)
		// if err == nil {
		// 	if access, isOk := smap["access_name"].(string); isOk {
		// 		vacc := []string{"superadmin", "admin", "client"}
		// 		if mtool.ArrayValueExist(vacc, access) {
		// 			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		// 		}
		// 	}
		// }

		//sessionCode := visitorInfo(r, w) //
		//fmt.Println(sessionCode)

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/resetpass.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			Yourname     string
		}{
			Title:        "LxRoot | Reset password",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
			CsrfToken:    ctoken,
			Yourname:     "Sign In",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		r.ParseForm()
		var errNo int = 1
		var errMsg string
		commonDataSet(r)
		//fmt.Println("post", r.Form)

		funcsMap := map[string]interface{}{
			"validCSRF": validCSRF,
		}
		rmap := make(map[string]interface{})
		for key := range r.Form {
			rmap[key] = r.FormValue(key)
		}
		response := CheckMultipleConditionTrue(rmap, funcsMap)

		if response == "OKAY" {

			username := r.FormValue("email")
			ipAddress := cleanIp(mtool.ReadUserIP(r))
			userAgent := r.UserAgent()

			//location := getLocationWithinSec(ipAddress)
			//fmt.Println(username, ipAddress, location)
			count := lxql.CheckCount("login", fmt.Sprintf(`username="%s"`, username), database.DB)
			if count == 1 {

				//generate a reset password email
				errNo = 0
				errMsg = "OK"
				_, _, browser := userAgntDetails(userAgent)
				resetPassNotificationEmail(username, ipAddress, browser)

			} else if count == 0 {

				errNo = 0
				fmt.Println("DO NOTHING - INVALID REQUEST")
			}

		} else {
			errNo = 0
			errMsg = "Validation error!"
			log.Println(response)
		}

		var row = make(map[string]interface{})
		row["error"] = errNo
		row["message"] = errMsg
		bs, err := json.Marshal(row)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(bs))
	}
}

func resetPassForm(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		// smap, err := getSessionInfo(r)
		// if err == nil {
		// 	if access, isOk := smap["access_name"].(string); isOk {
		// 		vacc := []string{"superadmin", "admin", "client"}
		// 		if mtool.ArrayValueExist(vacc, access) {
		// 			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		// 		}
		// 	}
		// }

		//sessionCode := visitorInfo(r, w) //
		//fmt.Println(sessionCode)

		r.ParseForm()
		var isValid bool
		token := r.FormValue("token")

		row, _ := tokenInfo(token)
		email, _ := row["email"].(string)
		isValid = checkTokenCodeValid(row)

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/reset-pass-form.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)
		invalidMsg := "Sorry, your token has expired!" //This reset link is invalid.

		if isValid {
			setCookie("vid", fmt.Sprint(row["id"]), 1800, w)
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			ValidToken   bool
			ErroMsg      string
			Username     string
			Yourname     string
		}{
			Title:        "LxRoot | Reset password form",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
			CsrfToken:    ctoken,
			ValidToken:   isValid,
			ErroMsg:      invalidMsg,
			Username:     email,
			Yourname:     "Sign In",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		r.ParseForm()
		var errNo int = 1
		var errMsg string
		commonDataSet(r)
		vid, _ := getCookie("vid", r)
		r.Form.Set("vid", vid)

		funcsMap := map[string]interface{}{
			"validCSRF":           validCSRF,
			"resetPassValidation": resetPassValidation,
			"vcodeIdValidation":   vcodeIdValidation,
		}
		rmap := make(map[string]interface{})
		for key := range r.Form {
			rmap[key] = r.FormValue(key)
		}
		response := CheckMultipleConditionTrue(rmap, funcsMap)
		//fmt.Println("response:", response)

		if response == "OKAY" {

			username := lxql.FieldByValue("verification", "username", fmt.Sprintf("id=%q", vid), database.DB)
			pass1 := r.FormValue("pass1")

			//location := getLocationWithinSec(ipAddress)
			count := lxql.CheckCount("login", fmt.Sprintf(`username="%s" AND access_name='client'`, username), database.DB)
			//fmt.Println(username, pass1, count)
			if count == 1 {

				//generate a reset password email
				delCookie("vid", r, w)

				sql := fmt.Sprintf("UPDATE %s SET status=1, update_date=%q WHERE id=%q;", tableToBucket("verification"), mtool.TimeNow(), vid)
				err = lxql.RawSQL(sql, database.DB)
				logError("verificationUpdate", err)

				//fmt.Println("now update login table with pass1")
				loginId := usernameToLoginId(username)
				sql = fmt.Sprintf("UPDATE %s SET passw=%q WHERE id=%q;", tableToBucket("login"), mtool.HashBcrypt(pass1), loginId)
				lxql.RawSQL(sql, database.DB)

				ipAddress := cleanIp(mtool.ReadUserIP(r))
				logDetails := fmt.Sprintf("username %s", username)
				addActiviyLog(loginId, RESET_PASS_ACTIVITY, "login", "", logDetails, ipAddress)
				errNo = 0
				errMsg = "Password reset successful."

			} else if count == 0 {
				errNo = 4
				errMsg = "Sorry, unauthorized user"
			}

		} else {
			errNo = 1
			errMsg = response
		}

		var row = make(map[string]interface{})
		row["error"] = errNo
		row["message"] = errMsg
		bs, err := json.Marshal(row)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(bs))
	}
}

func dashboard(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		smap, err := getSessionInfo(r)
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/footer2.gohtml",
			"wpages/dashboard.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		docNumber, _ := getCookie("docid", r)
		count := docToCartCount(docNumber)

		accountId, _ := smap["account_id"].(string)
		loginId, _ := smap["id"].(string)

		tickets := totalActiveTicketByUser(loginId)
		//fmt.Println(tickets, loginId)

		base := GetBaseURL(r)
		data := struct {
			Title          string
			Base           string
			BodyClass      string
			MainDivClass   string
			CartCount      int
			AccessName     string
			SessionMap     map[string]interface{}
			TotalOrders    int
			TotalInvoices  int
			UnpaidInvoices int
			ActiveTickets  int
		}{
			Title:          "LxRoot Dashboard",
			Base:           base,
			BodyClass:      "",
			MainDivClass:   "main min-h-[calc(100vh-52px)] bg-slate-200",
			CartCount:      count,
			SessionMap:     smap,
			TotalOrders:    totalOrdersByAccount(accountId),
			TotalInvoices:  totalInvoicesByAccount(accountId),
			UnpaidInvoices: totalUnpaidInvoicesByAccount(accountId),
			ActiveTickets:  tickets,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func logout(w http.ResponseWriter, r *http.Request) {

	smap, err := getSessionInfo(r)
	if err != nil {
		log.Println("UNABLE_TO_LOGOUT::", err)
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	//fmt.Println("logout:", smap)
	//sessionCode := smap["session_code"].(string)
	loginID := smap["id"].(string)
	qs := `UPDATE %s SET status=0,logout_time="%s" WHERE login_id="%s" AND status=1;`
	sql := fmt.Sprintf(qs, tableToBucket("login_session"), mtool.TimeNow(), loginID)
	_, err = database.DB.Exec(sql)
	logError("logout", err)
	if err == nil {
		//remove browser cookie
		token, err := getCookie("token", r)
		logError("logoutToken", err)

		sql = fmt.Sprintf("UPDATE %s SET status=0 WHERE token='%s';", tableToBucket("authc"), token)
		database.DB.Exec(sql)
		delCookie("token", r, w)
		delCookie("login_session", r, w)
		http.Redirect(w, r, "/signin", http.StatusSeeOther) //***
		return
	}
}

func profile(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		smap, err := getSessionInfo(r)
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/profile_left.gohtml",
			"templates/footer2.gohtml",
			"wpages/profile.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		//fmt.Println(smap)
		accountId, _ := smap["account_id"].(string)
		row := profileInfo(accountId)
		clientSince, _ := row["client_since"].(string)
		lastLoginDate, _ := row["last_login"].(string)

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			SessionMap   map[string]interface{}
			LastLogin    string
			ClientSince  string
			ProfileInfo  map[string]interface{}
			AddressList  []map[string]interface{}
		}{
			Title:        "LxRoot Profile",
			Base:         base,
			BodyClass:    "bg-slate-200",
			MainDivClass: "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:    ctoken,
			SessionMap:   smap,
			LastLogin:    lastLoginDate,
			ClientSince:  clientSince,
			ProfileInfo:  row,
			AddressList:  addressList(accountId),
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func security(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		smap, err := getSessionInfo(r)
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/profile_left.gohtml",
			"templates/footer2.gohtml",
			"wpages/security.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		loginId, _ := smap["id"].(string)
		lastLoginDate, clientSince := profileLastLogin(loginId)

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			SessionMap   map[string]interface{}
			LastLogin    string
			ClientSince  string
		}{
			Title:        "LxRoot Security",
			Base:         base,
			BodyClass:    "bg-slate-200",
			MainDivClass: "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:    ctoken,
			SessionMap:   smap,
			LastLogin:    lastLoginDate,
			ClientSince:  clientSince,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func tickets(w http.ResponseWriter, r *http.Request) {

	smap, err := getSessionInfo(r)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/footer2.gohtml",
			"wpages/tickets.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		//accountId, _ := smap["account_id"].(string)
		var ticketStatus string = "open"
		loginId, _ := smap["id"].(string)
		status := r.FormValue("status")
		if status == "" {
			status = "1"
		}
		if status == "0" {
			ticketStatus = "closed"
		}

		//closed tickets
		sql := fmt.Sprintf("SELECT id,department,subject,reference,ticket_status,ip_address,create_date,update_date FROM %s WHERE login_id=%q AND status=%s;", tableToBucket("ticket"), loginId, status)
		rows, err := lxql.GetRows(sql, database.DB)
		if err != nil {
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			SessionMap   map[string]interface{}
			Rows         []map[string]interface{}
			Count        int
			TicketStatus string
		}{
			Title:        "LxRoot Closed Tickets",
			Base:         base,
			BodyClass:    "bg-slate-200",
			MainDivClass: "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:    ctoken,
			SessionMap:   smap,
			Rows:         rows,
			Count:        len(rows),
			TicketStatus: ticketStatus,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		todo := parseMultipartTodo(r)
		var rmessage string = "Invalid request"
		if strings.ToUpper(todo) == "TICKET_REPLY" { //ticekt_reply

			loginId, _ := smap["id"].(string)
			ticketId := r.FormValue("tid")
			message := r.FormValue("message")
			ipAddress := cleanIp(r.RemoteAddr)
			addTicketResponse(ticketId, message, loginId, ipAddress)

			reference := xidToNumber(ticketId)
			rmessage = fmt.Sprintf("Ticket #%d response updated", reference)

			for _, mfh := range r.MultipartForm.File {
				for _, fh := range mfh {

					counter := xid.New().Counter()
					fileName := fmt.Sprintf("%d%s", counter, filepath.Ext(fh.Filename))
					fileAbsPath := filepath.Join("data", "ticket", fmt.Sprint(reference), fileName)
					err := saveFile(fh, fileAbsPath)
					logError("saveFile", err)
					_, err = addFileStore("ticket", ticketId, "", fileAbsPath, "")
					logError("addFileStore", err)
				}
			}
			fmt.Fprintln(w, rmessage)

		} else if strings.ToUpper(todo) == "CLOSE" {

			ticketId := r.FormValue("tid")
			//fmt.Println(ticketId, r.Referer()) //ticket?status=0
			sql := fmt.Sprintf("UPDATE %s SET ticket_status='closed',update_date=%q,status=0 WHERE id=%q;", tableToBucket("ticket"), mtool.TimeNow(), ticketId)
			lxql.RawSQL(sql, database.DB)
			http.Redirect(w, r, "/ticket?status=0", http.StatusSeeOther)
		}

	}
}

func ticketDetails(w http.ResponseWriter, r *http.Request) {

	smap, err := getSessionInfo(r)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {

		ticketId := chi.URLParam(r, "tid")
		//fmt.Println("ticketId:", ticketId)

		tmplt, err := template.New("base.gohtml").Funcs(FuncMap).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/footer2.gohtml",
			"wpages/ticket_details.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		sql := fmt.Sprintf("SELECT department,subject,message,reference,login_id,ticket_status,ip_address,create_date,update_date FROM %s WHERE id=%q;", tableToBucket("ticket"), ticketId)
		row, err := singleRow(sql)
		if err != nil {
			return
		}
		ticketOwner, _ := row["login_id"].(string)
		message := row["message"]
		ipAddress := row["ip_address"]
		createDate := row["create_date"]
		ticketStatus := row["ticket_status"].(string)
		//info, _ := loginToAccountRow(fmt.Sprint(ticketOwner))
		//accountName := info["account_name"]
		loginId, _ := smap["id"].(string)
		accountName, _ := smap["account_name"].(string)

		qs := `SELECT r.respond_by,r.message,r.ip_address,r.create_date,a.account_name FROM lxroot._default.ticket_response r LEFT JOIN lxroot._default.login l ON r.respond_by=l.id LEFT JOIN lxroot._default.account a ON a.id=l.account_id WHERE r.ticket_id=%q;`
		sql = fmt.Sprintf(qs, ticketId)
		rows, err := lxql.GetRows(sql, database.DB)
		if err != nil {
			log.Println(err, sql)
			return
		}

		var nrows = make([]map[string]interface{}, 0)
		if ticketOwner == loginId {
			trow := make(map[string]interface{})
			trow["respond_by"] = ticketOwner
			trow["message"] = message
			trow["ip_address"] = ipAddress
			trow["create_date"] = createDate
			trow["account_name"] = accountName
			nrows = append(nrows, trow)
		}
		nrows = append(nrows, rows...)
		//ptime, _ := time.Parse(DATE_TIME_FORMAT, "2024-06-12 07:47:00")
		//datetime = ptime.Format(outputFormat)
		//min := time.Since(ptime).Minutes()
		//hour := time.Since(ptime).Hours()

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			TicketId     string
			TicketStatus string
			SessionMap   map[string]interface{}
			TicketInfo   map[string]interface{}
			Responses    []map[string]interface{}
		}{
			Title:        fmt.Sprintf("LxRoot Ticket-%s", ticketId),
			Base:         base,
			BodyClass:    "bg-slate-200",
			MainDivClass: "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:    ctoken,
			TicketId:     ticketId,
			TicketStatus: ticketStatus,
			SessionMap:   smap,
			TicketInfo:   row,
			Responses:    nrows,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func orders(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		smap, err := getSessionInfo(r)
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}

		tmplt, err := template.New("base.gohtml").Funcs(FuncMap).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/footer2.gohtml",
			"wpages/orders.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		//purl, err := url.Parse(r.RequestURI)
		//fmt.Println(err, purl.Path)
		accountId, _ := smap["account_id"].(string)
		//loginId, _ := smap["id"].(string)

		rows, err := myOrders(accountId)
		logError("myOrders", err)
		//fmt.Println(rows)

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			SessionMap   map[string]interface{}
			Rows         []map[string]interface{}
		}{
			Title:        "LxRoot Orders",
			Base:         base,
			BodyClass:    "bg-slate-200",
			MainDivClass: "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:    ctoken,
			SessionMap:   smap,
			Rows:         rows,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func orderDetails(w http.ResponseWriter, r *http.Request) {

	smap, err := getSessionInfo(r)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(FuncMap).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/footer2.gohtml",
			"wpages/order_details.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}
		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		docId := chi.URLParam(r, "oid") //orderId
		accountId, _ := smap["account_id"].(string)
		//loginId, _ := smap["id"].(string)

		qs := `SELECT d.id, d.doc_status,d.doc_number,d.posting_date,d.receipt_url,d.payment_status,d.doc_name,d.doc_description,d.doc_ref,d.create_date,d.total_payable,t.item_info,t.price 
				FROM lxroot._default.doc_keeper d 
				LEFT JOIN  lxroot._default.transaction_record t ON d.doc_number=t.doc_number
				WHERE d.id="%s" AND d.account_id="%s";`
		sql := fmt.Sprintf(qs, docId, accountId)
		row, err := singleRow(cleanText(sql))
		logError("", err)
		//fmt.Println(row)

		docNumber, _ := row["doc_number"].(string)
		qs = `SELECT item_info,price,quantity,payable_amount FROM lxroot._default.transaction_record WHERE doc_number=%q;`
		sql = fmt.Sprintf(qs, docNumber)
		rows, err := lxql.GetRows(sql, database.DB)
		logError("itemRows", err)

		fmt.Println(row["receipt_url"])

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			SessionMap   map[string]interface{}
			Order        map[string]interface{}
			Items        []map[string]interface{}
		}{
			Title:        fmt.Sprintf("LxRoot #%s", toUpper(docNumber)),
			Base:         base,
			BodyClass:    "bg-slate-200",
			MainDivClass: "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:    ctoken,
			SessionMap:   smap,
			Order:        row,
			Items:        rows,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func invoices(w http.ResponseWriter, r *http.Request) {

	smap, err := getSessionInfo(r)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(FuncMap).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/footer2.gohtml",
			"wpages/invoices.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}
		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)
		//fmt.Println("smap:", smap)

		accountId, _ := smap["account_id"].(string)
		//loginId, _ := smap["id"].(string)

		qs := `SELECT d.id, d.doc_status,d.doc_number,d.posting_date,d.receipt_url,d.payment_status,d.doc_name,d.doc_description,d.doc_ref,d.create_date,d.total_payable,t.item_info,t.price 
				FROM lxroot._default.doc_keeper d 
				LEFT JOIN  lxroot._default.transaction_record t ON d.doc_number=t.doc_number
				WHERE d.account_id="%s" AND d.doc_type='sales' AND d.doc_status IN ['pending','complete'] AND d.status=1;`
		sql := fmt.Sprintf(qs, accountId)
		rows, err := lxql.GetRows(sql, database.DB)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			SessionMap   map[string]interface{}
			Rows         []map[string]interface{}
			Count        int
		}{
			Title:        "LxRoot Invoices",
			Base:         base,
			BodyClass:    "bg-slate-200",
			MainDivClass: "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:    ctoken,
			SessionMap:   smap,
			Rows:         rows,
			Count:        len(rows),
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func licenseKey(w http.ResponseWriter, r *http.Request) {

	smap, err := getSessionInfo(r)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/footer2.gohtml",
			"wpages/licensekey.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		accountId, _ := smap["account_id"].(string)
		//loginId, _ := smap["id"].(string)

		row, err := subscriptionDetailsByAccount(accountId)
		if err == nil {
			log.Println(err)
		}
		licenseKey, _ := row["license_key"].(string)
		paymentStatus, _ := row["payment_status"].(string)
		billing, _ := row["billing"].(string)
		price, _ := row["price"].(string)
		purchaseDate, _ := row["create_date"].(string)
		subscription_end, _ := row["subscription_end"].(string)

		dateFormat := "January 02, 2006"
		purchaseDate = mtool.DateTimeParser(purchaseDate, "2006-01-02 15:04:05", dateFormat)

		var licenseFound bool
		count := lxql.CheckCount("subscription", fmt.Sprintf("account_id=%q", accountId), database.DB)
		if count > 0 {
			licenseFound = true
		}

		base := GetBaseURL(r)
		data := struct {
			Title         string
			Base          string
			BodyClass     string
			MainDivClass  string
			CsrfToken     string
			SessionMap    map[string]interface{}
			LicenseKey    string
			PaymentStatus string
			Renews        string
			Price         string
			PurchaseDate  string
			ExpireDate    string
			LicenseFound  bool
		}{
			Title:         "LxRoot License Key",
			Base:          base,
			BodyClass:     "bg-slate-200",
			MainDivClass:  "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:     ctoken,
			SessionMap:    smap,
			LicenseKey:    licenseKey,
			PaymentStatus: paymentStatus,
			Renews:        billing,
			Price:         price,
			PurchaseDate:  purchaseDate,
			ExpireDate:    toTime(subscription_end).Format(dateFormat),
			LicenseFound:  licenseFound,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	}
}

func ticketNew(w http.ResponseWriter, r *http.Request) {

	smap, err := getSessionInfo(r)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header3.gohtml",
			"templates/footer2.gohtml",
			"wpages/ticketnew.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		ctoken := csrfToken()
		hashStr := hmacHash(ctoken, ENCDECPASS) //utility.ENCDECPASS
		setCookie("ctoken", hashStr, 1800, w)

		sql := fmt.Sprintf("SELECT id,name FROM %s WHERE status=1;", tableToBucket("department"))
		rows, err := lxql.GetRows(sql, database.DB)
		if err != nil {
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
			SessionMap   map[string]interface{}
			Rows         []map[string]interface{}
		}{
			Title:        "LxRoot New Ticket",
			Base:         base,
			BodyClass:    "bg-slate-200",
			MainDivClass: "main min-h-[calc(100vh-52px)] bg-slate-200",
			CsrfToken:    ctoken,
			SessionMap:   smap,
			Rows:         rows,
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {

		parseMultipartTodo(r)
		var message string = "OK"
		tokenPullNSet(r)

		funcsMap := map[string]interface{}{
			"validCSRF": validCSRF,
		}
		rmap := make(map[string]interface{})
		for key := range r.Form {
			rmap[key] = r.FormValue(key)
		}
		response := CheckMultipleConditionTrue(rmap, funcsMap)
		message = response
		//fmt.Println(message)

		if response == "OKAY" {

			subject := r.FormValue("subject")
			department := r.FormValue("department")
			ticketMessage := r.FormValue("message")
			ipAddress := cleanIp(r.RemoteAddr) //cleanIp(mtool.ReadUserIP())

			loginId, _ := smap["id"].(string)
			ticketId, _ := addTicket(loginId, department, subject, ticketMessage, "", ipAddress)
			reference := fmt.Sprint(xidToNumber(ticketId))

			for _, mfh := range r.MultipartForm.File {
				for _, fh := range mfh {

					counter := xid.New().Counter()
					fileName := fmt.Sprintf("%d%s", counter, filepath.Ext(fh.Filename))
					fileAbsPath := filepath.Join("data", "ticket", reference, fileName)
					err := saveFile(fh, fileAbsPath)
					logError("saveFile", err)
					_, err = addFileStore("ticket", ticketId, "", fileAbsPath, "")
					logError("addFileStore", err)
				}
			}
			message = fmt.Sprintf("Ticket #%s has been successfully created.", reference)
		}
		fmt.Fprintln(w, message)
	}
}

func invoice(w http.ResponseWriter, r *http.Request) {

	smap, err := getSessionInfo(r)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}
	inv := chi.URLParam(r, "inv")
	fmt.Println(inv, smap)

	http.ServeFile(w, r, "data/invoice/Receipt-2082-6216.pdf")
}
