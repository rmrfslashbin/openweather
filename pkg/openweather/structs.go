package openweather

import "net/url"

type ErrorResponse struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

// Location for the weather query
type Location struct {
	Lat float64
	Lon float64
}

// Rain holds the rain data
type Rain struct {
	OneH float64 `json:"1h"`
}

// Snow holds the snow data
type Snow struct {
	OneH float64 `json:"1h"`
}

// Weather returns the weather for the given location
type Weather struct {
	Units          string          `json:"units"`
	Lat            float64         `json:"lat"`
	Lon            float64         `json:"lon"`
	Timezone       int             `json:"timezone"`
	TimezoneOffset int             `json:"timezone_offset"`
	Current        *WeatherCurrent `json:"current"`
	Minutely       []struct {
		Dt            int     `json:"dt"`
		Precipitation float64 `json:"precipitation"`
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
	Rain       Rain            `json:"rain"`
	Snow       Snow            `json:"snow"`
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
	Rain       Rain            `json:"rain"`
	Snow       Snow            `json:"snow"`
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
