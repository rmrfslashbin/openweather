package geocode

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

// Options for the weather query
type Option func(c *Geocoder)

// Geocoder for the weather query
type Geocoder struct {
	log       *zerolog.Logger
	apikey    string
	lang      string
	directUrl *url.URL
	zipUrl    *url.URL
}

// New returns a new Config with the given options
func New(opts ...func(*Geocoder)) (*Geocoder, error) {
	cfg := &Geocoder{}

	// Default to English
	cfg.lang = "en"

	// Construct the direct query URL
	cfg.directUrl = &url.URL{
		// http://api.openweathermap.org/geo/1.0/direct?q={city name},{state code},{country code}&limit={limit}&appid={API key}
		Scheme: "https",
		Host:   "api.openweathermap.org",
		Path:   "/geo/1.0/direct",
	}

	// Construct the direct query URL
	cfg.zipUrl = &url.URL{
		// http://api.openweathermap.org/geo/1.0/zip?zip={zip code},{country code}&appid={API key}
		Scheme: "https",
		Host:   "api.openweathermap.org",
		Path:   "/geo/1.0/zip",
	}

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

	return cfg, nil
}

// WithAPIKey sets the API key
func WithAPIKey(apikey string) Option {
	return func(c *Geocoder) {
		c.apikey = apikey
	}
}

// WithDirectUrl sets the direct URL
func WithDirectUrl(directUrl *url.URL) Option {
	return func(c *Geocoder) {
		c.directUrl = directUrl
	}
}

// WithLogger sets the logger
func WithLogger(log *zerolog.Logger) Option {
	return func(c *Geocoder) {
		c.log = log
	}
}

// WithLanguage sets the language
func WithLanguage(lang string) Option {
	return func(c *Geocoder) {
		c.lang = lang
	}
}

// WithZipUrl sets the zip URL
func WithZipUrl(zipUrl *url.URL) Option {
	return func(c *Geocoder) {
		c.zipUrl = zipUrl
	}
}

func (c *Geocoder) ByCity(city string) (*DirectResponse, error) {
	// Construct the query
	q := c.directUrl.Query()
	q.Set("q", city)
	q.Set("appid", c.apikey)
	q.Set("limit", "5")
	c.directUrl.RawQuery = q.Encode()

	// Make the request
	c.log.Debug().
		Str("url", c.directUrl.String()).
		Msg("getting direct lookup data")

	httpResponse, err := http.Get(c.directUrl.String())
	if err != nil {
		c.log.Error().
			Str("url", c.directUrl.String()).
			Msg("error getting data")
		return nil, err
	}

	// Read the response
	defer httpResponse.Body.Close()
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		c.log.Error().
			Str("url", c.directUrl.String()).
			Msg("error reading data")
		return nil, err
	}

	// Parse the response
	directResponse := &DirectResponse{}
	directResponse.Entities = []*DirectResponseEntity{}
	err = json.Unmarshal(body, &directResponse.Entities)
	if err != nil {
		c.log.Error().
			Str("url", c.directUrl.String()).
			Msg("error unmarshalling data")
		return nil, err
	}

	return directResponse, nil
}

func (c *Geocoder) ByZip(zip string) (*ZipResponse, error) {
	// Construct the query
	q := c.zipUrl.Query()
	q.Set("zip", zip)
	q.Set("appid", c.apikey)
	c.zipUrl.RawQuery = q.Encode()

	// Make the request
	c.log.Debug().
		Str("url", c.zipUrl.String()).
		Msg("getting zip lookup data")

	httpResponse, err := http.Get(c.zipUrl.String())
	if err != nil {
		c.log.Error().
			Str("url", c.zipUrl.String()).
			Msg("error getting data")
		return nil, err
	}

	// Read the response
	defer httpResponse.Body.Close()
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		c.log.Error().
			Str("url", c.zipUrl.String()).
			Msg("error reading data")
		return nil, err
	}

	// Parse the response
	zipResponse := &ZipResponse{}
	err = json.Unmarshal(body, zipResponse)
	if err != nil {
		c.log.Error().
			Str("url", c.zipUrl.String()).
			Msg("error unmarshalling data")
		return nil, err
	}

	return zipResponse, nil
}

// ToJSON returns the zip response as a JSON byte array
func (z *ZipResponse) ToJSON() ([]byte, error) {
	return json.Marshal(z)
}

// ToToml returns the zip response as a TOML byte array
func (z *ZipResponse) ToToml() ([]byte, error) {
	return toml.Marshal(z)
}

// ToYAML returns the zip response as a YAML byte array
func (z *ZipResponse) ToYAML() ([]byte, error) {
	return yaml.Marshal(z)
}

// ToJSON returns the direct response as a JSON byte array
func (d *DirectResponse) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

// ToToml returns the direct response as a TOML byte array
func (d *DirectResponse) ToToml() ([]byte, error) {
	return toml.Marshal(d)
}

// ToYAML returns the direct response as a YAML byte array
func (d *DirectResponse) ToYAML() ([]byte, error) {
	return yaml.Marshal(d)
}
