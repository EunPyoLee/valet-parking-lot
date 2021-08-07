package impls

import (
	"fmt"
	"personal/valet-parking-lot/constant"
	servParking "personal/valet-parking-lot/services/parking"
	"personal/valet-parking-lot/services/valet"
)

type ParkingServiceImpl struct {
}

// Concrete Implementation of ParkingLotService's HandleRequest Behavior
func (p * ParkingServiceImpl) HandleRequest(req servParking.ServiceRequest) (servParking.ServiceResponse, error) {
	resp := servParking.ServiceResponse{}
	if req.OpType == servParking.OpTypeEnter || req.OpType == servParking.OpTypeEnterValet {
		// Create park vehicle of the requesting vehicle
		parkVehicle, err := createParkVehicle(req)
		if err != nil {
			// Todo: Error log
			resp.ErrCode = servParking.PLErrorFail
			return resp, err
		}
		if req.OpType == servParking.OpTypeEnterValet {
			// Decorate valet wrapper
			parkVehicle = &valet.Valet{
				ParkVehicle: parkVehicle,
			}
		}
		// Enter workflow
		resp, err = enterParkingLot(req, parkVehicle)
		if err != nil {
			// Todo: Error log
			if resp.ErrCode == servParking.PLErrorServiceUnavailable {
				notify(req.OpType, parkVehicle, true) // notfiy "Reject"
			}
		} else{
			notify(req.OpType, parkVehicle, false) // notify "Accept"
		}
		return resp, err
	} else if req.OpType == servParking.OpTypeExit || req.OpType == servParking.OpTypeExitValet {
		// Exit workflow
		return resp, nil
	} else {
		resp.ErrCode = servParking.PLErrorInvalidRequest
		err := fmt.Errorf("[Warn] invalid opType: received : %v\n", req.OpType)
		// Todo: Error log
		return resp, err
	} // Create ParkVehicle
}

func notify(opType servParking.OpType, parkVehicle servParking.ParkVehicle, isReject bool){
	if isReject { // Reject message
		fmt.Println("Reject")
		return
	}
	oneIdxPid := parkVehicle.GetPid() + 1
	if opType == servParking.OpTypeEnter || opType == servParking.OpTypeEnterValet { // Accept  message
		vname, err := parkVehicle.GetVtype().ConvertVtypeToName()
		if err != nil {
			return
		}
		msg := fmt.Sprintf(constant.AcceptFormat, vname, oneIdxPid)
		fmt.Println(msg)
	} else if opType == servParking.OpTypeExit || opType == servParking.OpTypeExitValet {
		return
	} else {
		return
	}
}
func isEnterRequestValid(req servParking.ServiceRequest) bool {
	parkingLot, pids, vidToVtype := req.ParkingLot, req.Pids, req.VidToVtype
	if parkingLot == nil || pids == nil || vidToVtype == nil {
		return false
	}
	if _, exist := (*vidToVtype)[req.Vid]; exist {
		// Duplicate entry without exit history between the two
		return false
	}
	return true
}

func isAvailable(req servParking.ServiceRequest) bool {
	parkingLot, pids := req.ParkingLot, req.Pids
	if len((*pids)[req.Vtype.GetInt64Val()]) == len((*parkingLot)[req.Vtype.GetInt64Val()]) {
		// Full
		return false
	}
	return true
}

// Workflow of Enter operation
func enterParkingLot(req servParking.ServiceRequest, parkVehicle servParking.ParkVehicle) (resp servParking.ServiceResponse, err error){
	defer req.Mutex.Unlock() // RAII-like pattern for unlock guarantee
	// The whole workflow is a critical section
	req.Mutex.Lock()
	// Check request validness
	if !isEnterRequestValid(req) {
		resp.ErrCode = servParking.PLErrorInvalidRequest
		err = fmt.Errorf("[Warn] invalid Enter request: request=[%+v]", req)
		return resp, err
	}
	// Check availability
	if !isAvailable(req) {
		resp.ErrCode = servParking.PLErrorServiceUnavailable
		err = fmt.Errorf("[Warn] Service Unavaialble: request=[%+v]", req)
		return resp, err
	}
	// Enter
	resp, err = parkCar(req, parkVehicle)
	return resp, err
}

// Write requesting vehicle to an empty slot with the minimum slot id
func parkCar(req servParking.ServiceRequest, parkVehicle servParking.ParkVehicle)(resp servParking.ServiceResponse, err error) {
	minPid := int64(-1)
	for i := 0; i < len((*req.Pids)[req.Vtype.GetInt64Val()]); i++ {
		if (*req.Pids)[req.Vtype.GetInt64Val()][int64(i)] == false {
			// Empty
			minPid = int64(i)
			(*req.Pids)[req.Vtype.GetInt64Val()][int64(i)] = true
			break
		}
	}
	if minPid == -1 {
		resp.ErrCode = servParking.PLErrorServiceUnavailable
		err = fmt.Errorf("[Warn] parking lot assignment service is unavaialble\n")
		return resp, err
	}
	parkVehicle.SetPid(minPid)
	(*req.ParkingLot)[req.Vtype.GetInt64Val()][parkVehicle.GetVid()] = parkVehicle
	(*req.VidToVtype)[req.Vid] = req.Vtype
	resp.ErrCode = servParking.PLErrorSuccess
	return resp, nil
}

// Return concrete ParkVehicle created by factory
func createParkVehicle(req servParking.ServiceRequest) (servParking.ParkVehicle, error) {
	parkVehicleFactory, err := servParking.GetParkVehicleFactory(req.Vtype)
	if err != nil {
		// Todo: Error Log
		return nil, err
	}
	return parkVehicleFactory.CreateParkVehicle(req.Vid, req.Vtype, req.Timestamp), nil
}

