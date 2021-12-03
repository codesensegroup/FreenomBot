package freenom

type PageData struct {
	Users []User
}

// User is translate freenom map data
type User struct {
	UserName   string
	Password   string
	CheckTimes int
	Token      string
	Domains    []Domain
}

// Domain is translate freenom map data
type Domain struct {
	DomainName string
	Days       int
	ID         string
	CheckTimes int
	RenewState int
}
