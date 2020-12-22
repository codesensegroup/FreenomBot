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

	"github.com/codesensegroup/FreenomBot/internal/config"
	"github.com/tidwall/gjson"
	"golang.org/x/net/publicsuffix"
)

//Domain struct
type Domain struct {
	DomainName string
	Days       int
	ID         string
	CheckTimes int
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
	cookiejar  *cookiejar.Jar
	client     *http.Client
	ConfigData *config.Config
}

var instance *Freenom
var once sync.Once

// GetInstance is the singleton
func GetInstance() *Freenom {
	once.Do(func() {
		instance = &Freenom{}
		instance.cookiejar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		instance.client = &http.Client{Timeout: config.Timeout * time.Second, Jar: instance.cookiejar}
		instance.ConfigData = config.GetData()
	})
	return instance
}

// Login on Freenom
func (f *Freenom) Login(uid int) *Freenom {
	_ = sendRequest(
		"POST",
		config.LoginURL,
		`{"headers":{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer": "`+config.RefererURL+`",
		},}`,
		url.Values{
			"username": {f.ConfigData.Accounts[uid].Username},
			"password": {f.ConfigData.Accounts[uid].Password},
		}.Encode(),
	)

	u, _ := url.Parse(config.BaseURL)
	for _, authcook := range f.cookiejar.Cookies(u) {
		if config.AuthKey == authcook.Name && authcook.Value == "" {
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
		config.DomainStatusURL,
		`{"headers":{
			"Referer": "`+config.RefererURL+`"
		},}`,
		"",
	)

	if !config.LoginStatusREGEX.Match(body) {
		log.Fatal("login state error no login")
	}

	//f.Users[uid].CheckTimes++

	token := getParams(config.TokenREGEX, string(body))[0]["token"]

	domains := getParams(config.DomainInfoREGEX, string(body))

	config.CreateDomains(uid, len(domains))
	for i, d := range domains {
		tmp, _ := d["days"]
		f.ConfigData.Accounts[uid].Domains[i].Days, _ = strconv.Atoi(tmp)
		f.ConfigData.Accounts[uid].Domains[i].ID, _ = d["id"]
		f.ConfigData.Accounts[uid].Domains[i].DomainName, _ = d["domain"]
		if f.ConfigData.Accounts[uid].Domains[i].Days <= config.DelayRegistionDelayTime {
			body := sendRequest(
				"POST",
				config.RenewDomainURL,
				`{"headers":{
					"Referer": "https://my.freenom.com/domains.php?a=renewdomain&domain=`+f.ConfigData.Accounts[uid].Domains[i].ID+`",
					"Content-Type": "application/x-www-form-urlencoded",
				},}`,
				url.Values{
					"token":     {token},
					"renewalid": {f.ConfigData.Accounts[uid].Domains[i].ID},
					"renewalperiod[" + f.ConfigData.Accounts[uid].Domains[i].ID + "]": {"12M"},
					"paymentmethod": {"credit"},
				}.Encode(),
			)
			if config.CheckRenew.Match(body) {
				f.ConfigData.Accounts[uid].Domains[i].RenewState = config.RenewYes
			} else {
				log.Fatalln("renew error")
				f.ConfigData.Accounts[uid].Domains[i].RenewState = config.RenewErr
			}
		} else {
			f.ConfigData.Accounts[uid].Domains[i].RenewState = config.RenewNo
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
 * sendRequest just all in one
 */
func sendRequest(method, furl, headers, datas string) []byte {
RETRY:
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
	resp, err := f.client.Do(req)
	if err != nil {
		log.Println("http response error: ", err)
		time.Sleep(3 * time.Second)
		goto RETRY
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
