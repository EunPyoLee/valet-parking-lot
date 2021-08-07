package impls

import (
	servParking "personal/valet-parking-lot/services/parking"
)

type ParkingServiceImpl struct {
}

// Concrete Implementation of ParkingLotService's HandleRequest Behavior
// Handle entry and Exit operation type
func (p * ParkingServiceImpl) HandleRequest(req servParking.ServiceRequest) (servParking.ServiceResponse, error) {
	resp := servParking.ServiceResponse{}
	return resp, nil
}
