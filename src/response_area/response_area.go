package response_area

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	WidthInChars     = 30
	MinHeightInLines = 4
	BorderColor      = 0x808080FF // Gray color in hex
	TextColor        = 0xFFFFFFFF // White color in hex
	BorderSize       = 1
	MarginX          = 3
)

type ResponseArea struct {
	content     string
	font        *ttf.Font
	minHeight   int32
	textHeight  int32
	textSurface *sdl.Surface
	width       int32
}

func NewResponseArea(font *ttf.Font) (*ResponseArea, error) {
	charWidth, lineHeight, err := font.SizeUTF8("M")
	if err != nil {
		return nil, err
	}

	minHeight := int32(MinHeightInLines) * lineHeight
	width := int32(WidthInChars)*charWidth + 2*MarginX

	return &ResponseArea{
		font:      font,
		minHeight: minHeight,
		width:     width,
	}, nil
}

func (ra *ResponseArea) Clear() {
	ra.content = ""
	ra.textSurface = nil
}

func (ra *ResponseArea) SetContent(content string, autowrap bool) error {
	ra.content = content
	if len(content) > 0 {
		var textSurface *sdl.Surface
		var err error

		if autowrap {
			textSurface, err = ra.font.RenderUTF8BlendedWrapped(content, sdl.Color{R: 255, G: 255, B: 255, A: 255}, uint32(ra.width))
		} else {
			// Handle newlines properly by using a width greater than the actual width to prevent unwanted wrapping
			textSurface, err = ra.font.RenderUTF8BlendedWrapped(content, sdl.Color{R: 255, G: 255, B: 255, A: 255}, uint32(ra.width+1))
		}

		if err != nil {
			return err
		}

		ra.textHeight = textSurface.H
		ra.textSurface = textSurface
	}

	return nil
}

func (ra *ResponseArea) Render(renderer *sdl.Renderer) error {
	textHeight := int32(math.Max(float64(ra.minHeight), float64(ra.textHeight)))
	totalHeight := textHeight + 2*BorderSize
	totalWidth := ra.width + 2*BorderSize

	// Create a surface with the appropriate width and height
	surface, err := sdl.CreateRGBSurfaceWithFormat(0, totalWidth, totalHeight, 24, sdl.PIXELFORMAT_RGB24)
	if err != nil {
		return err
	}
	defer surface.Free()

	// Fill the surface with the border color (Gray)
	surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: totalWidth, H: totalHeight}, BorderColor)

	// Draw an inner rectangle for the content area (black background)
	surface.FillRect(&sdl.Rect{X: BorderSize, Y: BorderSize, W: ra.width, H: textHeight}, sdl.MapRGB(surface.Format, 0, 0, 0))

	// Blit the text surface onto the response area if it exists
	if ra.textSurface != nil {
		surface.Blit(&sdl.Rect{X: BorderSize + MarginX, Y: BorderSize, W: ra.textSurface.W, H: ra.textSurface.H}, ra.textSurface, nil)
	}

	// Get the position of the response area on the screen
	positionY := int32(renderer.GetOutputHeight()) - totalHeight

	// Create a texture from the surface and render it on the canvas
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	renderer.Copy(texture, nil, &sdl.Rect{X: 250, Y: positionY, W: totalWidth, H: totalHeight})

	return nil
}
