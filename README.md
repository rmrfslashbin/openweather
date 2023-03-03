# openweather
Golang library to interface with OpenWeather Map (dot) org.

## API Key
An API key from https://home.openweathermap.org/api_keys is required to use the library.

## Location
A latitude and longitude pair representing the desired forecast are required to use the library.

## CLI
A CLI is available to interact with the library.

### Usage
- Use the `lookup` command to get the latitude and longitude for a location.
- Use the `current` command to get the current weather conditions for a location. Choose between metric, imperial, or standard units (the default is metric). Then choose an output format. `text` will print the output to the console in a human readable format- add `brief` to show a summary. `json`, `yaml`, and `toml` will print the output to the console in the specified format.


```
$ ./openweather --help                                                                                                                                                                                  1 ↵
Usage: openweather --apikey=STRING <command>

Flags:
  -h, --help               Show context-sensitive help.
      --apikey=STRING      The OpenWeatherMap API key ($APIKEY).
      --loglevel="info"    Set the log level ($LOGLEVEL).

Commands:
  current --apikey=STRING --metric --imperial --standard --lat=FLOAT-64 --lon=FLOAT-64 --json --yaml --toml --text
    Get current weather conditions.

  lookup --apikey=STRING --zip=STRING --city=STRING --json --yaml --toml --text
    Lookup lat/lon data for a location.

Run "openweather <command> --help" for more information on a command.
```
### Exmaple
```
$ ./openweather lookup --apikey ${APIKEY}  --city atlanta,ga,us --text                                                                                                         1 ↵
{"level":"info","app_name":"openweather","log_level":"info","time":"2023-03-03T18:45:52-05:00","message":"starting up"}
Name:     Atlanta
Country:  US
Lat:      33.7489924
Lon:      -84.3902644

$ ./openweather current --apikey ${APIKEY}  --imperial --lat 33.7489924 --lon="-84.3902644" --text --brief --loglevel panic
Current weather as of 2023-03-03 18:46:51 -0500 EST
  🌩️ Thunderstorm (thunderstorm) Temperature: 64.2°F Feels like: 64°F
  Wind speed: 10.4 mph from 260°
  Cloudiness: 100% UV index: 0.0

Friday (2023 March 03)
  🌧️ Rain (moderate rain) High 73.0°F Low 58.9°F with 1.0% chance of precipitation

Saturday (2023 March 04)
  ☀️ Clear (clear sky) High 66.0°F Low 51.2°F with 0.0% chance of precipitation

Sunday (2023 March 05)
  ☀️ Clear (clear sky) High 69.1°F Low 49.0°F with 0.0% chance of precipitation

Monday (2023 March 06)
  ⛅ Clouds (broken clouds) High 73.2°F Low 51.3°F with 0.0% chance of precipitation

Tuesday (2023 March 07)
  🌧️ Rain (light rain) High 70.5°F Low 56.8°F with 0.9% chance of precipitation

Wednesday (2023 March 08)
  🌧️ Rain (light rain) High 53.4°F Low 44.4°F with 0.5% chance of precipitation

Thursday (2023 March 09)
  ⛅ Clouds (overcast clouds) High 52.8°F Low 47.4°F with 0.2% chance of precipitation

Friday (2023 March 10)
  🌧️ Rain (light rain) High 64.9°F Low 51.2°F with 0.3% chance of precipitation

Friday (2023 March 03 18:00) 🌧️ light rain Temp: 64.3°F Wind: 14.9 mph Precip: 1.0%

Friday (2023 March 03 19:00) ⛅ overcast clouds Temp: 64.2°F Wind: 13.3 mph Precip: 0.8%

Friday (2023 March 03 20:00) ⛅ broken clouds Temp: 63.9°F Wind: 12.2 mph Precip: 0.0%

Friday (2023 March 03 21:00) ⛅ broken clouds Temp: 63.0°F Wind: 11.1 mph Precip: 0.0%

Friday (2023 March 03 22:00) 🌥️ scattered clouds Temp: 61.6°F Wind: 10.4 mph Precip: 0.0%

Friday (2023 March 03 23:00) 🌤️ few clouds Temp: 59.4°F Wind: 9.3 mph Precip: 0.0%

Saturday (2023 March 04 00:00) 🌙 clear sky Temp: 56.7°F Wind: 9.2 mph Precip: 0.0%

Saturday (2023 March 04 01:00) 🌙 clear sky Temp: 55.7°F Wind: 8.8 mph Precip: 0.0%

Saturday (2023 March 04 02:00) 🌙 clear sky Temp: 55.3°F Wind: 8.2 mph Precip: 0.0%

Saturday (2023 March 04 03:00) 🌙 clear sky Temp: 54.7°F Wind: 7.4 mph Precip: 0.0%

Saturday (2023 March 04 04:00) 🌙 clear sky Temp: 53.9°F Wind: 7.5 mph Precip: 0.0%

Saturday (2023 March 04 05:00) 🌙 clear sky Temp: 53.1°F Wind: 7.5 mph Precip: 0.0%
```