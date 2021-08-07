package valet

import "personal/valet-parking-lot/services/parking"

// Decorator for ParkVehicle
type Valet struct {
	parking.ParkVehicle
}

// Valet charge doubles - Though not instructed in the spec, the end result give expected value
func (v * Valet) GetRate(timestamp int64) float64 {
	return v.ParkVehicle.GetRate(timestamp) * 2
}