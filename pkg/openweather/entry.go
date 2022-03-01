package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

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
	langs map[string]bool
)

// Used to manage varidic options
type Option func(c *Config)

// database configs
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

type Location struct {
	Lat float64
	Lon float64
}

type Response struct {
	Lat            float64          `json:"lat"`
	Lon            float64          `json:"lon"`
	Timezone       string           `json:"timezone"`
	TimezoneOffset int              `json:"timezone_offset"`
	Current        *ResponseCurrent `json:"current"`
	Minutely       []struct {
		Dt            int `json:"dt"`
		Precipitation int `json:"precipitation"`
	} `json:"minutely"`
	Hourly *[]ResponseHourly `json:"hourly"`
	Daily  *[]ResponseDaily  `json:"daily"`
	Alerts *[]ResponseAlerts `json:"alerts"`
}

type ResponseCurrent struct {
	Dt         int64              `json:"dt"`
	Sunrise    int64              `json:"sunrise"`
	Sunset     int64              `json:"sunset"`
	Temp       float64            `json:"temp"`
	FeelsLike  float64            `json:"feels_like"`
	Pressure   int                `json:"pressure"`
	Humidity   int                `json:"humidity"`
	DewPoint   float64            `json:"dew_point"`
	Clouds     int                `json:"clouds"`
	Uvi        float64            `json:"uvi"`
	Visibility int                `json:"visibility"`
	WindSpeed  float64            `json:"wind_speed"`
	WindGust   float64            `json:"wind_gust"`
	WindDeg    int                `json:"wind_deg"`
	Rain       float64            `json:"rain"`
	Snow       float64            `json:"snow"`
	Weather    []*ResponseWeather `json:"weather"`
}

type ResponseHourly struct {
	Dt         int64              `json:"dt"`
	Temp       float64            `json:"temp"`
	FeelsLike  float64            `json:"feels_like"`
	Pressure   int                `json:"pressure"`
	Humidity   int                `json:"humidity"`
	DewPoint   float64            `json:"dew_point"`
	Uvi        float64            `json:"uvi"`
	Clouds     int                `json:"clouds"`
	Visibility int                `json:"visibility"`
	WindSpeed  float64            `json:"wind_speed"`
	WindGust   float64            `json:"wind_gust"`
	WindDeg    int                `json:"wind_deg"`
	Pop        float64            `json:"pop"`
	Rain       float64            `json:"rain"`
	Snow       float64            `json:"snow"`
	Weather    []*ResponseWeather `json:"weather"`
}

type ResponseDaily struct {
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
	Pressure  int                `json:"pressure"`
	Humidity  int                `json:"humidity"`
	DewPoint  float64            `json:"dew_point"`
	WindSpeed float64            `json:"wind_speed"`
	WindGust  float64            `json:"wind_gust"`
	WindDeg   int                `json:"wind_deg"`
	Clouds    int                `json:"clouds"`
	Uvi       float64            `json:"uvi"`
	Pop       float64            `json:"pop"`
	Rain      float64            `json:"rain"`
	Snow      float64            `json:"snow"`
	Weather   []*ResponseWeather `json:"weather"`
}

type ResponseAlerts struct {
	SenderName  string   `json:"sender_name"`
	Event       string   `json:"event"`
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type ResponseWeather struct {
	ID          int      `json:"id"`
	Main        string   `json:"main"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	IconURL     *url.URL `json:"icon_url"`
}

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
}

// New returns a new Config with the given options
func New(opts ...func(*Config)) (*Config, error) {
	config := &Config{}

	// Set up default logger
	config.log = logrus.New()
	config.units = "metric"
	config.lang = "en"

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

	// https://api.openweathermap.org/data/2.5/onecall?lat={lat}&lon={lon}&exclude={part}&appid={API key}
	config.rooturl = &url.URL{
		Scheme: "https",
		Host:   "api.openweathermap.org",
		Path:   "/data/2.5/onecall",
	}

	config.iconurlRoot = "https://openweathermap.org/img/wn/"
	return config, nil
}

func SetAppId(apikey string) Option {
	return func(c *Config) {
		c.apikey = apikey
	}
}

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

func SetLocation(location *Location) Option {
	return func(c *Config) {
		c.location = location
	}
}

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

func (c *Config) Run() error {

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

	resp, err := http.Get(c.rooturl.String())
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"url":   c.rooturl.String(),
		}).Error("error getting data")
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"url":   c.rooturl.String(),
		}).Error("error reading data")
		return err
	}

	res := &Response{}
	if err := json.Unmarshal(body, res); err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"url":   c.rooturl.String(),
		}).Error("error unmarshalling data")
		fmt.Println(string(body))
		return err
	}
	for _, v := range res.Current.Weather {
		if v.Icon != "" {
			v.IconURL, err = url.Parse(c.iconurlRoot + v.Icon + ".png")
			if err != nil {
				c.log.WithFields(logrus.Fields{
					"error": err,
					"url":   c.rooturl.String(),
				}).Error("error parsing icon url")
				return err
			}
		}
	}

	for _, v := range *res.Hourly {
		for _, vv := range v.Weather {
			if vv.Icon != "" {
				vv.IconURL, err = url.Parse(c.iconurlRoot + vv.Icon + ".png")
				if err != nil {
					c.log.WithFields(logrus.Fields{
						"error": err,
						"url":   c.rooturl.String(),
					}).Error("error parsing icon url")
					return err
				}
			}
		}
	}

	for _, v := range *res.Daily {
		for _, vv := range v.Weather {
			if vv.Icon != "" {
				vv.IconURL, err = url.Parse(c.iconurlRoot + vv.Icon + ".png")
				if err != nil {
					c.log.WithFields(logrus.Fields{
						"error": err,
						"url":   c.rooturl.String(),
					}).Error("error parsing icon url")
					return err
				}
			}
		}
	}

	spew.Dump(res)

	return nil
}
