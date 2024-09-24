package map_renderer

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/jessie846/myradar/src/custom_map"
	"github.com/jessie846/myradar/src/latlong"
	"github.com/jessie846/myradar/src/renderer"

	geojson "github.com/paulmach/go.geojson"
	"github.com/veandco/go-sdl2/sdl"
)

type MapRenderer struct {
	maps []*custom_map.Map // Ensure this is the correct type
}

// NewMapRenderer creates a new MapRenderer with a list of maps
func NewMapRenderer(maps []*custom_map.Map) *MapRenderer {
	return &MapRenderer{maps: maps}
}

// DrawMap draws the geojson map features using the provided renderer
func (mr *MapRenderer) DrawMap(m *custom_map.Map, r *renderer.Renderer) error {
	fc := m.GeoJson // Use directly without type assertion

	for _, feature := range fc.Features {
		if feature.Geometry != nil {
			if err := mr.DrawGeometry(feature.Geometry, r); err != nil {
				return err
			}
		}
	}
	return nil
}

// DrawGeometry draws different types of geometries using the provided renderer
func (mr *MapRenderer) DrawGeometry(geometry *geojson.Geometry, r *renderer.Renderer) error {
	switch geometry.Type {
	case "GeometryCollection":
		for _, geo := range geometry.Geometries {
			if err := mr.DrawGeometry(geo, r); err != nil {
				return err
			}
		}
	case "LineString":
		if len(geometry.LineString) == 0 {
			return errors.New("LineString geometry is empty")
		}

		var lastPoint *sdl.Point
		for _, coord := range geometry.LineString {
			latLong := latlong.LatLong{
				Latitude:  coord[1],
				Longitude: coord[0],
			}
			point := r.ScreenRelativePosition(renderer.LatLong(latLong)) // Convert to renderer.LatLong

			if lastPoint != nil {
				if err := r.DrawLine(*lastPoint, point, color.Gray{Y: 128}); err != nil {
					return fmt.Errorf("failed to draw line: %v", err)
				}
			}
			lastPoint = &point
		}
	default:
		return nil // Handle unsupported geometry types
	}

	return nil
}

// Render implements the renderable interface for MapRenderer
func (mr *MapRenderer) Render(r *renderer.Renderer) error {
	for _, m := range mr.maps {
		if err := mr.DrawMap(m, r); err != nil {
			return err
		}
	}
	return nil
}
