package map_renderer

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/paulmach/go.geojson"
	"github.com/veandco/go-sdl2/sdl" // Make sure you added this package using 'go get'
	"github.com/jessie846/myradar/src/latlong"
	"github.com/jessie846/myradar/src/renderer"
	"github.com/jessie846/myradar/src/custom_map"
)

type MapRenderer struct {
	maps []*map.Map
}

// NewMapRenderer creates a new MapRenderer with a list of maps
func NewMapRenderer(maps []*map.Map) *MapRenderer {
	return &MapRenderer{maps: maps}
}

// DrawMap draws the geojson map features using the provided renderer
func (mr *MapRenderer) DrawMap(m *map.Map, r *renderer.Renderer) error {
	switch geo := m.GeoJson.(type) {
	case *geojson.FeatureCollection:
		for _, feature := range geo.Features {
			if feature.Geometry != nil {
				if err := mr.DrawGeometry(feature.Geometry, r); err != nil {
					return err
				}
			}
		}
	case *geojson.Feature:
		if geo.Geometry != nil {
			if err := mr.DrawGeometry(geo.Geometry, r); err != nil {
				return err
			}
		}
	case *geojson.Geometry:
		return mr.DrawGeometry(geo, r)
	}
	return nil
}

// DrawGeometry draws different types of geometries using the provided renderer
func (mr *MapRenderer) DrawGeometry(geometry *geojson.Geometry, r *renderer.Renderer) error {
	switch geometry.Type {
	case geojson.TypeGeometryCollection:
		for _, geo := range geometry.Geometries {
			if err := mr.DrawGeometry(geo, r); err != nil {
				return err
			}
		}
	case geojson.TypeLineString:
		if len(geometry.LineString) == 0 {
			return errors.New("LineString geometry is empty")
		}

		var lastPoint *sdl.Point = nil
		for _, coord := range geometry.LineString {
			latLong := latlong.LatLong{
				Latitude:  coord[1],
				Longitude: coord[0],
			}
			point := r.ScreenRelativePosition(&latLong)

			if lastPoint != nil {
				if err := r.DrawLine(lastPoint, &point, color.Gray{Y: 128}); err != nil {
					return fmt.Errorf("failed to draw line: %v", err)
				}
			}
			lastPoint = &point
		}
	default:
		// Handle unsupported geometry types
		return nil
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
