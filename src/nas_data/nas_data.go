package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Data structures for deserialization

type MessageCollection struct {
	Messages []Message `xml:"message"`
}

type Message struct {
	Flight NasFlight `xml:"flight"`
}

type Gufi struct {
	Guid string `xml:",chardata"`
}

type NasFlightStatus struct {
	Status string `xml:"fdpsFlightStatus,attr"`
}

type AssignedAltitude struct {
	Simple       *UnitAndValue `xml:"simple"`
	VfrOnTopPlus *UnitAndValue `xml:"vfrOnTopPlus"`
	VfrPlus      *UnitAndValue `xml:"vfrPlus"`
}

func (a *AssignedAltitude) IsOTP() bool {
	return a.VfrOnTopPlus != nil
}

func (a *AssignedAltitude) IsVFR() bool {
	return a.VfrPlus != nil
}

func (a *AssignedAltitude) Value() (float64, error) {
	if a.Simple != nil {
		return strconv.ParseFloat(a.Simple.Value, 64)
	} else if a.VfrPlus != nil {
		return strconv.ParseFloat(a.VfrPlus.Value, 64)
	} else if a.VfrOnTopPlus != nil {
		return strconv.ParseFloat(a.VfrOnTopPlus.Value, 64)
	}
	return 0, errors.New("no altitude value found")
}

type InterimAltitude int

const (
	Missing InterimAltitude = iota
	Unset
	Set
)

type NasRoute struct {
	RouteText *string `xml:"nasRouteText,attr"`
}

type Agreed struct {
	Route NasRoute `xml:"route"`
}

type IcaoModelIdentifier struct {
	Value string `xml:",chardata"`
}

type AircraftType struct {
	IcaoModelIdentifier IcaoModelIdentifier `xml:"icaoModelIdentifier"`
}

func (a *AircraftType) ICAOModelIdentifier() string {
	return a.IcaoModelIdentifier.Value
}

type NasAircraft struct {
	EquipmentSuffix *string      `xml:"equipmentQualifier,attr"`
	AircraftType    AircraftType `xml:"aircraftType"`
}

type RequestedAirspeed struct {
	NasAirspeed UnitAndValue `xml:"nasAirspeed"`
}

func (r *RequestedAirspeed) Value() string {
	return r.NasAirspeed.Value
}

type NasFlight struct {
	Center               string                   `xml:"centre,attr"`
	Timestamp            string                   `xml:"timestamp,attr"`
	Agreed               *Agreed                  `xml:"agreed"`
	AircraftDescription  *NasAircraft             `xml:"aircraftDescription"`
	Arrival              *NasArrival              `xml:"arrival"`
	AssignedAltitude     *AssignedAltitude        `xml:"assignedAltitude"`
	ControllingUnit      *IdentifiedUnitReference `xml:"controllingUnit"`
	Departure            *NasDeparture            `xml:"departure"`
	EnRoute              *NasEnRoute              `xml:"enRoute"`
	FlightIdentification NasFlightIdentification  `xml:"flightIdentification"`
	FlightStatus         *NasFlightStatus         `xml:"flightStatus"`
	Gufi                 Gufi                     `xml:"gufi"`
	InterimAltitude      *NullableUnitAndValue    `xml:"interimAltitude"`
	RequestedAirspeed    *RequestedAirspeed       `xml:"requestedAirspeed"`
}

func (f *NasFlight) Guid() string {
	return f.Gufi.Guid
}

func (f *NasFlight) InterimAltitude() InterimAltitude {
	if f.InterimAltitude != nil {
		if f.InterimAltitude.IsNull() {
			return Unset
		}
		return Set
	}
	return Missing
}

type NasArrival struct {
	ArrivalPoint string `xml:"arrivalPoint,attr"`
}

type IdentifiedUnitReference struct {
	UnitIdentifier   string `xml:"unitIdentifier,attr"`
	SectorIdentifier string `xml:"sectorIdentifier,attr"`
}

type NasDeparture struct {
	DeparturePoint *string `xml:"departurePoint,attr"`
}

