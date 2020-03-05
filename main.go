package main

import (
	"github.com/fr05t1k/traktbot/handler"
	"github.com/fr05t1k/traktbot/traktapi"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"time"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TRAKT_BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	trakt := traktapi.NewClient(
		os.Getenv("TRAKT_CLIENT_ID"),
		os.Getenv("TRAKT_CLIENT_SECRET"),
	)

	b.Handle("/login", handler.CreateLoginHandler(b, trakt))

	b.Start()
}
