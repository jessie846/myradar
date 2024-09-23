package main

import (
	"fmt"
	"strings"

	"github.com/jessie846/myradar/src/renderer"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	widthInChars            = 30
	previewAreaHeight       = 15
	feedbackAreaSizeInLines = 4
	borderSize              = 1
	marginX                 = 3
	textColor               = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	borderColor             = sdl.Color{R: 128, G: 128, B: 128, A: 255}
)

type MCA struct {
	charWidth          int32
	currentInput       string
	feedback           *string
	font               *sdl.TTF_Font
	lineHeight         int32
	width              int32
	errorCharSurface   *sdl.Surface
	successCharSurface *sdl.Surface
}

// NewMCA creates a new instance of MCA with the given font.
func NewMCA(font *sdl.TTF_Font) *MCA {
	charWidth, lineHeight, _ := font.SizeText("M")
	errorCharSurface := renderer.RenderText("×", font, sdl.Color{R: 255, G: 0, B: 0, A: 255})
	successCharSurface := renderer.RenderText("✔", font, sdl.Color{R: 0, G: 255, B: 0, A: 255})

	return &MCA{
		charWidth:          int32(charWidth),
		currentInput:       "",
		feedback:           nil,
		font:               font,
		lineHeight:         int32(lineHeight),
		width:              widthInChars*int32(charWidth) + 2*marginX,
		errorCharSurface:   errorCharSurface,
		successCharSurface: successCharSurface,
	}
}

func (m *MCA) Backspace() {
	if len(m.currentInput) > 0 {
		m.currentInput = m.currentInput[:len(m.currentInput)-1]
	}
}

func (m *MCA) Clear() {
	m.ClearInput()
	m.ClearFeedback()
}

func (m *MCA) ClearFeedback() {
	m.feedback = nil
}

func (m *MCA) ClearInput() {
	m.currentInput = ""
}

func (m *MCA) HandleKeyboardInput(text string) {
	m.currentInput += strings.ToUpper(text)
}

func (m *MCA) renderErrorMessage(message string, surface *sdl.Surface) {
	m.renderTextWithSymbol(message, m.errorCharSurface, surface)
}

func (m *MCA) renderSuccessMessage(message string, surface *sdl.Surface) {
	m.renderTextWithSymbol(message, m.successCharSurface, surface)
}

func (m *MCA) renderTextWithSymbol(message string, symbolSurface *sdl.Surface, outputSurface *sdl.Surface) {
	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return
	}

	lineOffset := int32(previewAreaHeight + 2*borderSize)
	marginX := int32(marginX)

	symbolSurface.Blit(nil, outputSurface, &sdl.Rect{X: marginX, Y: lineOffset})
	renderer.RenderTextToSurface(lines[0], m.font, textColor, outputSurface, sdl.Point{X: 2*marginX + m.charWidth, Y: lineOffset})

	for i, line := range lines[1:] {
		lineOffset += m.lineHeight
		renderer.RenderTextToSurface(line, m.font, textColor, outputSurface, sdl.Point{X: marginX, Y: lineOffset + int32(i)*m.lineHeight})
	}
}

func (m *MCA) SetFeedback(text string) {
	m.feedback = &text
}

func (m *MCA) SetErrorFeedback(text string) {
	m.feedback = &text
}

func (m *MCA) Value() string {
	return m.currentInput
}

func (m *MCA) Render(r *renderer.Renderer) error {
	feedbackAreaHeight := feedbackAreaSizeInLines * m.lineHeight
	totalHeight := previewAreaHeight + feedbackAreaHeight + 3*borderSize
	totalWidth := m.width + 2*borderSize

	surface, err := sdl.CreateRGBSurface(0, totalWidth, totalHeight, 24, 0, 0, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to create surface: %v", err)
	}
	defer surface.Free()

	surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: totalWidth, H: totalHeight}, borderColor.Uint32())
	borderSizeI32 := int32(borderSize)

	// Draw inner rectangle for the preview area
	surface.FillRect(&sdl.Rect{X: borderSizeI32, Y: borderSizeI32, W: m.width, H: previewAreaHeight}, sdl.Color{R: 0, G: 0, B: 0, A: 255}.Uint32())

	// Draw another inner rectangle for the feedback area
	surface.FillRect(&sdl.Rect{X: borderSizeI32, Y: previewAreaHeight + 2*borderSizeI32, W: m.width, H: feedbackAreaHeight}, sdl.Color{R: 0, G: 0, B: 0, A: 255}.Uint32())

	// Append a cursor to the input
	text := m.currentInput + "_"
	renderer.RenderTextToSurface(text, m.font, textColor, surface, sdl.Point{X: marginX, Y: borderSizeI32})

	// Render feedback
	if m.feedback != nil {
		m.renderSuccessMessage(*m.feedback, surface)
	}

	// Render on the main canvas
	position := int32(r.Height()) - totalHeight
	r.RenderSurfaceToCanvas(surface, &sdl.Rect{X: 0, Y: position, W: m.width, H: totalHeight})

	return nil
}
