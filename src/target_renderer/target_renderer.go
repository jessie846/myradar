package target_renderer

import (
	"time"

	"myradar/src/flight"
	"myradar/src/renderer"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	flatTrackSize              = 5
	lineHeight                 = 12
	datablockColor             = 0xE4E400FF
	ldbColor                   = 0x888800FF
	fieldEFlashTime            = 500 * time.Millisecond
	fieldETimeshareTime        = 3 * time.Second
	acceptedHandoffDisplayTime = 3 * time.Minute
	dotSizeInPxs               = 2
)

type TargetRenderer struct {
	charWidth               int32
	currentPosition         flight.Owner
	datablockFont           *ttf.Font
	fieldEFlashTimer        utils.FlipFlopTimer
	fieldETimeshareTimer    utils.FlipFlopTimer
	flightList              []flight.Flight
	renderLDBs              bool
	notYourControlIndicator *sdl.Surface
	pointoutIndicator       *sdl.Surface
}

func NewTargetRenderer(font *ttf.Font, position flight.Owner) (*TargetRenderer, error) {
	charWidth, _, err := font.SizeUTF8("M")
	if err != nil {
		return nil, err
	}

	notYourControlIndicator, err := font.RenderUTF8Blended("R", sdl.Color{R: 228, G: 228, B: 0, A: 255})
	if err != nil {
		return nil, err
	}

	pointoutIndicator, err := font.RenderUTF8Blended("P", sdl.Color{R: 228, G: 228, B: 0, A: 255})
	if err != nil {
		return nil, err
	}

	return &TargetRenderer{
		charWidth:               int32(charWidth),
		currentPosition:         position,
		datablockFont:           font,
		fieldEFlashTimer:        utils.NewFlipFlopTimer(fieldEFlashTime),
		fieldETimeshareTimer:    utils.NewFlipFlopTimer(fieldETimeshareTime),
		flightList:              []flight.Flight{},
		notYourControlIndicator: notYourControlIndicator,
		pointoutIndicator:       pointoutIndicator,
		renderLDBs:              false,
	}, nil
}

func (tr *TargetRenderer) Draw(renderer *renderer.Renderer) error {
	for _, f := range tr.flightList {
		err := tr.drawTarget(&f, renderer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tr *TargetRenderer) drawTarget(flight *flight.Flight, renderer *renderer.Renderer) error {
	if flight.Position != nil {
		point := renderer.ScreenRelativePosition(flight.Position)
		tr.renderTarget(&point, flight, renderer)
		tr.renderDatablock(&point, flight, renderer)
	}
	return nil
}

func (tr *TargetRenderer) renderTarget(point *sdl.Point, flight *flight.Flight, renderer *renderer.Renderer) error {
	if tr.isShowingFDB(flight) {
		// Render flight symbol (square)
		points := []sdl.Point{
			{X: point.X, Y: point.Y - flatTrackSize},
			{X: point.X + flatTrackSize, Y: point.Y},
			{X: point.X, Y: point.Y + flatTrackSize},
			{X: point.X - flatTrackSize, Y: point.Y},
			{X: point.X, Y: point.Y - flatTrackSize},
		}
		renderer.DrawLines(points, sdl.Color{R: 228, G: 228, B: 0, A: 255})

		if flight.IsReducedSeparationEligible() {
			renderer.DrawDot(point, sdl.Color{R: 228, G: 228, B: 0, A: 255})
		} else {
			tr.renderCorrelatedTargetSymbol(point, renderer)
		}
	} else {
		renderer.DrawDot(point, sdl.Color{R: 0, G: 255, B: 0, A: 255})
	}
	return nil
}

func (tr *TargetRenderer) renderCorrelatedTargetSymbol(point *sdl.Point, renderer *renderer.Renderer) error {
	return renderer.DrawLine(
		&sdl.Point{X: point.X - flatTrackSize, Y: point.Y - flatTrackSize},
		&sdl.Point{X: point.X + flatTrackSize, Y: point.Y + flatTrackSize},
		sdl.Color{R: 228, G: 228, B: 0, A: 255},
	)
}

func (tr *TargetRenderer) isShowingFDB(flight *flight.Flight) bool {
	return flight.IsFDBOpen ||
		flight.IsBeingHandedOffTo(&tr.currentPosition) ||
		flight.IsBeingPointedOutTo(&tr.currentPosition) ||
		flight.IsTrackedBy(&tr.currentPosition)
}

func (tr *TargetRenderer) renderDatablock(point *sdl.Point, flight *flight.Flight, renderer *renderer.Renderer) error {
	if tr.isShowingFDB(flight) {
		tr.renderFullDatablock(point, flight, renderer)
	} else if tr.renderLDBs {
		tr.renderLimitedDatablock(point, flight, renderer)
	}
	return nil
}

func (tr *TargetRenderer) renderFullDatablock(point *sdl.Point, flight *flight.Flight, renderer *renderer.Renderer) error {
	// Example of drawing the full datablock
	// Add similar logic to handle flight details and positioning of datablock around target
	// You can use the SDL functions like `Surface.Blit`, and `RenderSurfaceToCanvas`
	return nil
}

func (tr *TargetRenderer) renderLimitedDatablock(point *sdl.Point, flight *flight.Flight, renderer *renderer.Renderer) error {
	// Example of rendering the limited datablock for a flight
	return nil
}

func (tr *TargetRenderer) UpdateFlights(flights []flight.Flight) {
	tr.flightList = flights
}

func (tr *TargetRenderer) UpdateCurrentPosition(position flight.Owner) {
	tr.currentPosition = position
}

func (tr *TargetRenderer) ToggleLDBRendering() {
	tr.renderLDBs = !tr.renderLDBs
}
