package line

import (
	"github.com/utahta/go-linenotify"
)

var token *string

// Init get Token
func Init(tk *string) {
	token = tk
}

// Send message
func Send(msg string) {
	c := linenotify.New()
	c.Notify(*token, msg, "", "", nil)
}
