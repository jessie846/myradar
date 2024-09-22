package main

import (
	"fmt"
	"os"
	"time"
)

// Constants for time-related operations
const DROP_AFTER = 300 * time.Second // Drop flights after 300 seconds

// FlightList manages the list of flights
type FlightList struct {
	flights       map[string]Flight
	acidToGuidMap map[string]string
	cidToGuidMap  map[string]string
}

// Flight is a placeholder for your actual Flight struct (implement as needed)
type Flight struct {
	acid         string
	cid          string
	lastSeenAt   time.Time
	flightStatus string
}

// NewFlightList creates and initializes a new FlightList
func NewFlightList() *FlightList {
	return &FlightList{
		flights:       make(map[string]Flight),
		acidToGuidMap: make(map[string]string),
		cidToGuidMap:  make(map[string]string),
	}
}

// FindByAcid finds a flight by ACID (Aircraft Identification)
func (fl *FlightList) FindByAcid(acid string) (*Flight, bool) {
	if guid, ok := fl.acidToGuidMap[acid]; ok {
		flight, exists := fl.flights[guid]
		if exists {
			return &flight, true
		}
	}
	return nil, false
}

// FindByCid finds a flight by CID (Computer Identification)
func (fl *FlightList) FindByCid(cid string) (*Flight, bool) {
	if guid, ok := fl.cidToGuidMap[cid]; ok {
		flight, exists := fl.flights[guid]
		if exists {
			return &flight, true
		}
	}
	return nil, false
}

// FindByFlid finds a flight by either CID or ACID (Fallback to either ID)
func (fl *FlightList) FindByFlid(flid string) (*Flight, bool) {
	if flight, ok := fl.FindByCid(flid); ok {
		return flight, ok
	}
	return fl.FindByAcid(flid)
}

// Update updates the list of flights with the provided data
func (fl *FlightList) Update(data string, currentPosition string) {
	// Assuming nasData.ParseData is a function to parse the input data into a list of Flight objects
	nasFlights := ParseData(data)

	for _, nasFlight := range nasFlights {
		// Logging for debugging purposes (if LOG_MESSAGE_TIMESTAMPS is set)
		if os.Getenv("LOG_MESSAGE_TIMESTAMPS") != "" {
			now := time.Now().UTC()
			fmt.Printf("[%s]: Processing flight with GUID: %s\n", now, nasFlight.guid)
		}

		flight, exists := fl.flights[nasFlight.guid]
		if exists {
			// Update the flight record
			flight.updateFromNas(&nasFlight, currentPosition)
			fl.acidToGuidMap[flight.acid] = nasFlight.guid
			fl.cidToGuidMap[flight.cid] = nasFlight.guid
			fl.flights[nasFlight.guid] = flight
		} else {
			// Create a new flight record
			flight := Flight{
				acid:       nasFlight.acid,
				cid:        nasFlight.cid,
				lastSeenAt: time.Now().UTC(),
			}
			fl.acidToGuidMap[flight.acid] = nasFlight.guid
			fl.cidToGuidMap[flight.cid] = nasFlight.guid
			fl.flights[nasFlight.guid] = flight
		}

		// Handle dropped or completed flights
		if nasFlight.flightStatus == "DROPPED" || nasFlight.flightStatus == "COMPLETED" || nasFlight.flightStatus == "CANCELLED" {
			delete(fl.flights, nasFlight.guid)
		}
	}

	fl.pruneDeadFlights()
}

// Prune dead flights that haven't been updated within DROP_AFTER duration
func (fl *FlightList) pruneDeadFlights() {
	for _, guid := range fl.deadFlights() {
		delete(fl.flights, guid)
	}
}

// deadFlights finds flights that are considered "dead" due to inactivity
func (fl *FlightList) deadFlights() []string {
	deadFlights := []string{}
	for guid, flight := range fl.flights {
		if time.Since(flight.lastSeenAt) > DROP_AFTER {
			deadFlights = append(deadFlights, guid)
		}
	}
	return deadFlights
}

// Placeholder function for parsing NAS data
func ParseData(data string) []Flight {
	// This function should implement logic to parse FAA NAS flight data
	// Returning dummy data for the sake of completeness
	return []Flight{
		{acid: "ABC123", cid: "XYZ456", lastSeenAt: time.Now().UTC()},
	}
}

// Update a flight with new NAS data
func (f *Flight) updateFromNas(nasFlight *Flight, currentPosition string) {
	f.acid = nasFlight.acid
	f.cid = nasFlight.cid
	f.lastSeenAt = time.Now().UTC()
	// Add more fields to update as necessary
}

func main() {
	// Example of initializing and updating FlightList
	flightList := NewFlightList()
	flightList.Update("some_flight_data", "current_position")

	// Checking if a flight can be found
	flight, found := flightList.FindByFlid("ABC123")
	if found {
		fmt.Printf("Found flight: %+v\n", flight)
	} else {
		fmt.Println("Flight not found")
	}
}
