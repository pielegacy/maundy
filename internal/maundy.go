package internal

import (
	"fmt"
	"maundy/internal/apple"
	"maundy/internal/spotify"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

const EnvVarBrowserPath = "MAUNDY_BROWSER_PATH"
const EnvVarApplePlaylistUrl = "MAUNDY_APPLE_PLAYLIST_URL"
const EnvVarSpotifyClientId = "MAUNDY_SPOTIFY_CLIENT_ID"
const EnvVarSpotifyClientSecret = "MAUNDY_SPOTIFY_CLIENT_SECRET"
const EnvVarSpotifyRefreshToken = "MAUNDY_SPOTIFY_REFRESH_TOKEN"
const EnvVarSpotifyPlaylistId = "MAUNDY_SPOTIFY_PLAYLIST_ID"

type SyncParams struct {
}

func Sync(syncParams SyncParams) {
	fmt.Println("-- Maundy --")
	godotenv.Load()

	browserPath := os.Getenv(EnvVarBrowserPath)
	playlistUrl := os.Getenv(EnvVarApplePlaylistUrl)

	url := launcher.New().Bin(browserPath).Set("--no-sandbox").MustLaunch()
	page := rod.New().ControlURL(url).MustConnect().MustPage(playlistUrl).MustWaitStable()

	fmt.Printf("Browser Path: %s\n", browserPath)
	fmt.Printf("[APPLE] Playlist URL: %s\n", playlistUrl)

	songs := make([]apple.Song, 0)
	for _, songElement := range apple.GetSongElements(page) {
		song := apple.ParseSongElement(songElement)
		fmt.Printf("[APPLE] Song loaded: '%s' by '%s'\n", song.Title, song.Artist)
		songs = append(songs, song)
	}

	context := spotify.GetContext(spotify.ContextRequest{
		ClientId:     os.Getenv(EnvVarSpotifyClientId),
		ClientSecret: os.Getenv(EnvVarSpotifyClientSecret),
		RefreshToken: os.Getenv(EnvVarSpotifyRefreshToken),
	})
	if !context.LoggedIn {
		panic("Spotify login failed")
	}

	playlistId := os.Getenv(EnvVarSpotifyPlaylistId)
	fmt.Printf("[SPOTIFY] Target Playlist ID: %s\n", playlistId)

	uris := make([]string, 0)

	for _, song := range songs {
		track, err := context.SearchTrack(spotify.SearchTrackRequest{
			Title:  song.Title,
			Artist: song.Artist,
		})
		if err != nil {
			panic(fmt.Sprintf("[SPOTIFY] Failed to match track : '%s' by '%s'", song.Title, song.Artist))
		}
		uris = append(uris, track.Uri)
	}

	fmt.Printf("[SPOTIFY] Updating playlist %s\n", playlistId)
	err := context.UpdatePlaylist(spotify.UpdatePlaylistRequest{
		Id:   playlistId,
		Uris: uris,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Playlist synced")

}
