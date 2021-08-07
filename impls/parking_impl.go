package impls

import (
	servParking "personal/valet-parking-lot/services/parking"
)

type ParkingLotImpl struct {
}

// Concrete Implementation of ParkingLotService's HandleRequest Behavior
// Handle entry and Exit operation type
func (p * ParkingLotImpl) HandleRequest(req servParking.ParkinglotServiceRequest) (servParking.ParkinglotServiceResponse, error) {
}