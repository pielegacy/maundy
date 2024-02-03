# Maundy

This is a small project I've used to learn Go, it probably still needs a bit of cleaning up and refactoring but I'm happy with the final result. Maundy exists so I can keep my supreme Apple Music playlists synced to Spotify for the peasants to enjoy ([thus the name](https://en.wikipedia.org/wiki/Royal_Maundy)).

## Why?

Just for fun, I was very keen to learn Go and this seemed like a good way to get acquanted with it. 

The project included the usage of external dependencies, HTTP calls, JSON serialization, some really basic logic for copying playlists over - it was nothing insanely complicated however the project overall was a good way to dip my toe into a new language. 

## Key Features

* Utilizes [rod](https://pkg.go.dev/github.com/go-rod/rod@v0.114.6) to pull Apple Music playlist data without requiring an Apple Developer Subscription
* Docker image includes Chrome baked in for optimal execution on cloud platforms
* Could probably sync more than 10 songs across (I'm too scared to try)

## Working Example

![PlayList Comparison](/docs/comparison.png)

Source Playlist: https://music.apple.com/au/playlist/the-current-ten/pl.u-DdANXBqIaJYNyqK

Target Playlist: https://open.spotify.com/playlist/28kUQ0lpoCgR3nZQPU1yOU
