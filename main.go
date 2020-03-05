package main

import (
	"fmt"
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
	traktClient := traktapi.NewClient(
		os.Getenv("TRAKT_CLIENT_ID"),
		os.Getenv("TRAKT_CLIENT_SECRET"),
	)

	b.Handle("/login", func(m *tb.Message) {
		code, err := traktClient.DeviceCode()
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Error: %s", err))
			return
		}

		b.Send(
			m.Sender,
			fmt.Sprintf("Please visit %s and input the code **%s**", code.VerificationUrl, code.UserCode),
			tb.ModeMarkdown,
		)

		times := code.ExpiresIn / code.Interval
		for ; times >= 0; times-- {
			<-time.After(time.Duration(code.Interval))
			token, err := traktClient.GetDeviceToken(code.DeviceCode)
			if err != nil {
				log.Println(err)
				continue
			}

			traktClient.SaveDeviceToken(token)
			b.Send(
				m.Sender,
				"You are logged in!",
			)
			return
		}

		b.Send(
			m.Sender,
			"Authentication expired",
		)
	})

	b.Start()
}
