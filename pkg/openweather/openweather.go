package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog"
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
	langs map[string]struct{}

	// Map to store emojis
	Emojis map[string]string
)

// Options for the weather query
type Option func(c *Openweather)

// Openweather for the weather query
type Openweather struct {
	log         *zerolog.Logger
	apikey      string
	location    *Location
	excludes    string
	units       string
	lang        string
	rooturl     *url.URL
	iconurlRoot string
}

// New returns a new Config with the given options
func New(opts ...func(*Openweather)) (*Openweather, error) {
	cfg := &Openweather{}

	// Default to metric
	cfg.units = "metric"

	// Default to English
	cfg.lang = "en"

	// Construct the root query URL
	cfg.rooturl = &url.URL{
		// https://api.openweathermap.org/data/3.0/onecall?lat={lat}&lon={lon}&exclude={part}&appid={API key}
		Scheme: "https",
		Host:   "api.openweathermap.org",
		Path:   "/data/2.5/onecall",
	}

	// OpenWeatherMap image URL
	cfg.iconurlRoot = "https://openweathermap.org/img/wn/"

	// apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// set up logger if not provided
	if cfg.log == nil {
		log := zerolog.New(os.Stderr).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		cfg.log = &log
	}

	// apikey must be set
	if cfg.apikey == "" {
		return nil, &ErrNoAPIKey{}
	}

	// location must be set
	if cfg.location == nil {
		return nil, &ErrNoLocation{}
	}

	return cfg, nil
}

// WithAPIKey sets the API key
func WithAPIKey(apikey string) Option {
	return func(c *Openweather) {
		c.apikey = apikey
	}
}

// WithExcludes sets the exclude list
func WithExcludes(excludes ...int) Option {
	return func(c *Openweather) {
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
			}
		}
		c.excludes = strings.Join(excludeList, ",")
	}
}

// WithLanguage sets the language
func WithLanguage(lang string) Option {
	return func(c *Openweather) {
		if _, ok := langs[lang]; ok {
			c.lang = lang
		}
	}
}

// WithLocation sets the location
func WithLocation(location *Location) Option {
	return func(c *Openweather) {
		c.location = location
	}
}

// WithLogger sets the logger
func WithLogger(log *zerolog.Logger) Option {
	return func(c *Openweather) {
		c.log = log
	}
}

// WithUnits sets the units
func WithUnits(units int) Option {
	return func(c *Openweather) {
		switch units {
		case Standard:
			c.units = "standard"
		case Metric:
			c.units = "metric"
		case Imperial:
			c.units = "imperial"
		}
	}
}

