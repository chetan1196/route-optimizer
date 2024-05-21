package main

import (
	"testing"
)

func TestHaversineDistanceCalculator(t *testing.T) {
	calculator := NewHaversineDistanceCalculator()

	loc1 := GeoLocation{Lat: floatPtr(12.916), Lon: floatPtr(12.594)}
	loc2 := GeoLocation{Lat: floatPtr(12.082), Lon: floatPtr(13.270)}

	expectedDistance := 118.26 // Approximate distance in kilometers
	distance := calculator.Calculate(loc1, loc2)

	if distance < expectedDistance-1 || distance > expectedDistance+1 {
		t.Errorf("Expected distance to be around %.2f km, but got %.2f km", expectedDistance, distance)
	}
}

func TestTravelTimeCalculator(t *testing.T) {
	distanceCalculator := NewHaversineDistanceCalculator()

	builder := NewRouteCalculatorBuilder().
		SetDistanceCalculator(distanceCalculator).
		SetSpeed(20)

	routeCalculator, _ := builder.Build()

	loc1 := GeoLocation{Lat: floatPtr(12.916), Lon: floatPtr(12.594)}
	loc2 := GeoLocation{Lat: floatPtr(12.082), Lon: floatPtr(13.270)}

	expectedTravelTime := 354.78 // in minutes
	travelTime := routeCalculator.TravelTimeCalculator.Calculate(loc1, loc2)

	if travelTime < expectedTravelTime-10 || travelTime > expectedTravelTime+10 {
		t.Errorf("Expected travel time to be around %.2f minutes, but got %.2f minutes", expectedTravelTime, travelTime)
	}
}

func TestRouteCalculator(t *testing.T) {
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

	start := GeoLocation{Lat: floatPtr(12.9249), Lon: floatPtr(12.6205)} // Starting point

	distanceCalculator := NewHaversineDistanceCalculator()
	builder := NewRouteCalculatorBuilder().
		SetDistanceCalculator(distanceCalculator).
		SetSpeed(20)

	routeCalculator, _ := builder.Build()

	time, steps, err := routeCalculator.ComputeBestRoute(orders, start)
	if err != nil {
		t.Fatalf("Error computing best route: %v", err)
	}

	if time <= 0 {
		t.Errorf("Expected positive travel time, but got %.2f", time)
	}

	expectedSteps := [][]string{
		{
			"Pick up from Restaurant B",
			"Deliver to Consumer B",
			"Pick up from Restaurant A",
			"Deliver to Consumer A",
		},
	}

	actualSteps := []string{}
	for _, step := range steps {
		actualSteps = append(actualSteps, step.Action)
	}

	validRoute := false
	for _, expSteps := range expectedSteps {
		if equalSteps(expSteps, actualSteps) {
			validRoute = true
			break
		}
	}

	if !validRoute {
		t.Errorf("Unexpected route steps: %v", actualSteps)
	}
}

func equalSteps(expected, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return false
		}
	}
	return true
}

func TestEmptyOrderList(t *testing.T) {
	orders := []Order{}

	start := GeoLocation{Lat: floatPtr(12.9249), Lon: floatPtr(13.6205)} // Starting point

	distanceCalculator := NewHaversineDistanceCalculator()

	builder := NewRouteCalculatorBuilder().
		SetDistanceCalculator(distanceCalculator).
		SetSpeed(20)

	routeCalculator, _ := builder.Build()

	_, _, err := routeCalculator.ComputeBestRoute(orders, start)
	if err == nil {
		t.Error("Expected error for empty order list, but got nil")
	}
}

func TestInvalidOrders(t *testing.T) {
	// Orders with missing geo-locations
	orders := []Order{
		{
			ConsumerName:   "Consumer A",
			RestaurantName: "Restaurant A",
			PrepTime:       10,
		},
	}

	err := ValidateOrders(orders)
	if err == nil {
		t.Error("Expected error for orders with missing geo-locations, but got nil")
	}

	// Orders with negative preparation time
	orders = []Order{
		{
			ConsumerName:   "Consumer A",
			RestaurantName: "Restaurant A",
			Consumer:       GeoLocation{Lat: floatPtr(12.916), Lon: floatPtr(12.594)},
			Restaurant:     GeoLocation{Lat: floatPtr(12.082), Lon: floatPtr(13.270)},
			PrepTime:       -5,
		},
	}

	err = ValidateOrders(orders)
	if err == nil {
		t.Error("Expected error for orders with negative preparation time, but got nil")
	}
}

func TestStartingLocationNotProvided(t *testing.T) {
	orders := []Order{
		{
			ConsumerName:   "Consumer A",
			RestaurantName: "Restaurant A",
			Consumer:       GeoLocation{Lat: floatPtr(12.916), Lon: floatPtr(12.594)},
			Restaurant:     GeoLocation{Lat: floatPtr(12.082), Lon: floatPtr(13.270)},
			PrepTime:       10,
		},
	}

	var start GeoLocation // Starting point not provided

	distanceCalculator := NewHaversineDistanceCalculator()

	builder := NewRouteCalculatorBuilder().
		SetDistanceCalculator(distanceCalculator).
		SetSpeed(20)

	routeCalculator, _ := builder.Build()

	_, _, err := routeCalculator.ComputeBestRoute(orders, start)
	if err == nil {
		t.Error("Expected error for starting location not provided, but got nil")
	}
}

func TestZeroOrNegativeSpeed(t *testing.T) {
	distanceCalculator := NewHaversineDistanceCalculator()
	builder := NewRouteCalculatorBuilder().
		SetDistanceCalculator(distanceCalculator).
		SetSpeed(0)

	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error for zero speed, but got nil")
	}

	builder.SetSpeed(-10) // Negative speed
	_, err = builder.Build()
	if err == nil {
		t.Error("Expected error for negative speed, but got nil")
	}
}

func TestNilDistanceCalculator(t *testing.T) {
	var distanceCalculator DistanceCalculator // Nil distance calculator
	builder := NewRouteCalculatorBuilder().
		SetDistanceCalculator(distanceCalculator).
		SetSpeed(20)

	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error for nil distance calculator, but got nil")
	}
}

type mockDistanceCalculator struct{}

func (m *mockDistanceCalculator) Calculate(start, end GeoLocation) float64 {
	// Return a very large distance to simulate an impossible route
	return 1000000.0 // 1000 km
}

func TestRouteCalculator_NoRoutePossible(t *testing.T) {
	orders := []Order{
		{
			ConsumerName:   "Consumer A",
			RestaurantName: "Restaurant A",
			Consumer:       GeoLocation{Lat: floatPtr(12.916), Lon: floatPtr(12.594)},
			Restaurant:     GeoLocation{Lat: floatPtr(12.082), Lon: floatPtr(13.270)},
			PrepTime:       10,
		},
	}

	distanceCalculator := &mockDistanceCalculator{}
	speed := 10.0 // km/h

	builder := NewRouteCalculatorBuilder().
		SetDistanceCalculator(distanceCalculator).
		SetSpeed(speed)

	routeCalculator, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build routeCalculator: %v", err)
	}

	// Set up test case where distance is too great for the given speed
	start := GeoLocation{Lat: floatPtr(0), Lon: floatPtr(0)}
	_, steps, _ := routeCalculator.ComputeBestRoute(orders, start)
	if steps != nil {
		t.Errorf("Expected nil steps, but got ['%v']", steps)
	}
}
