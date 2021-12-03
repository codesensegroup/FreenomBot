package freenom

import (
	"errors"
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

	"FreenomBot/config"

	"github.com/tidwall/gjson"
	"golang.org/x/net/publicsuffix"
)

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
		instance.client = &http.Client{Timeout: timeout * time.Second, Jar: instance.cookiejar}
		instance.ConfigData = config.GetData()
	})
	return instance
}

// Login on Freenom
func (f *Freenom) Login(user *User) error {
	_, err := f.sendRequest(
		"POST",
		loginURL,
		`{"headers":{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer": "`+refererURL+`",
		},}`,
		url.Values{
			"username": {user.UserName},
			"password": {user.Password},
		}.Encode(),
	)
	if err != nil {
		return err
	}

	u, _ := url.Parse(baseURL)
	for _, authcook := range f.cookiejar.Cookies(u) {
		if authKey == authcook.Name && authcook.Value == "" {
			return errors.New("AUTH error")
		}
		log.Println("log: cookie_id: ", authcook.Value)
	}

	return nil
}

func (f *Freenom) GetFreenomInfo(user *User) error {
	body, err := f.sendRequest(
		"GET",
		domainStatusURL,
		`{"headers":{
			"Referer": "`+refererURL+`"
		},}`,
		"",
	)
	if err != nil {
		return err
	}

	if !loginStatusREGEX.Match(body) {
		return errors.New("login state error no login")
	}

	user.Token = getParams(tokenREGEX, body)[0]["token"]
	domains := getParams(domainInfoREGEX, body)
	user.Domains = make([]Domain, len(domains))

	for i, d := range domains {
		tmp := d["days"]
		user.Domains[i].Days, _ = strconv.Atoi(tmp)
		user.Domains[i].ID = d["id"]
		user.Domains[i].DomainName = d["domain"]
	}

	return nil
}

//RenewDomains is renew domain name
func (f *Freenom) RenewDomains(user *User) error {
	for i := range user.Domains {
		if user.Domains[i].Days <= delayRegistionDelayTime {
			body, _ := f.sendRequest(
				"POST",
				renewDomainURL,
				`{"headers":{
					"Referer": "https://my.freenom.com/domains.php?a=renewdomain&domain=`+user.Domains[i].ID+`",
					"Content-Type": "application/x-www-form-urlencoded",
				},}`,
				url.Values{
					"token":     {user.Token},
					"renewalid": {user.Domains[i].ID},
					"renewalperiod[" + user.Domains[i].ID + "]": {"12M"},
					"paymentmethod": {"credit"},
				}.Encode(),
			)
			if checkRenew.Match(body) {
				user.Domains[i].RenewState = RenewYes
			} else {
				log.Fatalln("renew error")
				user.Domains[i].RenewState = RenewErr
			}
		} else {
			user.Domains[i].RenewState = RenewNo
		}
	}
	return nil
}

/**
 * Parses url with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func getParams(regEx *regexp.Regexp, body []byte) (paramsMaps map[int]map[string]string) {
	// log.Println(string(body))
	match := regEx.FindAllStringSubmatch(string(body), -1)
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
func (f *Freenom) sendRequest(method, furl, headers, datas string) ([]byte, error) {
RETRY:
	req, err := http.NewRequest(method, furl, strings.NewReader(datas))
	if err != nil {
		log.Fatal("Create http request error", err)
		time.Sleep(3 * time.Second)
		goto RETRY
	}

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
	return body, nil
}
