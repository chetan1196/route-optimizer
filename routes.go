package main

import (
	"errors"
)

var (
	ErrRouteImpossible    = errors.New("route is not possible due to unrealistic speed")
	ErrEmptyLocation      = errors.New("consumer and restaurant locations can't be empty")
	ErrEmptyStartLocation = errors.New("start location can't be empty")
)

type RouteCalculator struct {
	TravelTimeCalculator TravelTimeCalculator
	strategy             RouteCalculationStrategy
}

func (rc *RouteCalculator) SetRoutingAlgo(strategy RouteCalculationStrategy) {
	rc.strategy = strategy
}

// ComputeBestRoute computes the best route for the given orders and start location using the specified strategy.
func (rc *RouteCalculator) ComputeBestRoute(orders []Order, start GeoLocation) (float64, []RouteStep, error) {
	if len(orders) == 0 {
		return 0, nil, errors.New("no orders to process")
	}

	if start.IsNil() {
		return 0, nil, ErrEmptyStartLocation
	}

	// Delegate route calculation to the strategy
	return rc.strategy.CalculateRoute(orders, start)
}
