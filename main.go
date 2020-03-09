package main

import (
	"github.com/cyruzin/golang-tmdb"
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
	}
	trakt := traktapi.NewClient(
		os.Getenv("TRAKT_CLIENT_ID"),
		os.Getenv("TRAKT_CLIENT_SECRET"),
	)

	tmdbClient, err := tmdb.Init(os.Getenv("TRAKT_TMDB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/login", handler.CreateLoginHandler(b, trakt))
	b.Handle("/dvd-releases", handler.CreateDvdReleasesHandler(b, trakt, tmdbClient))

	b.Start()
}
