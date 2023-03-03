package openweather

// Init initializes the weather package
func init() {

	// Set up language list
	langs = make(map[string]struct{})
	availableLangs := []string{
		"af", "al", "ar", "az", "bg", "ca", "cz", "da", "de", "el", "en", "eu",
		"fa", "fi", "fr", "gl", "he", "hi", "hr", "hu", "id", "it", "ja", "kr",
		"la", "lt", "mk", "no", "nl", "pl", "pt", "pt_br", "ro", "ru", "sv", "se",
		"sk", "sl", "sp", "es", "sr", "th", "tr", "ua", "uk", "vi", "zh_cn", "zh_tw", "zu",
	}
	for _, lang := range availableLangs {
		langs[lang] = struct{}{}
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
