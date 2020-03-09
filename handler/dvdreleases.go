package handler

import (
	"fmt"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/fr05t1k/traktbot/traktapi"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

func CreateDvdReleasesHandler(bot *tb.Bot, trakt *traktapi.Client, tmdbClient *tmdb.Client) func(m *tb.Message) {
	return func(m *tb.Message) {
		releases, err := trakt.GetDvdReleases()
		if err != nil {
			log.Fatal(err)
			return
		}

		for i := range releases {

			bot.Send(
				m.Sender,
				buildReleasesMessage(tmdbClient, releases[i]),
				tb.ModeMarkdown,
			)
		}
	}
}

func buildReleasesMessage(tmdbClient *tmdb.Client, release traktapi.Release) interface{} {
	image := ""
	movie, err := tmdbClient.GetMovieDetails(release.Movie.Ids.Tmdb, map[string]string{"append_to_response": "images,videos"})
	if err != nil {
		log.Println(err)
		return nil
	}

	if len(movie.Images.Posters) > 0 {
		image = tmdb.GetImageURL(movie.Images.Posters[0].FilePath, "w500")
	}

	trailer := ""
	if len(movie.Videos.Results) > 0 {
		for i := range movie.Videos.Results {
			if movie.Videos.Results[i].Type == "Trailer" {
				trailer = tmdb.GetVideoURL(movie.Videos.Results[i].ID)
			}
		}

	}

	body := fmt.Sprintf(
		"%s\n%s (%d)\n%s",
		release.Released,
		release.Movie.Title,
		release.Movie.Year,
		movie.Overview,
	)
	if trailer != "" {
		body += fmt.Sprintf(" [Trailer](%s)", trailer)
	}

	if image != "" {
		return &tb.Photo{
			File:      tb.FromURL(image),
			Caption:   body,
			ParseMode: tb.ModeMarkdown,
		}
	}

	return body
}
