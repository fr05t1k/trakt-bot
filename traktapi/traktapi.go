package traktapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const ApiVersion = "2"
const ApiBaseUrl = "https://api.trakt.tv/"
const ApiOAuthDeviceCode = "oauth/device/code"
const ApiOauthDeviceToken = "oauth/device/token"
const ApiCalendarsAllDvd = "/calendars/all/dvd/%s/%s"
const ContentTypeJson = "application/json"

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUrl string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type Client struct {
	httpClient   *http.Client
	clientId     string
	clientSecret string
	accessToken  string
	refreshToken string
}

type DeviceCodeRequest struct {
	ClientId string `json:"client_id"`
}

type DeviceTokenRequest struct {
	ClientId     string `json:"client_id"`
	Code         string `json:"code"`
	ClientSecret string `json:"client_secret"`
}

type DeviceTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Release struct {
	Released string `json:"released"`
	Movie    Movie
}

type Movie struct {
	Title   string `json:"title"`
	Year    int    `json:"year"`
	Trailer string `json:"trailer"`
	Ids     Ids    `json:"ids"`
}

type Ids struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	Imdb  string `json:"imdb"`
	Tmdb  int    `json:"tmdb"`
}

func (c *Client) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

func (c *Client) DeviceCode() (deviceCode DeviceCodeResponse, err error) {
	body, err := json.Marshal(DeviceCodeRequest{ClientId: c.clientId})
	if err != nil {
		return
	}

	response, err := c.httpClient.Post(ApiBaseUrl+ApiOAuthDeviceCode, ContentTypeJson, bytes.NewReader(body))
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("response code: %d", response.StatusCode)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&deviceCode)

	return
}

func (c *Client) GetDeviceToken(code string) (tokenResponse DeviceTokenResponse, err error) {
	requestBody, err := json.Marshal(DeviceTokenRequest{
		ClientId:     c.clientId,
		Code:         code,
		ClientSecret: c.clientSecret,
	})
	if err != nil {
		return
	}

	response, err := c.httpClient.Post(ApiBaseUrl+ApiOauthDeviceToken, ContentTypeJson, bytes.NewReader(requestBody))
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		return tokenResponse, fmt.Errorf("response code: %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&tokenResponse)
	return
}

func (c *Client) GetDvdReleases() (releases []Release, err error) {
	request, err := c.makeRequestWithAccessToken(
		http.MethodGet,
		ApiBaseUrl+fmt.Sprintf(ApiCalendarsAllDvd, time.Now().AddDate(0, 0, 0).Format("2006-01-02"), "30"),
		nil,
	)
	if err != nil {
		return
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return
	}

	err = json.NewDecoder(response.Body).Decode(&releases)
	return
}

func (c *Client) SaveDeviceToken(deviceToken DeviceTokenResponse) {
	c.accessToken = deviceToken.AccessToken
	c.refreshToken = deviceToken.RefreshToken
}

func (c *Client) makeRequestWithAccessToken(method, url string, body io.Reader) (request *http.Request, err error) {
	request, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	request.Header.Set("trakt-api-version", ApiVersion)
	request.Header.Set("trakt-api-key", c.clientId)

	return request, err
}

func NewClient(clientId string, clientSecret string) *Client {

	return &Client{
		httpClient:   &http.Client{},
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}
