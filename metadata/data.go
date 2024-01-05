package metadata

import (
	"bytes"
	"fmt"
	"time"
)

//Licensed under Creative Commons CC0.
//
//To the extent possible under law, the author(s) have dedicated all copyright and related and
//neighboring rights to this software to the public domain worldwide.
//This software is distributed without any warranty.
//You should have received a copy of the CC0 Public Domain Dedication along with this software.
//If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

var TestSiteAddress = "http://127.0.0.1:7777"
var TestImplementation = "http://127.0.0.1:7777/1223c99f-70fe-40be-abe3-bf1c6ba1bdb6.txt"
var TestTitle = "Digital Marketing"
var TestPaymentUrl = ""
var TestSite = ""

var ProxySite = ""
var Placeholder = "<!-- Placeholder for digital advertisement -->"
var SiteName = TestSiteAddress
var PaymentUrl = ""
var SiteTitle = ""
var Terms = ""
var Contact = ""
var AdBlocker = "<div style=\"text-align: center\"><p>Block ads for your convenience. <a href=\"%s\">🐞(hop)</a><!--%s--> </p></div>\n"
var AdDescription = "🐞 Advertisement Technology"

var DefaultAdTime = 24 * time.Hour
var DefaultPurchaseTime = 10 * time.Minute

func GetDefaultImplementation() string {
	buf := bytes.Buffer{}
	if TestPaymentUrl != "" {
		buf.WriteString(fmt.Sprintf("Set the payment url to %s address.", TestPaymentUrl))
		buf.WriteByte('\n')
	}
	buf.WriteString(fmt.Sprintf("Set the title to %s text.", TestTitle))
	buf.WriteByte('\n')
	if TestSite != "" {
		buf.WriteString(fmt.Sprintf("Proxy the %s site.", TestSite))
		buf.WriteByte('\n')
	}

	return buf.String()
}
