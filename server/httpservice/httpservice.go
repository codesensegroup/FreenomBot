package httpservice

import (
	"log"
	"net/http"
	"regexp"
	"text/template"

	"github.com/codesensegroup/FreenomBot/internal/freenom"
)

type PageData struct {
	Users []User
}
type User struct {
	UserName string
	Domains  []Domain
}
type Domain struct {
	DomainName string
	Days       int
	ID         string
	RenewState int
}

var validPath = regexp.MustCompile("^/$")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *freenom.Freenom) {
	var pdata = &PageData{}
	pdata.Users = make([]User, len(data.Users))
	for i, user := range data.Users {
		pdata.Users[i].UserName = user.UserName
		pdata.Users[i].Domains = make([]Domain, len(user.Domains))
		for ii, domain := range user.Domains {
			pdata.Users[i].Domains[ii].DomainName = domain.DomainName
			pdata.Users[i].Domains[ii].Days = domain.Days
			pdata.Users[i].Domains[ii].ID = domain.ID
			pdata.Users[i].Domains[ii].RenewState = domain.RenewState
		}
	}

	t, err := template.ParseFiles("./resources/html/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, pdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Run server
func Run(data *freenom.Freenom) {
	http.HandleFunc("/", makeHandler(func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "status", data)
	}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
