package apple

import "github.com/go-rod/rod"

type Song struct {
	Title  string
	Artist string
}

func GetSongElements(page *rod.Page) rod.Elements {
	results, err := page.Elements(".songs-list-row__song-wrapper")
	if err != nil {
		panic("Failed to find song wrappers on provided page")
	}
	return results
}

func ParseSongElement(el *rod.Element) Song {
	songNameElement, songNameElementErr := el.Elements(".songs-list-row__song-name")
	if songNameElementErr != nil {
		panic("Failed to find song name element")
	}
	songName, songNameErr := songNameElement.First().Text()
	if songNameErr != nil {
		panic("Failed to parse song name")
	}

	songArtistElement, _ := el.Elements(".click-action")
	songArtist, _ := songArtistElement.First().Text()

	return Song{
		Title:  songName,
		Artist: songArtist,
	}
}
