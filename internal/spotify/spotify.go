package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const baseUrl string = "https://api.spotify.com/v1/"

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ContextRequest struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
}

type Context struct {
	LoggedIn    bool
	AccessToken string
	Client      http.Client
}

type Playlist struct {
	Name string `json:"name"`
}

type SearchTrackRequest struct {
	Title  string
	Artist string
}

type SearchTrackResponse struct {
	Tracks struct {
		Items []SearchTrackResponseItem
	} `json:"tracks"`
}

type SearchTrackResponseItem struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type UpdatePlaylistRequest struct {
	Id   string   `json:"-"`
	Uris []string `json:"uris"`
}

func GetContext(contextReq ContextRequest) *Context {
	result := Context{
		LoggedIn:    false,
		AccessToken: "",
		Client:      http.Client{},
	}

	formValues := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {contextReq.ClientId},
		"refresh_token": {contextReq.RefreshToken},
	}

	authCode := base64.StdEncoding.EncodeToString([]byte(contextReq.ClientId + ":" + contextReq.ClientSecret))

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(formValues.Encode()))
	req.Header.Add("Authorization", "Basic "+authCode)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &result
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &result
	}

	defer resp.Body.Close()
	accessTokenResponse := AccessTokenResponse{}
	jsonErr := json.Unmarshal(body, &accessTokenResponse)
	if jsonErr != nil {
		return &result
	}

	result.LoggedIn = true
	result.AccessToken = string(accessTokenResponse.AccessToken)
	return &result
}

func (context *Context) GetPlaylist(id string) (*Playlist, error) {
	path := fmt.Sprintf("playlists/%s", id)
	req, err := context.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	res, err := context.Client.Do(req)
	if err != nil {
		return nil, err
	}

	playlist := Playlist{}
	err = fromJson(res.Body, &playlist)
	if err != nil {
		return nil, err
	}

	return &playlist, nil
}

func (context *Context) SearchTrack(request SearchTrackRequest) (*SearchTrackResponseItem, error) {
	search := url.QueryEscape(request.Title + " " + request.Artist)
	path := fmt.Sprintf("search?type=track&q=%s", search)
	req, err := context.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	res, err := context.Client.Do(req)
	if err != nil {
		return nil, err
	}

	response := SearchTrackResponse{}
	err = fromJson(res.Body, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Tracks.Items) == 0 {
		return nil, errors.New("provided search yielded no results")
	}

	return &response.Tracks.Items[0], nil
}

func (context *Context) UpdatePlaylist(request UpdatePlaylistRequest) error {
	path := fmt.Sprintf("playlists/%s/tracks", request.Id)
	reqBody, _ := json.Marshal(request)
	req, err := context.createRequest("PUT", path, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	res, err := context.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("invalid response code returned (%d)", res.StatusCode)
	}

	return nil
}

func (context *Context) createRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := baseUrl + path
	request, error := http.NewRequest(method, url, body)
	if error == nil {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", context.AccessToken))
		request.Header.Add("Content-Type", "application/json")
	}
	return request, error
}

// Common serialization pattern
func fromJson(readCloser io.ReadCloser, result any) error {
	defer readCloser.Close()
	body, err := io.ReadAll(readCloser)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	return nil
}
