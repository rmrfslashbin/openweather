package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rmrfslashbin/openweather/pkg/geocode"
	"github.com/rmrfslashbin/openweather/pkg/openweather"
	"github.com/rs/zerolog"
)

const (
	// APP_NAME is the name of the application
	APP_NAME = "openweather"
)

// Context is used to pass context/global configs to the commands
type Context struct {
	// log is the logger
	log    *zerolog.Logger
	apikey string
}

// CurrentCmd updates the GTFS feed specs
type CurrentCmd struct {
	Metric   bool    `name:"metric" required:"" group:"unit" xor:"unit" help:"Use metric units."`
	Imperial bool    `name:"imperial" required:"" group:"unit" xor:"unit" help:"Use imperial units."`
	Standard bool    `name:"standard" required:"" group:"unit" xor:"unit" help:"Use standard units."`
	Lat      float64 `name:"lat" env:"LAT" required:"" help:"Latitude."`
	Lon      float64 `name:"lon" env:"LON" required:"" help:"Longitude."`
	Json     bool    `name:"json" required:"" group:"output" xor:"output" help:"Output the results as JSON."`
	Yaml     bool    `name:"yaml" required:"" group:"output" xor:"output" help:"Output the results as YAML."`
	Toml     bool    `name:"toml" required:"" group:"output" xor:"output" help:"Output the results as TOML."`
	Text     bool    `name:"text" required:"" group:"output" xor:"output" help:"Output the results as text."`
	Brief    bool    `name:"brief"  help:"Output brief text results."`
}

// Run is the entry point for the CurrentCmd command
func (r *CurrentCmd) Run(ctx *Context) error {
	units := openweather.Metric
	if r.Metric {
		units = openweather.Metric
	} else if r.Imperial {
		units = openweather.Imperial
	} else if r.Standard {
		units = openweather.Standard
	}

	// Set up the OpenWeatherMap client
	ow, err := openweather.New(
		openweather.WithAPIKey(ctx.apikey),
		openweather.WithLocation(&openweather.Location{
			Lat: r.Lat,
			Lon: r.Lon,
		}),
		openweather.WithLogger(ctx.log),
		openweather.WithUnits(units),
	)
	if err != nil {
		return err
	}

	// Fetch the current weather conditions
	weather, err := ow.GetOneCallWeather()
	if err != nil {
		return err
	}

	if r.Json {
		if bytes, err := weather.ToJSON(); err != nil {
			return err
		} else {
			fmt.Println(string(bytes))
		}
	} else if r.Yaml {
		if bytes, err := weather.ToYAML(); err != nil {
			return err
		} else {
			fmt.Println(string(bytes))
		}
	} else if r.Toml {
		if bytes, err := weather.ToToml(); err != nil {
			return err
		} else {
			fmt.Println(string(bytes))
		}
	} else if r.Text {
		weather.Text(r.Brief)
	}

	return nil
}

// GeoLookupCmd looks up the location of a zip/post code or city
type GeoLookupCmd struct {
	Zip  string `name:"zip" required:"" group:"by" xor:"by" help:"Zip/post code and country code divided by comma. Please use ISO 3166 country codes. (ex: 30318 or 30318,US)"`
	City string `name:"city" required:"" group:"by" xor:"by" help:"City name, state code (only for the US) and country code divided by comma. Please use ISO 3166 country codes. (ex: Atlanta or Atlanta,US or Atlanta,GA,US)"`
	Json bool   `name:"json" required:"" group:"output" xor:"output" help:"Output the results as JSON."`
	Yaml bool   `name:"yaml" required:"" group:"output" xor:"output" help:"Output the results as YAML."`
	Toml bool   `name:"toml" required:"" group:"output" xor:"output" help:"Output the results as TOML."`
	Text bool   `name:"text" required:"" group:"output" xor:"output" help:"Output the results as text."`
	Lang string `name:"lang" default:"en" help:"Language code for the output."`
}

// Run is the entry point for the GeoLookupCmd command
func (r *GeoLookupCmd) Run(ctx *Context) error {
	// Set up the OpenWeatherMap client
	gc, err := geocode.New(
		geocode.WithAPIKey(ctx.apikey),
	)
	if err != nil {
		return err
	}

	if r.Zip != "" {
		loc, err := gc.ByZip(r.Zip)
		if err != nil {
			return err
		}
		if r.Json {
			if bytes, err := loc.ToJSON(); err != nil {
				return err
			} else {
				fmt.Println(string(bytes))
			}
		} else if r.Yaml {
			if bytes, err := loc.ToYAML(); err != nil {
				return err
			} else {
				fmt.Println(string(bytes))
			}
		} else if r.Toml {
			if bytes, err := loc.ToToml(); err != nil {
				return err
			} else {
				fmt.Println(string(bytes))
			}
		} else if r.Text {
			fmt.Println("Name:    ", loc.Name)
			fmt.Println("Country: ", loc.Country)
			fmt.Println("Zip:     ", loc.Zip)
			fmt.Println("Lat:     ", loc.Lat)
			fmt.Println("Lon:     ", loc.Lon)
		}

	} else if r.City != "" {
		loc, err := gc.ByCity(r.City)
		if err != nil {
			return err
		}
		if r.Json {
			if bytes, err := loc.ToJSON(); err != nil {
				return err
			} else {
				fmt.Println(string(bytes))
			}
		} else if r.Yaml {
			if bytes, err := loc.ToYAML(); err != nil {
				return err
			} else {
				fmt.Println(string(bytes))
			}
		} else if r.Toml {
			if bytes, err := loc.ToToml(); err != nil {
				return err
			} else {
				fmt.Println(string(bytes))
			}
		} else if r.Text {
			for _, l := range loc.Entities {
				fmt.Println("Name:    ", l.Name)
				fmt.Println("Country: ", l.Country)
				fmt.Println("Lat:     ", l.Lat)
				fmt.Println("Lon:     ", l.Lon)
				if len(loc.Entities) > 1 {
					fmt.Println()
				}
			}
		}
	}

	return nil
}

// CLI is the main CLI struct
type CLI struct {
	// Global flags/args
	APIKey   string `name:"apikey" env:"APIKEY" required:"" help:"The OpenWeatherMap API key."`
	LogLevel string `name:"loglevel" env:"LOGLEVEL" default:"error" enum:"panic,fatal,error,warn,info,debug,trace" help:"Set the log level."`

	Current CurrentCmd   `cmd:"" help:"Get current weather conditions."`
	Lookup  GeoLookupCmd `cmd:"" help:"Lookup lat/lon data for a location."`
}

func main() {
	var err error

	// Set up the logger
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Parse the command line
	var cli CLI
	ctx := kong.Parse(&cli)

	// Set up the logger's log level
	// Default to info via the CLI args
	switch cli.LogLevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	// Log some start up stuff for debugging
	log.Info().
		Str("app_name", APP_NAME).
		Str("log_level", cli.LogLevel).
		Msg("starting up")

	// Call the Run() method of the selected parsed command.
	err = ctx.Run(&Context{
		apikey: cli.APIKey,
		log:    &log,
	})

	// FatalIfErrorf terminates with an error message if err != nil
	ctx.FatalIfErrorf(err)
}
