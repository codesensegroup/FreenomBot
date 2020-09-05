package main

import (
	"fmt"

	freenom "github.com/codesensegroup/FreenomBot/internal/freenom"
)

func main() {

	f := freenom.GetInstance().Login().RenewDomains()
	for _, d := range f.Domains {
		fmt.Println(d)
	}
}