type NasHandoff struct {
	Event            *string                  `xml:"event,attr"`
	ReceivingUnit    IdentifiedUnitReference  `xml:"receivingUnit"`
	TransferringUnit *IdentifiedUnitReference `xml:"transferringUnit"`
}

type NasUnitBoundary struct {
	Handoff *NasHandoff `xml:"handoff"`
}

type NasPointout struct {
	OriginatingUnit IdentifiedUnitReference `xml:"originatingUnit"`
	ReceivingUnit   IdentifiedUnitReference `xml:"receivingUnit"`
}

type BeaconCodeAssignment struct {
	CurrentBeaconCode    *string `xml:"currentBeaconCode"`
	PreviousBeaconCode   *string `xml:"previousBeaconCode"`
	ReassignedBeaconCode *string `xml:"reassignedBeaconCode"`
}

type NasEnRoute struct {
	BeaconCodeAssignment *BeaconCodeAssignment `xml:"beaconCodeAssignment"`
	BoundaryCrossings    *NasUnitBoundary      `xml:"boundaryCrossings"`
	Cleared              *NasClearedFlightInfo `xml:"cleared"`
	Pointout             *NasPointout          `xml:"pointout"`
	Position             *NasAircraftPosition  `xml:"position"`
}

type NasClearedFlightInfo struct {
	ClearanceHeading *string `xml:"clearanceHeading,attr"`
	ClearanceSpeed   *string `xml:"clearanceSpeed,attr"`
	ClearanceText    *string `xml:"clearanceText,attr"`
}

type NasFlightIdentification struct {
	CID  string `xml:"computerId,attr"`
	ACID string `xml:"aircraftIdentification,attr"`
}

type NasAircraftPosition struct {
	ActualSpeed    *ActualSpeed   `xml:"actualSpeed"`
	Altitude       *UnitAndValue  `xml:"altitude"`
	Position       *LocationPoint `xml:"position"`
	TargetPosition *Position      `xml:"targetPosition"`
}

func (p *NasAircraftPosition) HasLatLong() bool {
	return p.Position != nil
}

func (p *NasAircraftPosition) Latitude() string {
	return p.Position.Latitude()
}

func (p *NasAircraftPosition) Longitude() string {
	return p.Position.Longitude()
}

func (p *NasAircraftPosition) CurrentAltitude() float64 {
	if p.Altitude != nil {
		altitude, _ := strconv.ParseFloat(p.Altitude.Value, 64)
		return altitude
	}
	return 0
}

func (p *NasAircraftPosition) Speed() float64 {
	if p.ActualSpeed != nil {
		speed, _ := strconv.ParseFloat(p.ActualSpeed.Surveillance.Value, 64)
		return speed
	}
	return 0
}

type ActualSpeed struct {
	Surveillance UnitAndValue `xml:"surveillance"`
}

type UnitAndValue struct {
	Unit  string `xml:"uom,attr"`
	Value string `xml:",chardata"`
}

type NullableUnitAndValue struct {
	Unit  string  `xml:"uom,attr"`
	Value *string `xml:",chardata"`
	Nil   *string `xml:"nil,attr"`
}

func (n *NullableUnitAndValue) IsNull() bool {
	return n.Nil != nil && *n.Nil == "true"
}

type LocationPoint struct {
	Location Position `xml:"location"`
}

type Position struct {
	Pos string `xml:"pos"`
}

func (p *Position) Latitude() string {
	coords := strings.Split(p.Pos, " ")
	return coords[0]
}

func (p *Position) Longitude() string {
	coords := strings.Split(p.Pos, " ")
	return coords[1]
}

// XML Parsing Functions

func ParseFile(filename string) ([]NasFlight, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	fmt.Printf("Parsing file %s...\n", filename)
	return ParseData(string(data))
}

func ParseData(data string) ([]NasFlight, error) {
	var result MessageCollection
	err := xml.Unmarshal([]byte(data), &result)
	if err != nil {
		_ = ioutil.WriteFile("failed-parsing.xml", []byte(data), 0644)
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	var flights []NasFlight
	for _, message := range result.Messages {
		flights = append(flights, message.Flight)
	}
	return flights, nil
}
