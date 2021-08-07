package impls

import (
	"fmt"
	servParking "personal/valet-parking-lot/services/parking"
)

type ParkingServiceImpl struct {
}

// Concrete Implementation of ParkingLotService's HandleRequest Behavior
func (p * ParkingServiceImpl) HandleRequest(req servParking.ServiceRequest) (servParking.ServiceResponse, error) {
	resp := servParking.ServiceResponse{}
	if req.OpType == servParking.OpTypeEnter || req.OpType == servParking.OpTypeEnterValet {
		// Enter workflow
	} else if req.OpType == servParking.OpTypeExit || req.OpType == servParking.OpTypeExitValet {
		// Exit workflow
	} else {
		resp.ErrCode = servParking.PLErrorInvalidRequest
		err := fmt.Errorf("invalid opType: received : %v\n", req.OpType)
		return resp, err
	} // Create ParkVehicle
}
