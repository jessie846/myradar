package main

// LatLong represents a geographical coordinate with latitude and longitude
type LatLong struct {
	Latitude  float64
	Longitude float64
}

// Clone creates a copy of the LatLong object
func (ll *LatLong) Clone() *LatLong {
	return &LatLong{
		Latitude:  ll.Latitude,
		Longitude: ll.Longitude,
	}
}

func main() {
	// Example usage
	coordinate := LatLong{Latitude: 30.2672, Longitude: -97.7431}
	clone := coordinate.Clone()

	println("Original:", coordinate.Latitude, coordinate.Longitude)
	println("Clone:", clone.Latitude, clone.Longitude)
}
