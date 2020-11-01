/*
Origin : https://github.com/ramsgoli/Golang-OpenWeatherMap
Hardcoded unit from imperial to metric
*/

package openweathermap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ernesto-jimenez/httplogger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
)

/*
Define API response fields
*/
type OpenWeatherMap struct {
	APIKEY string
}

/*
Return response fields of City data
*/
type City struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

/*
Geographic spherical coordinate for input
*/
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

/*
Return weather forecast
*/
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

/*
Return wind condition
*/
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

/*
Return cloud condition
*/
type Clouds struct {
	All int `json:"all"`
}

/*
Return rain condition
*/
type Rain struct {
	Threehr int `json:"3h"`
}

/*
Return temperature data
*/
type Main struct {
	Temp     float64 `json:"temp"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}

/*
Define API response objects (compose of the above fields)
*/
type CurrentWeatherResponse struct {
	Coord   `json:"coord"`
	Weather []Weather `json:"weather"`
	Main    `json:"main"`
	Wind    `json:"wind"`
	Rain    `json:"rain"`
	Clouds  `json:"clouds"`
	DT      int    `json:"dt"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
}

/*
Response from openweathermap
*/
type ForecastResponse struct {
	City    `json:"city"`
	Coord   `json:"coord"`
	Country string `json:"country"`
	List    []struct {
		DT      int `json:"dt"`
		Main    `json:"main"`
		Weather `json:"weather"`
		Clouds  `json:"clouds"`
		Wind    `json:"wind"`
	} `json:"list"`
}

/*
httpLogger log http request response
*/
type httpLogger struct {
	log *log.Logger
}

/*
openweathermap endpoint
*/
const (
	APIURL string = "api.openweathermap.org"
)

func newLogger() *httpLogger {
	return &httpLogger{
		log: log.New(os.Stderr, "log - ", log.LstdFlags),
	}
}

func (l *httpLogger) LogRequest(req *http.Request) {
	l.log.Printf(
		"Request %s %s\n",
		req.Method,
		req.URL.String(),
	)
	l.log.Printf(
		"Request %+q %s",
		req.Header["User-Agent"],
		req.URL.String(),
	)
}

func (l *httpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		l.log.Println(err)
	} else {
		l.log.Printf(
			"Response method=%s status=%d durationMs=%d %s",
			req.Method,
			res.StatusCode,
			duration,
			req.URL.String(),
		)
	}
}

/*
Build request to openweathermap
*/
func makeAPIRequest(ctx context.Context, url string) ([]byte, error) {
	// Build an http client so we can have control over timeout
	client := &http.Client{
		Timeout:   time.Second * 2,
		Transport: httplogger.NewLoggedTransport(http.DefaultTransport, newLogger()),
	}

	req, getErr := http.NewRequest("GET", url, nil)
	if getErr != nil {
		return nil, getErr
	}

	ctx, req = otelhttptrace.W3C(ctx, req)
	otelhttptrace.Inject(ctx, req)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// defer the closing of the res body
	defer res.Body.Close()

	// read the http response body into a byte stream
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}

/*
Get current weather from a city - openweathermap
*/
func (owm *OpenWeatherMap) CurrentWeatherFromCity(city string) (*CurrentWeatherResponse, error) {
	ctx := context.Background()

	if owm.APIKEY == "" {
		// No API keys present, return error
		return nil, errors.New("No API keys present")
	}
	url := fmt.Sprintf("http://%s/data/2.5/weather?q=%s&units=metric&APPID=%s", APIURL, city, owm.APIKEY)

	body, err := makeAPIRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	var cwr CurrentWeatherResponse

	// unmarshal the byte stream into a Go data type
	jsonErr := json.Unmarshal(body, &cwr)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &cwr, nil
}

/*
Get current weather from a coordinate - openweathermap
*/
func (owm *OpenWeatherMap) CurrentWeatherFromCoordinates(lat, long float64) (*CurrentWeatherResponse, error) {
	ctx := context.Background()

	if owm.APIKEY == "" {
		// No API keys present, return error
		return nil, errors.New("No API keys present")
	}

	url := fmt.Sprintf("http://%s/data/2.5/weather?lat=%f&lon=%f&units=metric&APPID=%s", APIURL, lat, long, owm.APIKEY)

	body, err := makeAPIRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var cwr CurrentWeatherResponse

	// unmarshal the byte stream into a Go data type
	jsonErr := json.Unmarshal(body, &cwr)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &cwr, nil
}

/*
Return current weather from a zip code - openweathermap
*/
func (owm *OpenWeatherMap) CurrentWeatherFromZip(zip int) (*CurrentWeatherResponse, error) {
	ctx := context.Background()

	if owm.APIKEY == "" {
		// No API keys present, return error
		return nil, errors.New("No API keys present")
	}
	url := fmt.Sprintf("http://%s/data/2.5/weather?zip=%d&units=metric&APPID=%s", APIURL, zip, owm.APIKEY)

	body, err := makeAPIRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	var cwr CurrentWeatherResponse

	// unmarshal the byte stream into a Go data type
	jsonErr := json.Unmarshal(body, &cwr)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &cwr, nil
}

/*
Return current weather from a city id - openweathermap
*/
func (owm *OpenWeatherMap) CurrentWeatherFromCityId(id int) (*CurrentWeatherResponse, error) {
	ctx := context.Background()

	if owm.APIKEY == "" {
		// No API keys present, return error
		return nil, errors.New("No API keys present")
	}
	url := fmt.Sprintf("http://%s/data/2.5/weather?id=%d&units=metric&APPID=%s", APIURL, id, owm.APIKEY)

	body, err := makeAPIRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	var cwr CurrentWeatherResponse

	// unmarshal the byte stream into a Go data type
	jsonErr := json.Unmarshal(body, &cwr)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &cwr, nil
}
