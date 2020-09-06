package freenom

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"golang.org/x/net/publicsuffix"
)

const (
	version         = "v0.0.5"
	timeout         = 34
	baseURL         = "https://my.freenom.com"
	refererURL      = "https://my.freenom.com/clientarea.php"
	loginURL        = "https://my.freenom.com/dologin.php"
	domainStatusURL = "https://my.freenom.com/domains.php?a=renewals"
	renewDomainURL  = "https://my.freenom.com/domains.php?submitrenewals=true"
	authKey         = "WHMCSZH5eHTGhfvzP"
)

var (
	tokenREGEX       = regexp.MustCompile(`name="token"\svalue="(?P<token>[^"]+)"`)
	domainInfoREGEX  = regexp.MustCompile(`<tr><td>(?P<domain>[^<]+)<\/td><td>[^<]+<\/td><td>[^<]+<span class="[^"]+">(?P<days>\d+)[^&]+&domain=(?P<id>\d+)"`)
	loginStatusREGEX = regexp.MustCompile(`<li.*?Logout.*?<\/li>`)
	checkRenew       = regexp.MustCompile(`(?i)Order Confirmation`)
)

//Domain struct
type Domain struct {
	DomainName string
	Days       int
	ID         string
	RenewState int
}

//User data
type User struct {
	UserName   string
	PassWord   string
	CheckTimes int
	Domains    map[int]*Domain
}

// Freenom for opterate FreenomAPI
type Freenom struct {
	cookiejar *cookiejar.Jar
	client    *http.Client
	Users     map[int]*User
}

var instance *Freenom
var once sync.Once

var (
	renewNo  int = 0
	renewYes int = 1
	renewErr int = 3
)

// GetInstance is get  instance
func GetInstance() *Freenom {
	once.Do(func() {
		instance = &Freenom{}
		instance.Users = make(map[int]*User)
		instance.cookiejar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		instance.client = &http.Client{Timeout: timeout * time.Second, Jar: instance.cookiejar}
	})
	return instance
}

// InputAccount input user data
func (f *Freenom) InputAccount(uid int, UserName, PassWord string) *Freenom {
	f.Users[uid] = &User{
		UserName:   UserName,
		PassWord:   PassWord,
		CheckTimes: 0,
	}
	return f
}

// Login on Freenom
func (f *Freenom) Login(uid int) *Freenom {
	_ = sendRequest(
		"POST",
		loginURL,
		`{"headers":{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer": "`+refererURL+`",
		},}`,
		url.Values{
			"username": {f.Users[uid].UserName},
			"password": {f.Users[uid].PassWord},
		}.Encode(),
	)

	u, _ := url.Parse(baseURL)
	for _, authcook := range f.cookiejar.Cookies(u) {
		if authKey == authcook.Name && authcook.Value == "" {
			log.Println("AUTH error")
		}
		log.Println("log: cookie_id: ", authcook.Value)
	}
	return f
}

//RenewDomains is renew domain name
func (f *Freenom) RenewDomains(uid int) *Freenom {
	body := sendRequest(
		"GET",
		domainStatusURL,
		`{"headers":{
			"Referer": "`+refererURL+`"
		},}`,
		"",
	)

	if !loginStatusREGEX.Match(body) {
		log.Fatal("login state error no login")
	}

	f.Users[uid].CheckTimes++

	var token = getParams(tokenREGEX, string(body))[0]["token"]

	domains := getParams(domainInfoREGEX, string(body))
	f.Users[uid].Domains = make(map[int]*Domain)
	for i, d := range domains {
		tmp, _ := d["days"]
		f.Users[uid].Domains[i] = &Domain{}
		f.Users[uid].Domains[i].Days, _ = strconv.Atoi(tmp)
		f.Users[uid].Domains[i].ID, _ = d["id"]
		f.Users[uid].Domains[i].DomainName, _ = d["domain"]
		if f.Users[uid].Domains[i].Days <= 14 {
			body := sendRequest(
				"POST",
				renewDomainURL,
				`{"headers":{
					"Referer": "https://my.freenom.com/domains.php?a=renewdomain&domain=`+f.Users[uid].Domains[i].ID+`",
					"Content-Type": "application/x-www-form-urlencoded",
				},}`,
				url.Values{
					"token":     {token},
					"renewalid": {f.Users[uid].Domains[i].ID},
					"renewalperiod[" + f.Users[uid].Domains[i].ID + "]": {"12M"},
					"paymentmethod": {"credit"},
				}.Encode(),
			)
			if checkRenew.Match(body) {
				f.Users[uid].Domains[i].RenewState = renewYes
			} else {
				log.Fatalln("renew error")
				f.Users[uid].Domains[i].RenewState = renewErr
			}
		} else {
			f.Users[uid].Domains[i].RenewState = renewNo
		}
	}
	return f
}

/**
 * Parses url with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func getParams(regEx *regexp.Regexp, url string) (paramsMaps map[int]map[string]string) {
	match := regEx.FindAllStringSubmatch(url, -1)
	paramsMaps = map[int]map[string]string{}

	for j := 0; j < len(match); j++ {
		paramsMaps[j] = make(map[string]string)
		for i, name := range regEx.SubexpNames() {
			if i > 0 && i <= len(match[j]) {
				paramsMaps[j][name] = match[j][i]
			}
		}
	}
	return
}

/**
 * getRequest just all in one
 */
func sendRequest(method, furl, headers, datas string) []byte {
	req, err := http.NewRequest(method, furl, strings.NewReader(datas))
	if err != nil {
		log.Fatal("Create http request error", err)
	}
	f := GetInstance()
	if headers != "" {
		headerObj := gjson.Get(headers, "headers")
		headerObj.ForEach(func(key, value gjson.Result) bool {
			req.Header.Add(key.String(), value.String())
			return true
		})
	}
	resp, _ := f.client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}
