package main

import (
	"math"
)

const (
	maxReasonableTravelTime = 24 * 60 // 24 hours in minutes
)

type DistanceCalculator interface {
	Calculate(start, end GeoLocation) float64
}

type HaversineDistanceCalculator struct{}

func NewHaversineDistanceCalculator() *HaversineDistanceCalculator {
	return &HaversineDistanceCalculator{}
}

func (h *HaversineDistanceCalculator) Calculate(start, end GeoLocation) float64 {
	const earthRadius = 6371 // Radius of the Earth in kilometers
	lat1, lon1 := *(start.Lat)*math.Pi/180, *(start.Lon)*math.Pi/180
	lat2, lon2 := *(end.Lat)*math.Pi/180, *(end.Lon)*math.Pi/180

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

type TravelTimeCalculator struct {
	DistanceCalculator DistanceCalculator
	Speed              float64 // in km/h
}

func (ttc *TravelTimeCalculator) Calculate(start, end GeoLocation) float64 {
	distance := ttc.DistanceCalculator.Calculate(start, end)
	return (distance / ttc.Speed) * 60 // Travel time in minutes
}