// GetOneCallWeather returns the current, minute, hourly, and daily weather plus alerts
func (c *Openweather) GetOneCallWeather() (*Weather, error) {

	// Construct the query URL
	query := c.rooturl.Query()
	query.Add("lat", fmt.Sprintf("%f", c.location.Lat))
	query.Add("lon", fmt.Sprintf("%f", c.location.Lon))
	query.Add("exclude", c.excludes)
	query.Add("units", c.units)
	query.Add("lang", c.lang)
	query.Add("appid", c.apikey)
	c.rooturl.RawQuery = query.Encode()
	c.log.Debug().
		Str("url", c.rooturl.String()).
		Msg("requesting data")

	// Make the request
	httpResponse, err := http.Get(c.rooturl.String())
	if err != nil {
		c.log.Error().
			Str("url", c.rooturl.String()).
			Msg("error getting data")
		return nil, err
	}

	// Read the response
	defer httpResponse.Body.Close()
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		c.log.Error().
			Str("url", c.rooturl.String()).
			Msg("error reading data")
		return nil, err
	}
	if httpResponse.StatusCode != http.StatusOK {
		errMsg := &ErrorResponse{}
		if err := json.Unmarshal(body, errMsg); err != nil {
			c.log.Error().
				Str("url", c.rooturl.String()).
				Str("status", httpResponse.Status).
				Str("body", string(body)).
				Msg("error proccessing http error")
			return nil, err
		}
		c.log.Error().
			Str("url", c.rooturl.String()).
			Msg("error getting data")
		return nil, &ErrAPIError{
			Code: httpResponse.StatusCode,
			Msg:  errMsg.Message,
		}
	}

	// Parse the response
	weather := &Weather{}
	if err := json.Unmarshal(body, weather); err != nil {
		c.log.Error().
			Str("url", c.rooturl.String()).
			Msg("error unmarshalling data")
		fmt.Println(string(body))
		return nil, err
	}
	fmt.Println(string(body))
	spew.Dump(weather)

	// Add weather icons to current forcast
	for _, v := range weather.Current.Weather {
		if v.Icon != "" {
			v.IconURL, err = url.Parse(c.iconurlRoot + v.Icon + ".png")
			if err != nil {
				c.log.Error().
					Str("url", c.rooturl.String()).
					Msg("error parsing icon url")
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
					c.log.Error().
						Str("url", c.rooturl.String()).
						Msg("error parsing icon url")
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
					c.log.Error().
						Str("url", c.rooturl.String()).
						Msg("error parsing icon url")
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

func (weather *Weather) Text(brief bool) error {
	// Get the times
	dt := time.Unix(weather.Current.Dt, 0)
	sunrise := time.Unix(weather.Current.Sunrise, 0)
	sunset := time.Unix(weather.Current.Sunset, 0)

	// Set up the units output
	unit := "°C"
	speed := "m/s"
	if weather.Units == "standard" {
		unit = "°K"
	} else if weather.Units == "imperial" {
		unit = "°F"
		speed = "mph"
	}

	if brief {
		fmt.Printf("Current weather as of %s\n", dt.Local())
		fmt.Printf("  %s %s (%s) Temperature: %.1f%s Feels like: %1.f%s\n",
			Emojis[weather.Current.Weather[0].Icon],
			weather.Current.Weather[0].Main,
			weather.Current.Weather[0].Description,
			weather.Current.Temp, unit,
			weather.Current.FeelsLike, unit,
		)
		fmt.Printf("  Wind speed: %.1f %s from %d°\n", weather.Current.WindSpeed, speed, weather.Current.WindDeg)
		fmt.Printf("  Cloudiness: %d%% UV index: %.1f\n", weather.Current.Clouds, weather.Current.Uvi)

		for _, day := range *weather.Daily {
			ts := time.Unix(day.Dt, 0)
			fmt.Printf("\n%s (%d %s %02d)\n", ts.Local().Weekday(), ts.Local().Year(), ts.Local().Month(), ts.Local().Day())
			fmt.Printf("  %s %s (%s) High %.1f%s Low %.1f%s with %.1f%% chance of precipitation\n",
				Emojis[day.Weather[0].Icon],
				day.Weather[0].Main,
				day.Weather[0].Description,
				day.Temp.Max, unit, day.Temp.Min, unit,
				day.Pop,
			)
		}

		for i, hour := range *weather.Hourly {
			// Only show the next 12 hours
			if i >= 12 {
				break
			}
			ts := time.Unix(hour.Dt, 0)
			fmt.Printf("\n%s (%d %s %02d %02d:%02d) %s %s Temp: %.1f%s Wind: %.1f %s Precip: %.1f%%\n",
				ts.Local().Weekday(), ts.Local().Year(), ts.Local().Month(), ts.Local().Day(), ts.Local().Hour(), ts.Local().Minute(),
				Emojis[hour.Weather[0].Icon],
				hour.Weather[0].Description,
				hour.Temp, unit,
				hour.WindSpeed, speed,
				hour.Pop,
			)

		}

	} else {

		// Print the current weather conditions
		fmt.Printf("Current weather for %f, %f as of %s\n", weather.Lat, weather.Lon, dt.Local())
		fmt.Printf("%s %s (%s)\n",
			Emojis[weather.Current.Weather[0].Icon],
			weather.Current.Weather[0].Main,
			weather.Current.Weather[0].Description,
		)
		fmt.Printf("  Temperature: %.1f%s\n", weather.Current.Temp, unit)
		fmt.Printf("  Feels like: %.1f%s\n", weather.Current.FeelsLike, unit)
		fmt.Printf("  Humidity: %d%%\n", weather.Current.Humidity)
		fmt.Printf("  Pressure: %d hPa\n", weather.Current.Pressure)
		fmt.Printf("  Due point: %.1f%s\n", weather.Current.DewPoint, unit)
		fmt.Printf("  Wind speed: %.1f %s\n", weather.Current.WindSpeed, speed)
		fmt.Printf("  Wind gust: %.1f %s\n", weather.Current.WindGust, speed)
		fmt.Printf("  Wind direction: %d°\n", weather.Current.WindDeg)
		fmt.Printf("  Cloudiness: %d%%\n", weather.Current.Clouds)
		fmt.Printf("  Rain: %.1f mm\n", weather.Current.Rain)
		fmt.Printf("  Snow: %.1f mm\n", weather.Current.Snow)
		fmt.Printf("  UV index: %.1f\n", weather.Current.Uvi)
		fmt.Printf("  Visibility: %d m\n", weather.Current.Visibility)
		fmt.Printf("  Sunrise: %s\n", sunrise.Local())
		fmt.Printf("  Sunset: %s\n", sunset.Local())

		// If alerts, print them
		if weather.Alerts != nil {
			fmt.Println("\nAlerts:")
			for _, alert := range *weather.Alerts {
				start := time.Unix(alert.Start, 0)
				end := time.Unix(alert.End, 0)
				fmt.Println("---")
				fmt.Printf("  %s :: %s\n", alert.SenderName, alert.Event)
				fmt.Printf("  From %s :: Until %s\n", start.Local(), end.Local())
				fmt.Printf("  %s\n", alert.Description)
				fmt.Println("---")
			}
		}
		fmt.Println()

		// If daily forecast, print it
		for _, day := range *weather.Daily {
			ts := time.Unix(day.Dt, 0)
			sunrise := time.Unix(day.Sunrise, 0)
			sunset := time.Unix(day.Sunset, 0)
			moonrise := time.Unix(day.Moonrise, 0)
			moonset := time.Unix(day.Moonset, 0)
			fmt.Printf("%s (%d %s %02d)\n", ts.Local().Weekday(), ts.Local().Year(), ts.Local().Month(), ts.Local().Day())
			fmt.Printf("  %s %s (%s)\n",
				Emojis[day.Weather[0].Icon],
				day.Weather[0].Main,
				day.Weather[0].Description,
			)
			fmt.Printf("  High %.1f%s Low %.1f%s\n", day.Temp.Max, unit, day.Temp.Min, unit)
			fmt.Printf("  Morning: %.1f%s (%.1f%s)\n", day.Temp.Morn, unit, day.FeelsLike.Morn, unit)
			fmt.Printf("  Day: %.1f%s (%.1f%s)\n", day.Temp.Day, unit, day.FeelsLike.Day, unit)
			fmt.Printf("  Evening: %.1f%s (%.1f%s)\n", day.Temp.Eve, unit, day.FeelsLike.Eve, unit)
			fmt.Printf("  Night: %.1f%s (%.1f%s)\n", day.Temp.Night, unit, day.FeelsLike.Night, unit)
			fmt.Printf("  Wind speed: %.1f %s (gust %.1f%s) from %d°\n", day.WindSpeed, speed, day.WindGust, speed, day.WindDeg)
			fmt.Printf("  Cloudiness: %d%% UV: %.1f\n", day.Clouds, day.Uvi)
			fmt.Printf("  Probability of precipitation: %.1f%%\n", day.Pop)
			fmt.Printf("  Rain: %.1f mm Snow: %.1f mm\n", day.Rain, day.Snow)
			fmt.Printf("  Sunrise (%s) Sunset (%s)\n", sunrise.Local(), sunset.Local())
			fmt.Printf("  Moonrise (%s) Moonset (%s)\n", moonrise.Local(), moonset.Local())

			fmt.Println()
		}
		for _, hour := range *weather.Hourly {
			ts := time.Unix(hour.Dt, 0)
			fmt.Printf("\n%s (%d %s %02d %02d:%02d) %s %s Temp: %.1f%s Wind: %.1f %s Precip: %.1f%%\n",
				ts.Local().Weekday(), ts.Local().Year(), ts.Local().Month(), ts.Local().Day(), ts.Local().Hour(), ts.Local().Minute(),
				Emojis[hour.Weather[0].Icon],
				hour.Weather[0].Description,
				hour.Temp, unit,
				hour.WindSpeed, speed,
				hour.Pop,
			)

		}
	}
	return nil
}
