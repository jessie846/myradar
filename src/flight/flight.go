package flight

import (
	"fmt"
	"time"
)

// DatablockPosition represents the different positions in which data can be displayed
type DatablockPosition string

const (
	N  DatablockPosition = "N"
	NE DatablockPosition = "NE"
	E  DatablockPosition = "E"
	SE DatablockPosition = "SE"
	S  DatablockPosition = "S"
	SW DatablockPosition = "SW"
	W  DatablockPosition = "W"
	NW DatablockPosition = "NW"
)

const DefaultDatablockPosition = SE

// FourthLine represents the fourth line of the datablock, which contains heading, speed, and free text
type FourthLine struct {
	Heading  *string
	Speed    *string
	FreeText *string
}

// Owner represents a facility and sector that controls a flight
type Owner struct {
	Facility string
	Sector   string
}

// FromNas creates an Owner from a NAS controlling unit
func OwnerFromNas(unitIdentifier, sectorIdentifier string) Owner {
	return Owner{
		Facility: unitIdentifier,
		Sector:   sectorIdentifier,
	}
}

// HandoffStatus represents the various states of a handoff
type HandoffStatus string

const (
	Acceptance  HandoffStatus = "ACCEPTANCE"
	Failure     HandoffStatus = "FAILURE"
	Initiation  HandoffStatus = "INITIATION"
	Retraction  HandoffStatus = "RETRACTION"
	TakeControl HandoffStatus = "TAKE_CONTROL"
	Update      HandoffStatus = "UPDATE"
)

// Handoff represents a handoff event for a flight
type Handoff struct {
	Status    *HandoffStatus
	From      *Owner
	To        Owner
	EventTime time.Time
}

// Pointout represents a pointout event for a flight
type Pointout struct {
	From Owner
	To   Owner
}

// Flight represents a flight with associated data like altitude, speed, and ownership
type Flight struct {
	Acid               string
	Cid                string
	Arrival            *string
	Departure          *string
	AssignedAltitude   *float32
	CurrentAltitude    *float32
	InterimAltitude    *float32
	AssignedBeaconCode *string
	CurrentBeaconCode  *string
	Speed              *float32
	Position           *LatLong
	FourthLine         FourthLine
	Handoff            *Handoff
	Pointout           *Pointout
	Owner              *Owner
	IsFDBOpen          bool
	LastSeenAt         time.Time
	AircraftType       *string
	EquipmentSuffix    *string
	FiledCruiseSpeed   *float32
	Route              *string
	DatablockPosition  DatablockPosition
	DatablockLeaderLen uint8
}

// LatLong represents a latitude and longitude position (substitute for missing struct)
type LatLong struct {
	Latitude  float64
	Longitude float64
}

// NewFlight creates a new flight from NAS flight data (dummy NASFlight struct assumed)
func NewFlight(nas NasFlight, currentPosition Owner) Flight {
	var flight Flight
	flight.Acid = nas.FlightIdentification.Acid
	flight.Cid = nas.FlightIdentification.Cid
	flight.Arrival = nas.Arrival
	flight.CurrentAltitude = nas.CurrentAltitude
	flight.InterimAltitude = nas.InterimAltitude
	flight.Speed = nas.Speed
	flight.Route = nas.Route
	flight.AircraftType = nas.AircraftType
	flight.EquipmentSuffix = nas.EquipmentSuffix
	flight.FiledCruiseSpeed = nas.FiledCruiseSpeed
	flight.Position = nas.Position

	// Handle handoff creation
	if nas.Handoff != nil {
		flight.Handoff = &Handoff{
			From:      &nas.Handoff.From,
			To:        nas.Handoff.To,
			Status:    &nas.Handoff.Status,
			EventTime: time.Now(),
		}
	}

	// Handle pointout
	if nas.Pointout != nil {
		flight.Pointout = &Pointout{
			From: nas.Pointout.From,
			To:   nas.Pointout.To,
		}
	}

	flight.Owner = &Owner{
		Facility: currentPosition.Facility,
		Sector:   currentPosition.Sector,
	}
	flight.IsFDBOpen = (flight.Owner.Facility == currentPosition.Facility)

	// Flight timing and status
	flight.LastSeenAt = time.Now()
	flight.DatablockPosition = DefaultDatablockPosition
	flight.DatablockLeaderLen = 1

	return flight
}

// UpdateFromNas updates the flight from new NAS data
func (f *Flight) UpdateFromNas(nas NasFlight, currentPosition Owner) {
	if nas.Position != nil {
		f.Position = nas.Position
	}
	if nas.CurrentAltitude != nil {
		f.CurrentAltitude = nas.CurrentAltitude
	}
	if nas.Speed != nil {
		f.Speed = nas.Speed
	}

	f.LastSeenAt = time.Now()
	if nas.Handoff != nil {
		f.Handoff = &Handoff{
			From:      &nas.Handoff.From,
			To:        nas.Handoff.To,
			Status:    &nas.Handoff.Status,
			EventTime: time.Now(),
		}
	}
}

// HasFourthLine checks if the flight has a fourth line of information
func (f *Flight) HasFourthLine() bool {
	return f.FourthLine.Heading != nil || f.FourthLine.Speed != nil || f.FourthLine.FreeText != nil
}

// IsBeingHandedOffTo checks if the flight is being handed off to a specific owner
func (f *Flight) IsBeingHandedOffTo(owner Owner) bool {
	return f.Handoff != nil && f.Handoff.To == owner
}

// IsBeingPointedOutTo checks if the flight is being pointed out to a specific owner
func (f *Flight) IsBeingPointedOutTo(owner Owner) bool {
	return f.Pointout != nil && f.Pointout.To == owner
}

// IsReducedSeparationEligible checks if the flight is eligible for reduced separation
func (f *Flight) IsReducedSeparationEligible() bool {
	return f.CurrentAltitude != nil && *f.CurrentAltitude <= 23000.0
}

// IsTrackedBy checks if the flight is tracked by the given owner
func (f *Flight) IsTrackedBy(owner Owner) bool {
	return f.Owner != nil && *f.Owner == owner
}

func main() {
	// Placeholder for main function
	fmt.Println("Flight system initialized")
}

// Example of a dummy NasFlight struct for illustrative purposes
type NasFlight struct {
	FlightIdentification struct {
		Acid string
		Cid  string
	}
	Arrival          *string
	CurrentAltitude  *float32
	Speed            *float32
	Route            *string
	AircraftType     *string
	EquipmentSuffix  *string
	FiledCruiseSpeed *float32
	Position         *LatLong
	Handoff          *Handoff
	Pointout         *Pointout
	InterimAltitude  *float32
	AssignedAltitude *float32
}
