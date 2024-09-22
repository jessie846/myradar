package crc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// CRCData represents the root structure of the data.
type CRCData struct {
	ID       string          `json:"id"`
	Facility CRCFacilityData `json:"facility"`
}

// CRCFacilityData represents the facility data structure.
type CRCFacilityData struct {
	ERAMConfiguration CRCFacilityERAMConfigurationData `json:"eramConfiguration"`
}

// CRCFacilityERAMConfigurationData holds the ERAM configuration.
type CRCFacilityERAMConfigurationData struct {
	NasID   string          `json:"nasId"`
	GeoMaps []CRCGeoMapData `json:"geoMaps"`
}

// CRCGeoMapData represents the geographic map data.
type CRCGeoMapData struct {
	ID          string                        `json:"id"`
	Name        string                        `json:"name"`
	LabelLine1  string                        `json:"labelLine1"`
	LabelLine2  string                        `json:"labelLine2"`
	FilterMenu  []CRCGeoMapFilterMenuItemData `json:"filterMenu"`
	BcgMenu     []string                      `json:"bcgMenu"`
	VideoMapIDs []string                      `json:"videoMapIds"`
}

// CRCGeoMapFilterMenuItemData represents filter menu item data for geographic maps.
type CRCGeoMapFilterMenuItemData struct {
	ID         string `json:"id"`
	LabelLine1 string `json:"labelLine1"`
	LabelLine2 string `json:"labelLine2"`
}

// LoadData reads and deserializes the JSON data from a file.
func LoadData(filename string) (*CRCData, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read file contents
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Deserialize JSON into CRCData struct
	var data CRCData
	if err := json.Unmarshal(contents, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &data, nil
}
