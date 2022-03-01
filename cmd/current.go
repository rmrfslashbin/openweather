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
package cmd

import (
	"fmt"
	"time"

	"github.com/rmrfslashbin/openweather/pkg/openweather"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current weather conditions",
	Long:  `Show the current weather conditions along with daily forecast.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Catch errors
		var err error
		defer func() {
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("main crashed")
			}
		}()
		if err := getCurrent(); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error")
		}
	},
}

// Init adds the flags and configures the command
func init() {
	rootCmd.AddCommand(currentCmd)
}

// getCurrent fetches the current weather conditions
func getCurrent() error {
	// Get/set the requested units
	var units int
	switch viper.GetString("units") {
	case "metric":
		units = openweather.Metric
	case "imperial":
		units = openweather.Imperial
	case "standard":
		units = openweather.Standard
	default:
		units = openweather.Metric
	}

	// Set up the OpenWeatherMap client
	ow, err := openweather.New(
		openweather.SetAPIKey(viper.GetString("apikey")),
		openweather.SetLocation(&openweather.Location{
			Lat: viper.GetFloat64("lat"),
			Lon: viper.GetFloat64("lon"),
		}),
		openweather.SetUnits(units),
	)
	if err != nil {
		return err
	}

	// Fetch the current weather conditions
	weather, err := ow.GetOneCallWeather()
	if err != nil {
		return err
	}

	// If JSON, YAML, or TOML output is requested, print the weather in that format
	if viper.GetBool("json") {
		if bytes, err := weather.ToJSON(); err != nil {
			return err
		} else {
			fmt.Println(string(bytes))
			return nil
		}
	} else if viper.GetBool("yaml") {
		if bytes, err := weather.ToYAML(); err != nil {
			return err
		} else {
			fmt.Println(string(bytes))
			return nil
		}
	} else if viper.GetBool("toml") {
		if bytes, err := weather.ToToml(); err != nil {
			return err
		} else {
			fmt.Println(string(bytes))
			return nil
		}
	}

	// Otherwise, print the weather in a human-readable format

	// Get the times
	dt := time.Unix(weather.Current.Dt, 0)
	sunrise := time.Unix(weather.Current.Sunrise, 0)
	sunset := time.Unix(weather.Current.Sunset, 0)

	// Set up the units output
	unit := "°C"
	if weather.Units == "standard" {
		unit = "°K"
	} else if weather.Units == "imperial" {
		unit = "°F"
	}

	// Print the current weather conditions
	fmt.Printf("Current weather for %f, %f as of %s\n", weather.Lat, weather.Lon, dt.Local())
	fmt.Printf("%s %s (%s)\n",
		openweather.Emojis[weather.Current.Weather[0].Icon],
		weather.Current.Weather[0].Main,
		weather.Current.Weather[0].Description,
	)
	fmt.Printf("Temperature: %.1f%s\n", weather.Current.Temp, unit)
	fmt.Printf("Feels like: %.1f%s\n", weather.Current.FeelsLike, unit)
	fmt.Printf("Humidity: %d%%\n", weather.Current.Humidity)
	fmt.Printf("Pressure: %d hPa\n", weather.Current.Pressure)
	fmt.Printf("Due point: %.1f%s\n", weather.Current.DewPoint, unit)
	fmt.Printf("Wind speed: %.1f m/s\n", weather.Current.WindSpeed)
	fmt.Printf("Wind gust: %.1f m/s\n", weather.Current.WindGust)
	fmt.Printf("Wind direction: %d°\n", weather.Current.WindDeg)
	fmt.Printf("Cloudiness: %d%%\n", weather.Current.Clouds)
	fmt.Printf("Rain: %.1f mm\n", weather.Current.Rain)
	fmt.Printf("Snow: %.1f mm\n", weather.Current.Snow)
	fmt.Printf("UV index: %.1f\n", weather.Current.Uvi)
	fmt.Printf("Visibility: %d m\n", weather.Current.Visibility)
	fmt.Printf("Sunrise: %s\n", sunrise.Local())
	fmt.Printf("Sunset: %s\n", sunset.Local())

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
			openweather.Emojis[day.Weather[0].Icon],
			day.Weather[0].Main,
			day.Weather[0].Description,
		)
		fmt.Printf("  High %.1f%s Low %.1f%s\n", day.Temp.Max, unit, day.Temp.Min, unit)
		fmt.Printf("  Morning: %.1f%s (%.1f%s)\n", day.Temp.Morn, unit, day.FeelsLike.Morn, unit)
		fmt.Printf("  Day: %.1f%s (%.1f%s)\n", day.Temp.Day, unit, day.FeelsLike.Day, unit)
		fmt.Printf("  Evening: %.1f%s (%.1f%s)\n", day.Temp.Eve, unit, day.FeelsLike.Eve, unit)
		fmt.Printf("  Night: %.1f%s (%.1f%s)\n", day.Temp.Night, unit, day.FeelsLike.Night, unit)
		fmt.Printf("  Wind speed: %.1f m/s (gust %.1f m/s) from %d°\n", day.WindSpeed, day.WindGust, day.WindDeg)
		fmt.Printf("  Cloudiness: %d%% UV: %.1f\n", day.Clouds, day.Uvi)
		fmt.Printf("  Probability of precipitation: %.1f%%\n", day.Pop)
		fmt.Printf("  Rain: %.1f mm Snow: %.1f mm\n", day.Rain, day.Snow)
		fmt.Printf("  Sunrise (%s) Sunset (%s)\n", sunrise.Local(), sunset.Local())
		fmt.Printf("  Moonrise (%s) Moonset (%s)\n", moonrise.Local(), moonset.Local())

		fmt.Println()
	}

	return nil
}
