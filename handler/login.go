package handler

import (
	"fmt"
	"github.com/fr05t1k/traktbot/traktapi"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

func CreateLoginHandler(bot *tb.Bot, trakt *traktapi.Client) func(m *tb.Message) {
	return func(m *tb.Message) {
		code, err := trakt.DeviceCode()
		if err != nil {
			bot.Send(m.Sender, fmt.Sprintf("Error: %s", err))
			return
		}

		bot.Send(
			m.Sender,
			fmt.Sprintf("Please visit %s and input the code **%s**", code.VerificationUrl, code.UserCode),
			tb.ModeMarkdown,
		)

		times := code.ExpiresIn / code.Interval
		for ; times >= 0; times-- {
			<-time.After(time.Duration(code.Interval))
			token, err := trakt.GetDeviceToken(code.DeviceCode)
			if err != nil {
				log.Println(err)
				continue
			}

			trakt.SaveDeviceToken(token)
			bot.Send(
				m.Sender,
				"You are logged in!",
			)
			return
		}

		bot.Send(
			m.Sender,
			"Authentication expired",
		)
	}
}
