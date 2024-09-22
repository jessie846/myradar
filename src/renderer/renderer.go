package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const dotSizeInPxs = 2

type ScreenSize struct {
	Width, Height int32
}

type LatLong struct {
	Latitude, Longitude float64
}

type Renderer struct {
	canvas     *sdl.Renderer
	center     LatLong
	scale      float64
	screenSize ScreenSize
}

func NewRenderer(window *sdl.Window, center LatLong, screenSize ScreenSize, scale float64) (*Renderer, error) {
	canvas, err := window.GetRenderer()
	if err != nil {
		return nil, fmt.Errorf("unable to get renderer: %v", err)
	}

	return &Renderer{
		canvas:     canvas,
		center:     center,
		scale:      scale,
		screenSize: screenSize,
	}, nil
}

func (r *Renderer) Recenter(center LatLong) {
	r.center = center
}

func (r *Renderer) Resize(screenSize ScreenSize) {
	r.screenSize = screenSize
}

func (r *Renderer) Scale(scale float64) {
	r.scale = scale
}

func (r *Renderer) Height() int32 {
	return r.screenSize.Height
}

func (r *Renderer) Width() int32 {
	return r.screenSize.Width
}

func (r *Renderer) DrawDot(point sdl.Point, color sdl.Color) error {
	r.canvas.SetDrawColor(color.R, color.G, color.B, color.A)
	x, y := point.X, point.Y
	rect := sdl.Rect{
		X: x - (dotSizeInPxs / 2),
		Y: y - (dotSizeInPxs / 2),
		W: dotSizeInPxs,
		H: dotSizeInPxs,
	}
	return r.canvas.FillRect(&rect)
}

func (r *Renderer) DrawLines(points []sdl.Point, color sdl.Color) error {
	r.canvas.SetDrawColor(color.R, color.G, color.B, color.A)
	return r.canvas.DrawLines(points)
}

func (r *Renderer) DrawLine(from, to sdl.Point, color sdl.Color) error {
	r.canvas.SetDrawColor(color.R, color.G, color.B, color.A)
	return r.canvas.DrawLine(from.X, from.Y, to.X, to.Y)
}

func (r *Renderer) Draw(
	targetRenderer TargetRenderer,
	mapRenderer MapRenderer,
	mca MCA,
	responseArea ResponseArea,
) error {
	width, height := r.Width(), r.Height()

	r.canvas.SetLogicalSize(width, height)
	r.canvas.SetDrawColor(0, 0, 0, 255) // Black background
	r.canvas.Clear()

	// Render the map, targets, MCA, and response area
	if err := mapRenderer.Render(r); err != nil {
		return fmt.Errorf("failed to render map: %v", err)
	}
	if err := targetRenderer.Render(r); err != nil {
		return fmt.Errorf("failed to render target: %v", err)
	}
	if err := mca.Render(r); err != nil {
		return fmt.Errorf("failed to render MCA: %v", err)
	}
	if err := responseArea.Render(r); err != nil {
		return fmt.Errorf("failed to render response area: %v", err)
	}

	r.canvas.Present()

	return nil
}

func (r *Renderer) RenderSurfaceToCanvas(surface *sdl.Surface, rect sdl.Rect) error {
	texture, err := r.canvas.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("unable to create texture from surface: %v", err)
	}
	defer texture.Destroy()
	return r.canvas.Copy(texture, nil, &rect)
}

func (r *Renderer) SetWindowTitle(window *sdl.Window, title string) error {
	return window.SetTitle(title)
}

func (r *Renderer) ScreenRelativePosition(latLong LatLong) sdl.Point {
	width, height := float64(r.Width()), float64(r.Height())
	dx := latLong.Longitude - r.center.Longitude
	dy := latLong.Latitude - r.center.Latitude
	x := int32(width/2 + math.Round(dx*r.scale))
	y := int32(height/2 - math.Round(dy*r.scale))
	return sdl.Point{X: x, Y: y}
}

func (r *Renderer) PositionFromScreen(point sdl.Point) LatLong {
	width, height := float64(r.Width()), float64(r.Height())
	x, y := float64(point.X), float64(point.Y)
	longitude := (x-width/2)/r.scale + r.center.Longitude
	latitude := (height/2-y)/r.scale + r.center.Latitude
	return LatLong{Latitude: latitude, Longitude: longitude}
}

// Utility functions for rendering text

func RenderText(text string, font *ttf.Font, color sdl.Color) (*sdl.Surface, error) {
	return font.RenderUTF8Blended(text, color)
}

func RenderWrappedText(text string, font *ttf.Font, color sdl.Color, width int) (*sdl.Surface, error) {
	return font.RenderUTF8BlendedWrapped(text, color, width)
}

func RenderTextToSurface(text string, font *ttf.Font, color sdl.Color, surface *sdl.Surface, position sdl.Point) error {
	if len(text) == 0 {
		return nil
	}
	textSurface, err := RenderText(text, font, color)
	if err != nil {
		return fmt.Errorf("unable to render text: %v", err)
	}
	defer textSurface.Free()

	return textSurface.Blit(nil, surface, &sdl.Rect{X: position.X, Y: position.Y})
}

func RenderWrappedTextToSurface(text string, font *ttf.Font, color sdl.Color, width int, surface *sdl.Surface, position sdl.Point) error {
	if len(text) == 0 {
		return nil
	}
	textSurface, err := RenderWrappedText(text, font, color, width)
	if err != nil {
		return fmt.Errorf("unable to render wrapped text: %v", err)
	}
	defer textSurface.Free()

	return textSurface.Blit(nil, surface, &sdl.Rect{X: position.X, Y: position.Y})
}
