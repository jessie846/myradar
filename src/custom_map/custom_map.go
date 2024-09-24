package custom_map

import (
	"fmt"
	"io/ioutil" // For reading file content
	"os"

	geojson "github.com/paulmach/go.geojson"
)

type Map struct {
	GeoJson *geojson.FeatureCollection
}

// LoadMap loads the GeoJSON map file
func LoadMap(filename string) (*geojson.FeatureCollection, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open map file: %v", err)
	}
	defer file.Close()

	// Read file content into a byte slice
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read map file: %v", err)
	}

	// Unmarshal the byte slice into a FeatureCollection
	fc, err := geojson.UnmarshalFeatureCollection(byteValue)
	if err != nil {
		return nil, fmt.Errorf("failed to parse geojson: %v", err)
	}

	return fc, nil
}
