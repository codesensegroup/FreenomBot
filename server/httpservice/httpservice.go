package httpservice

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"FreenomBot/common/freenom"
	getConfig "FreenomBot/config"
)

var validPath = regexp.MustCompile("^/$")

func makeHandler(fn func(http.ResponseWriter, *http.Request), config *getConfig.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		basicAuthPrefix := "Basic "
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):],
			)
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && bytes.Equal(pair[0], []byte(config.System.Account)) && bytes.Equal(pair[1], []byte(config.System.Password)) {
					fn(w, r)
				}
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *freenom.PageData) {
	//var pdata = getPageData(&PageData{}, data)

	t, err := template.ParseFiles("./resources/html/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Run server
func Run(data *freenom.PageData, config *getConfig.Config) {
	http.HandleFunc("/", makeHandler(func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "status", data)
	}, config))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
