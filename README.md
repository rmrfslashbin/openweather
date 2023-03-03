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
