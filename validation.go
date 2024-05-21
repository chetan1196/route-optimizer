package main

import "errors"

func ValidateOrders(orders []Order) error {
	if len(orders) == 0 {
		return errors.New("no orders provided")
	}
	for _, order := range orders {
		if order.ConsumerName == "" || order.RestaurantName == "" {
			return errors.New("orders must have consumer and restaurant names")
		}
		if order.Restaurant.IsNil() || order.Consumer.IsNil() {
			return ErrEmptyLocation
		}
		if order.PrepTime <= 0 {
			return errors.New("preparation time must be positive")
		}
	}
	return nil
}
