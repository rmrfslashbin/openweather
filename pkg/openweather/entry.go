/*
Copyright © 2022 Robert Sigler <sigler@improvisedscience.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Const enums for query options/settings
const (
	Current = iota
	Minutely
	Hourly
	Daily
	Alerts
	Standard
	Metric
	Imperial
)

var (
	// Map to store valid languages
	langs map[string]bool

	// Map to store emojis
	Emojis map[string]string
)

// Options for the weather query
type Option func(c *Config)

// Config for the weather query
type Config struct {
	log         *logrus.Logger
	apikey      string
	location    *Location
	excludes    string
	units       string
	lang        string
	rooturl     *url.URL
	iconurlRoot string
}

// Location for the weather query
type Location struct {
	Lat float64
	Lon float64
}

// Weather returns the weather for the given location
type Weather struct {
	Units          string          `json:"units"`
	Lat            float64         `json:"lat"`
	Lon            float64         `json:"lon"`
	Timezone       string          `json:"timezone"`
	TimezoneOffset int             `json:"timezone_offset"`
	Current        *WeatherCurrent `json:"current"`
	Minutely       []struct {
		Dt            int `json:"dt"`
		Precipitation int `json:"precipitation"`
	} `json:"minutely"`
	Hourly *[]WeatherHourly `json:"hourly"`
	Daily  *[]WeatherDaily  `json:"daily"`
	Alerts *[]WeatherAlerts `json:"alerts"`
}

// WeatherCurrent holds the current weather data
type WeatherCurrent struct {
	Dt         int64           `json:"dt"`
	Sunrise    int64           `json:"sunrise"`
	Sunset     int64           `json:"sunset"`
	Temp       float64         `json:"temp"`
	FeelsLike  float64         `json:"feels_like"`
	Pressure   int             `json:"pressure"`
	Humidity   int             `json:"humidity"`
	DewPoint   float64         `json:"dew_point"`
	Clouds     int             `json:"clouds"`
	Uvi        float64         `json:"uvi"`
	Visibility int             `json:"visibility"`
	WindSpeed  float64         `json:"wind_speed"`
	WindGust   float64         `json:"wind_gust"`
	WindDeg    int             `json:"wind_deg"`
	Rain       float64         `json:"rain"`
	Snow       float64         `json:"snow"`
	Weather    []*WeatherStats `json:"weather"`
}

// WeatherHourly holds the hourly weather data
type WeatherHourly struct {
	Dt         int64           `json:"dt"`
	Temp       float64         `json:"temp"`
	FeelsLike  float64         `json:"feels_like"`
	Pressure   int             `json:"pressure"`
	Humidity   int             `json:"humidity"`
	DewPoint   float64         `json:"dew_point"`
	Uvi        float64         `json:"uvi"`
	Clouds     int             `json:"clouds"`
	Visibility int             `json:"visibility"`
	WindSpeed  float64         `json:"wind_speed"`
	WindGust   float64         `json:"wind_gust"`
	WindDeg    int             `json:"wind_deg"`
	Pop        float64         `json:"pop"`
	Rain       float64         `json:"rain"`
	Snow       float64         `json:"snow"`
	Weather    []*WeatherStats `json:"weather"`
}

// WeatherDaily holds the daily weather data
type WeatherDaily struct {
	Dt        int64   `json:"dt"`
	Sunrise   int64   `json:"sunrise"`
	Sunset    int64   `json:"sunset"`
	Moonrise  int64   `json:"moonrise"`
	Moonset   int64   `json:"moonset"`
	MoonPhase float64 `json:"moon_phase"`
	Temp      struct {
		Morn  float64 `json:"morn"`
		Day   float64 `json:"day"`
		Eve   float64 `json:"eve"`
		Night float64 `json:"night"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
	} `json:"temp"`
	FeelsLike struct {
		Morn  float64 `json:"morn"`
		Day   float64 `json:"day"`
		Eve   float64 `json:"eve"`
		Night float64 `json:"night"`
	} `json:"feels_like"`
	Pressure  int             `json:"pressure"`
	Humidity  int             `json:"humidity"`
	DewPoint  float64         `json:"dew_point"`
	WindSpeed float64         `json:"wind_speed"`
	WindGust  float64         `json:"wind_gust"`
	WindDeg   int             `json:"wind_deg"`
	Clouds    int             `json:"clouds"`
	Uvi       float64         `json:"uvi"`
	Pop       float64         `json:"pop"`
	Rain      float64         `json:"rain"`
	Snow      float64         `json:"snow"`
	Weather   []*WeatherStats `json:"weather"`
}

// WeatherAlerts holds the weather alerts
type WeatherAlerts struct {
	SenderName  string   `json:"sender_name"`
	Event       string   `json:"event"`
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// WeatherStats holds the weather stats
type WeatherStats struct {
	ID          int      `json:"id"`
	Main        string   `json:"main"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	IconURL     *url.URL `json:"icon_url"`
}

// Init initializes the weather package
func init() {

	// Set up language list
	langs = make(map[string]bool)
	availableLangs := []string{
		"af", "al", "ar", "az", "bg", "ca", "cz", "da", "de", "el", "en", "eu",
		"fa", "fi", "fr", "gl", "he", "hi", "hr", "hu", "id", "it", "ja", "kr",
		"la", "lt", "mk", "no", "nl", "pl", "pt", "pt_br", "ro", "ru", "sv", "se",
		"sk", "sl", "sp", "es", "sr", "th", "tr", "ua", "uk", "vi", "zh_cn", "zh_tw", "zu",
	}
	for _, lang := range availableLangs {
		langs[lang] = true
	}

	// Set up emoji list
	Emojis = make(map[string]string)
	Emojis["01d"] = "☀️"
	Emojis["01n"] = "🌙"
	Emojis["02d"] = "🌤️"
	Emojis["02n"] = "🌤️"
	Emojis["03d"] = "🌥️"
	Emojis["03n"] = "🌥️"
	Emojis["04d"] = "⛅"
	Emojis["04n"] = "⛅"
	Emojis["09d"] = "⛈️"
	Emojis["09n"] = "⛈️"
	Emojis["10d"] = "🌧️"
	Emojis["10n"] = "🌧️"
	Emojis["11d"] = "🌩️"
	Emojis["11n"] = "🌩️"
	Emojis["13d"] = "❄️"
	Emojis["13n"] = "❄️"
	Emojis["50d"] = "🌫️"
	Emojis["50n"] = "🌫️"
	Emojis["moon_new"] = "🌑"           // > 0.75 <= 1.0 && 0
	Emojis["moon_first_quarter"] = "🌓" // >0 <= 0.25
	Emojis["moon_full"] = "🌕"          // > 0.25 <= 0.5
	Emojis["moon_last_quarter"] = "🌗"  // > 0.5 <= 0.75
}

// New returns a new Config with the given options
func New(opts ...func(*Config)) (*Config, error) {
	config := &Config{}

	// Set up default logger
	config.log = logrus.New()

	// Default to metric
	config.units = "metric"

	// Default to English
	config.lang = "en"

	// Construct the root query URL
	config.rooturl = &url.URL{
		// https://api.openweathermap.org/data/2.5/onecall?lat={lat}&lon={lon}&exclude={part}&appid={API key}
		Scheme: "https",
		Host:   "api.openweathermap.org",
		Path:   "/data/2.5/onecall",
	}

	// OpenWeatherMap image URL
	config.iconurlRoot = "https://openweathermap.org/img/wn/"

	// apply options
	for _, opt := range opts {
		opt(config)
	}

	// apikey must be set
	if config.apikey == "" {
		return nil, fmt.Errorf("apikey is required")
	}

	// location must be set
	if config.location == nil {
		return nil, fmt.Errorf("location is required")
	}

	return config, nil
}

// SetAPIKey sets the API key
func SetAPIKey(apikey string) Option {
	return func(c *Config) {
		c.apikey = apikey
	}
}

// SetExclude sets the exclude list
func SetExcludes(excludes ...int) Option {
	return func(c *Config) {
		excludeList := []string{}
		for _, exclude := range excludes {
			switch exclude {
			case Current:
				excludeList = append(excludeList, "current")
			case Minutely:
				excludeList = append(excludeList, "minutely")
			case Hourly:
				excludeList = append(excludeList, "hourly")
			case Daily:
				excludeList = append(excludeList, "daily")
			case Alerts:
				excludeList = append(excludeList, "alerts")
			default:
				c.log.WithFields(logrus.Fields{
					"excludeProvided": exclude,
				}).Warn("invalid exclude option")
			}
		}
		c.excludes = strings.Join(excludeList, ",")
	}
}

// SetLanguage sets the language
func SetLanguage(lang string) Option {
	return func(c *Config) {
		if langs[lang] {
			c.lang = lang
		} else {
			c.log.WithFields(logrus.Fields{
				"languageProvided": lang,
			}).Warn("invalid language option; defaulting to en")
			c.lang = "en"
		}
	}
}

// SetLocation sets the location
func SetLocation(location *Location) Option {
	return func(c *Config) {
		c.location = location
	}
}

// SetUnits sets the units
func SetUnits(units int) Option {
	return func(c *Config) {
		switch units {
		case Standard:
			c.units = "standard"
		case Metric:
			c.units = "metric"
		case Imperial:
			c.units = "imperial"
		default:
			c.log.WithFields(logrus.Fields{
				"unitProvided": units,
			}).Warn("invalid unit option; defaulting to metric")
			c.units = "metric"
		}
	}
}

// GetOneCallWeather returns the current, minute, hourly, and daily weather plus alerts
func (c *Config) GetOneCallWeather() (*Weather, error) {

	// Construct the query URL
	query := c.rooturl.Query()
	query.Add("lat", fmt.Sprintf("%f", c.location.Lat))
	query.Add("lon", fmt.Sprintf("%f", c.location.Lon))
	query.Add("exclude", c.excludes)
	query.Add("units", c.units)
	query.Add("lang", c.lang)
	query.Add("appid", c.apikey)
	c.rooturl.RawQuery = query.Encode()
	c.log.WithFields(logrus.Fields{
		"url": c.rooturl.String(),
	}).Debug("requesting data")

	// Make the request
	httpResponse, err := http.Get(c.rooturl.String())
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"url":   c.rooturl.String(),
		}).Error("error getting data")
		return nil, err
	}

	// Read the response
	defer httpResponse.Body.Close()
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"url":   c.rooturl.String(),
		}).Error("error reading data")
		return nil, err
	}

	// Parse the response
	weather := &Weather{}
	if err := json.Unmarshal(body, weather); err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"url":   c.rooturl.String(),
		}).Error("error unmarshalling data")
		fmt.Println(string(body))
		return nil, err
	}

	// Add weather icons to current forcast
	for _, v := range weather.Current.Weather {
		if v.Icon != "" {
			v.IconURL, err = url.Parse(c.iconurlRoot + v.Icon + ".png")
			if err != nil {
				c.log.WithFields(logrus.Fields{
					"error": err,
					"url":   c.rooturl.String(),
				}).Error("error parsing icon url")
				return nil, err
			}
		}
	}

	// Add weather icons to hourly forcast
	for _, v := range *weather.Hourly {
		for _, vv := range v.Weather {
			if vv.Icon != "" {
				vv.IconURL, err = url.Parse(c.iconurlRoot + vv.Icon + ".png")
				if err != nil {
					c.log.WithFields(logrus.Fields{
						"error": err,
						"url":   c.rooturl.String(),
					}).Error("error parsing icon url")
					return nil, err
				}
			}
		}
	}

	// Add weather icons to daily forcast
	for _, v := range *weather.Daily {
		for _, vv := range v.Weather {
			if vv.Icon != "" {
				vv.IconURL, err = url.Parse(c.iconurlRoot + vv.Icon + ".png")
				if err != nil {
					c.log.WithFields(logrus.Fields{
						"error": err,
						"url":   c.rooturl.String(),
					}).Error("error parsing icon url")
					return nil, err
				}
			}
		}
	}

	weather.Units = c.units

	return weather, nil
}

// ToJSON returns the weather as a JSON byte array
func (w *Weather) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}

// ToToml returns the weather as a TOML byte array
func (w *Weather) ToToml() ([]byte, error) {
	return toml.Marshal(w)
}

// ToYAML returns the weather as a YAML byte array
func (w *Weather) ToYAML() ([]byte, error) {
	return yaml.Marshal(w)
}
