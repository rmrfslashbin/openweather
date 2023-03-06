package openweather

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
)

func TestGetOneCallWeather(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/3.0/onecall" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		fqpn := filepath.Clean("../../testdata/onecall-v3.0.json")
		fh, err := os.Open(fqpn)
		if err != nil {
			t.Fatalf("failed to open testdata (%s): %v", fqpn, err)
		}
		defer fh.Close()
		oneCallTestData, err := io.ReadAll(fh)
		if err != nil {
			t.Fatalf("failed to read testdata (%s): %v", fqpn, err)
		}
		fmt.Fprint(w, string(oneCallTestData))
	}))
	defer ts.Close()

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	url, _ := url.Parse(ts.URL)
	url.Path = "data/3.0/onecall"
	ow, err := New(
		WithAPIKey("123ABC"),
		WithLocation(&Location{
			Lat: 0.0,
			Lon: 0.0,
		}),
		WithLogger(&log),
		WithUnits(Metric),
		WithRootURL(url),
	)
	if err != nil {
		t.Fatalf("failed to create Openweather instance: %v", err)
	}

	weather, err := ow.GetOneCallWeather()
	if err != nil {
		t.Fatalf("failed to get weather: %v", err)
	}
	if weather.Lat != 33.749 {
		t.Errorf("expected lat to be 33.749, got %f", weather.Lat)
	}
}
