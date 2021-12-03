package line

import (
	"context"
	"log"

	"github.com/utahta/go-linenotify"
)

var c *linenotify.Client

func Init() {
	c = linenotify.NewClient()
}

// Send message
func Send(token *string, msg string) {
	if rep, err := c.NotifyMessage(context.Background(), *token, msg); err != nil {
		log.Fatalln(err)
	} else {
		log.Println(rep.Message)
	}
}
