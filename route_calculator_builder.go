package main

import "errors"

var (
	ErrInvalidSpeed            = errors.New("invalid speed")
	ErrNoDistanceCalculatorSet = errors.New("distance calculator is not set")
)

type RouteCalculatorBuilder struct {
	distanceCalculator DistanceCalculator
	speed              float64
}

func NewRouteCalculatorBuilder() *RouteCalculatorBuilder {
	return &RouteCalculatorBuilder{}
}

func (b *RouteCalculatorBuilder) SetDistanceCalculator(calculator DistanceCalculator) *RouteCalculatorBuilder {
	b.distanceCalculator = calculator
	return b
}

func (b *RouteCalculatorBuilder) SetSpeed(speed float64) *RouteCalculatorBuilder {
	b.speed = speed
	return b
}

// Build builds the routeCalculator
func (b *RouteCalculatorBuilder) Build() (*RouteCalculator, error) {
	if b.speed <= 0 {
		return nil, ErrInvalidSpeed
	}
	if b.distanceCalculator == nil {
		return nil, ErrNoDistanceCalculatorSet
	}
	routeCalculator := &RouteCalculator{
		TravelTimeCalculator: TravelTimeCalculator{
			DistanceCalculator: b.distanceCalculator,
			Speed:              b.speed,
		},
	}
	// Default Strategy
	routeCalculator.strategy = &ConcurrentNaiveStrategy{routeCalculator.TravelTimeCalculator}
	return routeCalculator, nil
}
