package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

// RouteCalculationStrategy defines the interface for route calculation strategies.
type RouteCalculationStrategy interface {
	CalculateRoute(orders []Order, start GeoLocation) (float64, []RouteStep, error)
}
type ConcurrentNaiveStrategy struct {
	travelTimeCalculator TravelTimeCalculator
}

// ConcurrentNaiveStrategy implements the delivery route calculation using a naive algorithm.
// However it uses concurrency to speed up the process.
func (s *ConcurrentNaiveStrategy) CalculateRoute(orders []Order, start GeoLocation) (float64, []RouteStep, error) {
	routeChan := make(chan struct {
		totalTime float64
		steps     []RouteStep
	}, len(orders))

	var wg sync.WaitGroup

	const numWorkers = 5
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for _, perm := range generatePermutations(len(orders)) {
				totalTime, steps, err := s.calculateRouteTime(perm, orders, start)
				if err != nil {
					continue
				}
				routeChan <- struct {
					totalTime float64
					steps     []RouteStep
				}{totalTime, steps}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(routeChan)
	}()

	bestTime := math.MaxFloat64
	var bestSteps []RouteStep
	for result := range routeChan {
		if result.totalTime < bestTime {
			bestTime = result.totalTime
			bestSteps = result.steps
		}
	}

	return bestTime, bestSteps, nil
}

// calculateRouteTime calculates time to deliever all the orders
func (s *ConcurrentNaiveStrategy) calculateRouteTime(orderIndices []int, orders []Order, start GeoLocation) (float64, []RouteStep, error) {
	totalTime := 0.0
	currentLocation := start
	var steps []RouteStep

	for _, index := range orderIndices {
		order := orders[index]

		// Travel to the restaurant
		travelTime := s.travelTimeCalculator.Calculate(currentLocation, order.Restaurant)
		currentLocation = order.Restaurant
		totalTime += travelTime

		if totalTime > maxReasonableTravelTime {
			return -1, nil, ErrRouteImpossible
		}

		steps = append(steps, RouteStep{Action: "Pick up from " + order.RestaurantName, Location: currentLocation})

		totalTime += order.PrepTime

		// Travel to the consumer
		travelTime = s.travelTimeCalculator.Calculate(currentLocation, order.Consumer)
		currentLocation = order.Consumer
		totalTime += travelTime

		steps = append(steps, RouteStep{Action: "Deliver to " + order.ConsumerName, Location: currentLocation})
	}

	return totalTime, steps, nil
}

func generatePermutations(n int) [][]int {
	if n == 1 {
		return [][]int{{0}}
	}
	var result [][]int
	for i := 0; i < n; i++ {
		subPerms := generatePermutations(n - 1)
		for _, perm := range subPerms {
			newPerm := append([]int{i}, perm...)
			result = append(result, newPerm)
		}
	}
	return result
}

// DynamicProgrammingStrategy implements the delivery route calculation using dynamic programming.
type DynamicProgrammingStrategy struct {
	travelTimeCalculator TravelTimeCalculator
}

// CalculateRoute calculates the delivery route using dynamic programming.
func (s *DynamicProgrammingStrategy) CalculateRoute(orders []Order, start GeoLocation) (float64, []RouteStep, error) {
	if len(orders) == 0 {
		return 0, nil, errors.New("no orders provided")
	}

	optimalRoute, totalTime, err := s.findOptimalRoute(orders, start)
	if err != nil {
		return 0, nil, err
	}

	routeSteps := generateRouteSteps(orders, optimalRoute)

	if totalTime > maxReasonableTravelTime {
		return -1, nil, ErrRouteImpossible
	}

	return totalTime, routeSteps, nil
}

// findOptimalRoute finds optimal route to deliever all the orders in minimal time using dp.
func (s *DynamicProgrammingStrategy) findOptimalRoute(orders []Order, start GeoLocation) ([]int, float64, error) {
	n := len(orders)
	if n == 0 {
		return nil, 0, errors.New("no orders provided")
	}

	memo := make(map[[2]int]float64)

	calculateTravelTime := func(from, to GeoLocation) float64 {
		return s.travelTimeCalculator.Calculate(from, to)
	}

	var dp func(int, int) float64
	dp = func(visited int, last int) float64 {
		if visited == (1<<n)-1 {
			return 0
		}
		if val, ok := memo[[2]int{visited, last}]; ok {
			return val
		}
		res := math.MaxFloat64
		for i := 0; i < n; i++ {
			if visited&(1<<i) == 0 {
				newVisited := visited | (1 << i)
				res = math.Min(res, dp(newVisited, i)+calculateTravelTime(orders[last].Consumer, orders[i].Restaurant)+orders[i].PrepTime+calculateTravelTime(orders[i].Restaurant, orders[i].Consumer))
			}
		}
		memo[[2]int{visited, last}] = res
		return res
	}

	// Find optimal route
	last := -1
	bestTime := math.MaxFloat64
	for i := 0; i < n; i++ {
		time := dp(1<<i, i) + calculateTravelTime(start, orders[i].Restaurant) + orders[i].PrepTime + calculateTravelTime(orders[i].Restaurant, orders[i].Consumer)
		if time < bestTime {
			bestTime = time
			last = i
		}

	}

	// Reconstruct the route
	optimalRoute := make([]int, 0, n)
	visited := 1 << last
	for len(optimalRoute) < n {
		if last == -1 {
			break // No next node found
		}
		optimalRoute = append(optimalRoute, last)
		bestNext := -1
		bestNextTime := math.MaxFloat64
		for i := 0; i < n; i++ {
			if visited&(1<<i) == 0 {
				nextVisited := visited | (1 << i)
				nextTime := dp(nextVisited, i) + calculateTravelTime(orders[last].Consumer, orders[i].Restaurant) + orders[i].PrepTime + calculateTravelTime(orders[i].Restaurant, orders[i].Consumer)
				if nextTime < bestNextTime {
					bestNextTime = nextTime
					bestNext = i
				}
			}
		}
		if bestNext == -1 {
			break // No next node found
		}
		visited |= 1 << bestNext
		last = bestNext
	}

	return optimalRoute, bestTime, nil
}

// generateRouteSteps generates route steps from optimal route.
func generateRouteSteps(orders []Order, optimalRoute []int) []RouteStep {
	routeSteps := make([]RouteStep, 0, len(optimalRoute)*2)
	for _, index := range optimalRoute {
		order := orders[index]
		routeSteps = append(routeSteps, RouteStep{
			Action:   fmt.Sprintf("Pick up from %s", order.RestaurantName),
			Location: order.Restaurant,
		})
		routeSteps = append(routeSteps, RouteStep{
			Action:   fmt.Sprintf("Deliver to %s", order.ConsumerName),
			Location: order.Consumer,
		})
	}
	return routeSteps
}
