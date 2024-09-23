package windows

import (
	"fmt"
	"time"

	"myradar/src/flight"
	"myradar/src/flight_list"
	"myradar/src/lat_long"
	"myradar/src/mca"
	"myradar/src/renderer"
	"myradar/src/response_area"
	"myradar/src/target_renderer"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowTitle               = "Bomboclaat-radar"
	scaleFactor               = 0.75
	panFactor                 = 0.75
	datablockFontName         = "Inconsolata-VariableFont_wdth,wght.ttf"
	datablockFontSize         = 12
	clickTargetSize   float32 = 5.0
	visibilitySlop    int     = 50
)

// potentiallyVisible checks if a flight is visible within the rendering window.
func potentiallyVisible(r *renderer.Renderer, f *flight.Flight, screenSize *renderer.ScreenSize) bool {
	width, height := screenSize.Width, screenSize.Height
	if f.Position != nil {
		planePos := r.ScreenRelativePosition(f.Position)
		if planePos.X < -visibilitySlop || planePos.Y < -visibilitySlop {
			return false
		}
		adjustedWidth, adjustedHeight := int(width)+visibilitySlop, int(height)+visibilitySlop
		if planePos.X > adjustedWidth || planePos.Y > adjustedHeight {
			return false
		}
		return true
	}
	return false
}

func flightDetails(flight *flight.Flight, currentPosition *flight.Owner) string {
	sector := "??"
	if flight.Owner != nil {
		if flight.Owner.Facility != currentPosition.Facility {
			sector = fmt.Sprintf("%c%s", facilityChar(flight.Owner.Facility), flight.Owner.Sector)
		} else {
			sector = flight.Owner.Sector
		}
	}

	route := flight.Route
	if route == nil {
		route = ""
	}

	acType := flight.AircraftType
	if acType == "" {
		acType = "UNK"
	}

	equipmentSuffix := flight.EquipmentSuffix
	if equipmentSuffix == "" {
		equipmentSuffix = ""
	}

	assignedAltitude := "0"
	if flight.AssignedAltitude != nil {
		if alt := flight.AssignedAltitude.Value(); alt != nil {
			assignedAltitude = fmt.Sprintf("%03d", *alt/100)
		}
	}

	filedCruiseSpeed := "0"
	if speed := flight.FiledCruiseSpeed; speed != nil {
		filedCruiseSpeed = fmt.Sprintf("%.0f", *speed)
	}

	return fmt.Sprintf("%s\n%s %s(%s) %s/%s %s %s %s",
		time.Now().Format("1504"),
		flight.CID,
		flight.ACID,
		sector,
		acType,
		equipmentSuffix,
		flight.AssignedBeaconCode(),
		filedCruiseSpeed,
		assignedAltitude,
		route,
	)
}

func show(
	currentPosition *flight.Owner,
	maps []Map,
	messageReceiver MessageReceiver,
) error {
	flightList := flight_list.NewFlightList()

	window, renderer := initializeSDL() // SDL and font initialization

	defer window.Destroy()
	defer renderer.Destroy()

	scale := 200.0
	var center = lat_long.LatLong{Latitude: 40.2024022, Longitude: -74.4950261}
	eventPump := sdl.GetEventPump()

	var width, height uint32 = 800, 600

	// Initialize Renderers
	mapRenderer := NewMapRenderer(&maps)
	targetRenderer := target_renderer.NewTargetRenderer(datablockFont, *currentPosition)

	mca := mca.NewMCA(&mcaFont)
	responseArea := response_area.NewResponseArea(&responseAreaFont)

	didPan := false

	// Main loop
	for {
		// Update visible flights
		visibleFlights := updateVisibleFlights(&renderer, &flightList, width, height)

		// Message handling
		if message := messageReceiver.Listen(); message != nil {
			flightList.Update(*message, currentPosition)
		}

		for event := eventPump.PollEvent(); event != nil; event = eventPump.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				return nil

			case *sdl.KeyDownEvent:
				switch ev.Keysym.Sym {
				case sdl.K_ESCAPE:
					mca.Clear()
				case sdl.K_l:
					targetRenderer.ToggleLdbRendering()
				}

			case *sdl.MouseMotionEvent:
				if ev.State&sdl.BUTTON_LEFT != 0 {
					didPan = true
					center = panMap(center, ev.XRel, ev.YRel, scale)
					renderer.Recenter(center)
				}

			case *sdl.MouseWheelEvent:
				scale = zoomMap(scale, ev.Y)
				renderer.Scale(scale)

			case *sdl.WindowEvent:
				if ev.Event == sdl.WINDOWEVENT_RESIZED {
					width, height = uint32(ev.Data1), uint32(ev.Data2)
					renderer.Resize(width, height)
				}
			}
		}

		updateAndDrawFlights(visibleFlights, &renderer, targetRenderer, mapRenderer, mca, responseArea)

		sdl.Delay(16)
	}
}

func initializeSDL() (*sdl.Window, *renderer.Renderer) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_RESIZABLE|sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}

	renderer, err := window.GetRenderer()
	if err != nil {
		panic(err)
	}

	return window, renderer
}

func panMap(center lat_long.LatLong, xrel, yrel int32, scale float64) lat_long.LatLong {
	return lat_long.LatLong{
		Latitude:  center.Latitude + float64(yrel)*panFactor/scale,
		Longitude: center.Longitude - float64(xrel)*panFactor/scale,
	}
}

func zoomMap(scale float64, y int32) float64 {
	if y < 0 {
		return scale * scaleFactor
	} else if y > 0 {
		return scale * (1.0 / scaleFactor)
	}
	return scale
}

func updateVisibleFlights(renderer *renderer.Renderer, flightList *flight_list.FlightList, width, height uint32) []string {
	var visibleFlights []string
	for guid, flight := range flightList.Flights {
		if potentiallyVisible(renderer, flight, &renderer.ScreenSize{Width: width, Height: height}) {
			visibleFlights = append(visibleFlights, guid)
		}
	}
	return visibleFlights
}

func updateAndDrawFlights(visibleFlights []string, renderer *renderer.Renderer, targetRenderer *target_renderer.TargetRenderer, mapRenderer *MapRenderer, mca *mca.MCA, responseArea *response_area.ResponseArea) {
	// Update flight rendering list
	var flights []flight.Flight
	for _, guid := range visibleFlights {
		if flight, ok := flightList.Flights[guid]; ok {
			flights = append(flights, *flight)
		}
	}
	targetRenderer.UpdateFlights(flights)
	renderer.Draw(targetRenderer, mapRenderer, mca, responseArea)
}
