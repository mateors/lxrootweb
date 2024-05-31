package main

import (
	"encoding/json"
	"fmt"
	"log"
	"lxrootweb/database"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/mateors/mtool"
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

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Shop | LxRoot",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func product(w http.ResponseWriter, r *http.Request) {

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

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Product Details | LxRoot",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func checkout(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/checkout.gohtml", //
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
			Title:        "Checkout | LxRoot",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
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
			"validPermission": validCSRF,
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

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
			CsrfToken    string
		}{
			Title:        "Signup | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
			CsrfToken:    ctoken,
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
			"validPermission":  validCSRF,
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
				addLogin(accountId, accessId, accessName, username, passwd)
				code := uuid.NewV4().String()
				_, err = addVerification(username, "signup", code, "")
				logError("addVerification", err)

				location := "Dhaka, Bangladesh"
				verfyUrl := fmt.Sprintf("https://lxroot.com/verify?email=%s&token=%s", email, code)
				err = signupEmail(email, accountName, location, verfyUrl)
				logError("signupEmail", err)

				errNo = 0
				errMsg = "Congratulations! You have successfully registered with LxRoot. 🎉"
			}

		} else {
			errNo = 9
			errMsg = "Unknown error! please report"
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

		email := r.FormValue("email")
		token := r.FormValue("token")
		ainfo, err := usernameToAccounInfo(email)
		if err != nil {
			logError("usernameToAccounInfo", err)
			return
		}

		var message string
		err = verifySignup(email, token)
		if err != nil {
			message = fmt.Sprintf(`<h1 style="text-align:center;color:red;font-size:64px;">%s</h1>`, err.Error())

		} else {

			fmt.Println(ainfo, len(ainfo))
			accountId, isOk := ainfo["account_id"].(string)
			if !isOk {
				log.Println("unable to parse ainfo")
				return
			}

			sql := fmt.Sprintf("UPDATE %s SET status=1,update_date=%q,remarks='verified' WHERE id=%q;", tableToBucket("account"), mtool.TimeNow(), accountId)
			database.DB.Exec(sql)

			loginId := ainfo["login_id"].(string)
			sql = fmt.Sprintf("UPDATE %s SET status=1 WHERE id=%q;", tableToBucket("login"), loginId)
			database.DB.Exec(sql)

			sql = fmt.Sprintf("UPDATE %s SET status=1,update_date=%q WHERE id=%q;", tableToBucket("verification"), mtool.TimeNow(), loginId)
			database.DB.Exec(sql)

			ipAddress := cleanIp(mtool.ReadUserIP(r))
			logDetails := fmt.Sprintf("username %s verified", email)
			addActiviyLog(loginId, "UPDATE", "account", "", logDetails, ipAddress)
			message = fmt.Sprintf(`<h1 style="text-align:center;color:green;font-size:64px;">%s <a href="/signin">Sign in</a></h1>`, "Verified")
		}
		fmt.Fprintln(w, message)
	}
}

func signin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

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

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Signin | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}
