package main

import (
	"log"
)

func floatPtr(f float64) *float64 {
	return &f
}

func main() {
	orders := []Order{
		{
			ConsumerName:   "Consumer A",
			RestaurantName: "Restaurant A",
			Consumer:       GeoLocation{Lat: floatPtr(12.916), Lon: floatPtr(12.594)},
			Restaurant:     GeoLocation{Lat: floatPtr(12.082), Lon: floatPtr(13.270)},
			PrepTime:       10,
		},
		{
			ConsumerName:   "Consumer B",
			RestaurantName: "Restaurant B",
			Consumer:       GeoLocation{Lat: floatPtr(12.937), Lon: floatPtr(12.894)},
			Restaurant:     GeoLocation{Lat: floatPtr(12.982), Lon: floatPtr(13.670)},
			PrepTime:       8,
		},
	}

	start := GeoLocation{Lat: floatPtr(12.9249), Lon: floatPtr(13.6205)} // Starting point

	err := ValidateOrders(orders)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	builder := NewRouteCalculatorBuilder().
		SetDistanceCalculator(NewHaversineDistanceCalculator()).
		SetSpeed(20)

	routeCalculator, err := builder.Build()
	if err != nil {
		log.Fatalf("Error building routeCalculator: %v", err)
	}

	routeCalculator.SetRoutingAlgo(&ConcurrentNaiveStrategy{routeCalculator.TravelTimeCalculator})

	time, steps, err := routeCalculator.ComputeBestRoute(orders, start)
	if err != nil {
		log.Fatalf("Error computing best route: %v", err)
	}

	log.Printf("Total travel time: %.2f minutes", time)
	for i, step := range steps {
		log.Printf("Step %d: %s at (%.4f, %.4f)", i+1, step.Action, *(step.Location.Lat), *(step.Location.Lon))
	}
}
