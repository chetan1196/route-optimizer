package main

import "fmt"

// GeoLocation represents the geographical location in lat and long
type GeoLocation struct {
	Lat *float64
	Lon *float64
}

func (g *GeoLocation) IsNil() bool {
	return g.Lat == nil || g.Lon == nil
}

func (g *GeoLocation) String() string {
	return fmt.Sprintf("Lat: %v, Lon: %v", *g.Lat, *g.Lon)
}

// Order reprsents the customer order
type Order struct {
	ConsumerName   string
	RestaurantName string
	Consumer       GeoLocation
	Restaurant     GeoLocation
	PrepTime       float64
}

type RouteStep struct {
	Action   string
	Location GeoLocation
}
