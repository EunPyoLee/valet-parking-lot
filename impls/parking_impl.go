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
				notify(req, parkVehicle, true) // notfiy "Reject"
			}
		} else{
			notify(req, parkVehicle, false) // notify "Accept"
		}
		return resp, err
	} else if req.OpType == servParking.OpTypeExit || req.OpType == servParking.OpTypeExitValet {
		// Exit workflow
		resp, parkVehicle, err := exitParkingLot(req)
		if err != nil {
			// Todo: Error log
		} else {
			notify(req, parkVehicle, false) // notfiy "Exit"
		}
		return resp, nil
	} else {
		resp.ErrCode = servParking.PLErrorInvalidRequest
		err := fmt.Errorf("[Warn] invalid opType: received : %v\n", req.OpType)
		// Todo: Error log
		return resp, err
	} // Create ParkVehicle
}

func notify(req servParking.ServiceRequest, parkVehicle servParking.ParkVehicle, isReject bool){
	if isReject { // Reject message
		fmt.Println("Reject")
		return
	}
	oneIdxPid := parkVehicle.GetPid() + 1
	vname, err := parkVehicle.GetVtype().ConvertVtypeToName()
	if err != nil {
		return
	}
	if req.OpType == servParking.OpTypeEnter || req.OpType == servParking.OpTypeEnterValet { // Accept  message
		msg := fmt.Sprintf(constant.AcceptFormat, vname, oneIdxPid)
		fmt.Println(msg)
	} else if req.OpType == servParking.OpTypeExit || req.OpType == servParking.OpTypeExitValet { // Exit message
		msg := fmt.Sprintf(constant.ExitFormat, vname, oneIdxPid, parkVehicle.GetRate(req.Timestamp))
		fmt.Println(msg)
		return
	} else {
		// Todo: Service Unavailable
		return
	}
}

func isExitRequestValid(req servParking.ServiceRequest) bool {
	parkingLot, pids, vidToVtype := req.ParkingLot, req.Pids, req.VidToVtype
	if parkingLot == nil || pids == nil || vidToVtype == nil {
		return false
	}
	if _, exist := (*vidToVtype)[req.Vid]; !exist {
		// No vehicle with this Vid exists in the parking lot
		// Todo: Warn log
		return false
	} else if v, exist := (*parkingLot)[req.Vtype.GetInt64Val()][req.Vid]; !exist ||
		v.GetEnterTimeStamp() > req.Timestamp {
		// Invalid timestamp
		// Todo: Warn log
		return false
	}
	return true
}

func exitCar(req servParking.ServiceRequest) (servParking.ParkVehicle, error) {
	parkingLot, pids, vidToVtype := req.ParkingLot, req.Pids, req.VidToVtype
	parkVehicle := (*parkingLot)[req.Vtype.GetInt64Val()][req.Vid]
	// Unpark and reset the used parking lot id as empty again
	delete((*parkingLot)[req.Vtype.GetInt64Val()], req.Vid)
	(*pids)[req.Vtype.GetInt64Val()][parkVehicle.GetPid()] = false
	delete(*vidToVtype, req.Vid)
	return parkVehicle, nil
}

func exitParkingLot(req servParking.ServiceRequest) (resp servParking.ServiceResponse,
	parkVehicle servParking.ParkVehicle, err error){
	defer req.Mutex.Unlock() // RAII-like pattern for unlock guarantee
	// The whole workflow is a critical section
	req.Mutex.Lock()
	// Check request validness
	if !isExitRequestValid(req) {
		resp.ErrCode = servParking.PLErrorInvalidRequest
		err = fmt.Errorf("[Warn] invalid request: received request=[%+v'", req)
		// Todo: Warn log
		return resp, nil, err
	}
	// Exit
	parkVehicle, err = exitCar(req)
	if err != nil {
		// Todo: Error log
		resp.ErrCode = servParking.PLErrorFail
		err = fmt.Errorf("[Error] request failed: received request=[%+v'", req)
		return resp, nil, err
	} else {
		resp.ErrCode = servParking.PLErrorSuccess
	}
	return resp, parkVehicle, err
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
		// Todo: Warn log
		err = fmt.Errorf("[Warn] invalid request: received request=[%+v'", req)
		return resp, err
	}
	// Check availability
	if !isAvailable(req) {
		resp.ErrCode = servParking.PLErrorServiceUnavailable
		// Todo: Warn log
		err = fmt.Errorf("[Warn] parking service unavailable: received request=[%+v'", req)
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