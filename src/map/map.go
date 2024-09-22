package main

import (
	"fmt"
	"io/ioutil"
	"os"

	geojson "github.com/paulmach/go.geojson"
)

// Map structure that holds GeoJSON data.
type Map struct {
	Geojson *geojson.FeatureCollection
}

// FromFile reads a GeoJSON file and returns a Map struct.
func FromFile(filename string) (*Map, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Read the file contents
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Parse the contents into GeoJSON format
	geojsonData, err := geojson.UnmarshalFeatureCollection(contents)
	if err != nil {
		return nil, fmt.Errorf("failed to parse geojson from file %s: %w", filename, err)
	}

	// Return the Map struct with the parsed GeoJSON
	return &Map{Geojson: geojsonData}, nil
}
