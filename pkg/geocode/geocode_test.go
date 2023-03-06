package geocode

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

func TestByZip(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/geo/1.0/zip" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		fqpn := filepath.Clean("../../testdata/lookupByZip-v1.0.json")
		fh, err := os.Open(fqpn)
		if err != nil {
			t.Fatalf("failed to open testdata (%s): %v", fqpn, err)
		}
		defer fh.Close()
		lookupByZip, err := io.ReadAll(fh)
		if err != nil {
			t.Fatalf("failed to read testdata (%s): %v", fqpn, err)
		}
		fmt.Fprint(w, string(lookupByZip))
	}))
	defer ts.Close()

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	url, _ := url.Parse(ts.URL)
	url.Path = "/geo/1.0/zip"
	gc, err := New(
		WithAPIKey("123ABC"),
		WithLogger(&log),
		WithZipUrl(url),
	)
	if err != nil {
		t.Fatalf("failed to create Geocoder instance: %v", err)
	}

	zipData, err := gc.ByZip("12345")
	if err != nil {
		t.Fatalf("failed to get geocode by zip: %v", err)
	}
	if zipData.Lat != 33.7865 {
		t.Errorf("expected lat to be 33.7865, got %f", zipData.Lat)
	}
}

func TestByCity(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/geo/1.0/direct" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		fqpn := filepath.Clean("../../testdata/lookupByCity-v1.0.json")
		fh, err := os.Open(fqpn)
		if err != nil {
			t.Fatalf("failed to open testdata (%s): %v", fqpn, err)
		}
		defer fh.Close()
		lookupByCity, err := io.ReadAll(fh)
		if err != nil {
			t.Fatalf("failed to read testdata (%s): %v", fqpn, err)
		}
		fmt.Fprint(w, string(lookupByCity))
	}))
	defer ts.Close()

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	url, _ := url.Parse(ts.URL)
	url.Path = "/geo/1.0/direct"
	gc, err := New(
		WithAPIKey("123ABC"),
		WithLogger(&log),
		WithDirectUrl(url),
	)
	if err != nil {
		t.Fatalf("failed to create Geocoder instance: %v", err)
	}

	cityData, err := gc.ByCity("Atlanta")
	if err != nil {
		t.Fatalf("failed to get geocode by city: %v", err)
	}
	if cityData.Entities == nil {
		t.Fatalf("expected city entities to be non-nil")
	}
	if cityData.Entities[0].Name != "Atlanta" {
		t.Errorf("expected lat to be 'Atlanta', got %s", cityData.Entities[0].Name)
	}
}
