package wttrin

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/acarl005/stripansi"
)

const (
	baseURL = "https://wttr.in/"

	curlUserAgent = "curl/7.54.0"
)

// WeathersText returns weathers for 3 days in plain text
func WeathersText(place string) (result string, err error) {
	return httpGet(baseURL+url.QueryEscape(place), false)
}

// WeathersHTML returns weathers for 3 days in HTML
func WeathersHTML(place string) (result string, err error) {
	return httpGet(baseURL+url.QueryEscape(place), true)
}

// WeatherTextForToday returns today's weather for given place in plain text
func WeatherTextForToday(place string) (result string, err error) {
	return httpGet(baseURL+url.QueryEscape(place)+"?0", false)
}

// WeatherHTMLForToday returns today's weather for given place in HTML
func WeatherHTMLForToday(place string) (result string, err error) {
	return httpGet(baseURL+url.QueryEscape(place)+"?0", true)
}

func httpGet(url string, asHTML bool) (result string, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 300 * time.Second,
			}).Dial,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	var req *http.Request

	if req, err = http.NewRequest("GET", url, nil); err == nil {
		if !asHTML {
			req.Header.Set("User-Agent", curlUserAgent)
		}

		var resp *http.Response
		if resp, err = client.Do(req); err == nil {
			defer resp.Body.Close()

			var body []byte
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				if !asHTML {
					return stripansi.Strip(string(body)), nil
				}

				return string(body), nil
			}
		}
	}

	return "", err
}
