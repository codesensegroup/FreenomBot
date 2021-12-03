package freenom

import "regexp"

const (
	Version                 = "v0.0.5"
	timeout                 = 34
	delayRegistionDelayTime = 14
	baseURL                 = "https://my.freenom.com"
	refererURL              = "https://my.freenom.com/clientarea.php"
	loginURL                = "https://my.freenom.com/dologin.php"
	domainStatusURL         = "https://my.freenom.com/domains.php?a=renewals"
	renewDomainURL          = "https://my.freenom.com/domains.php?submitrenewals=true"
	authKey                 = "WHMCSZH5eHTGhfvzP"
)

var (
	tokenREGEX       = regexp.MustCompile(`(?s)name="token"\svalue="(?P<token>[^"]+)"`)
	domainInfoREGEX  = regexp.MustCompile(`<tr><td>(?P<domain>[^<]+)<\/td><td>[^<]+<\/td><td>[^<]+<span class="[^"]+">(?P<days>\d+)[^&]+&domain=(?P<id>\d+)"`)
	loginStatusREGEX = regexp.MustCompile(`<li.*?Logout.*?<\/li>`)
	checkRenew       = regexp.MustCompile(`(?i)Order Confirmation`)
)

const (
	RenewNo  int = 0
	RenewYes int = 1
	RenewErr int = 3
)
