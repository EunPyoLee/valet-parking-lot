package parking

import (
	"fmt"
	"personal/valet-parking-lot/types"
)

type PLErrorCode int64

const (
	PLErrorSuccess PLErrorCode = 0
	PLErrorFail PLErrorCode = 1
	PLErrorInvalidRequest PLErrorCode = 2
	PLErrorServiceUnavailable PLErrorCode = 3
	PLErrorUnknown PLErrorCode = 4
)

func (e PLErrorCode) Error() string {
	switch e {
	case PLErrorSuccess:
		return "Success"
	case PLErrorFail:
		return "Fail"
	case PLErrorInvalidRequest:
		return "Invalid Request"
	case PLErrorServiceUnavailable:
		return "Service Unavailable"
	default:
		return "Unknown Error"
	}
}

func (e PLErrorCode) GetInt64Val() int64 {
	return int64(e)
}

type OpType int64

const (
	OpTypeEnter OpType = 0
	OpTypeExit OpType = 1
	OpTypeEnterValet OpType = 2
	OpTypeExitValet OpType = 3

)

func GetOpType (opTypeStr string, isValet bool) (OpType, error) {
	switch opTypeStr {
	case "Enter":
		if isValet {
			return OpType(2), nil
		} else {
			return OpType(0), nil
		}
	case "Exit":
		if isValet {
			return OpType(3), nil
		} else {
			return OpType(1), nil
		}
	default:
		return OpType(-1), fmt.Errorf("invalid OpType: received [%s]", opTypeStr)
	}
}

func (o OpType) GetInt64Val() int64 {
	return int64(o)
}

type Service interface {
	HandleRequest(req ServiceRequest) (ServiceResponse, error)
}

type ServiceRequest struct {
	OpType OpType // 0 - Enter , 1 - Exit, 2 - ValetEnter, 3 - ValetExit
	Vtype types.Vtype // 0 - Car, 1- Motorcycle
	Timestamp int64
	Vid string // vehicle number
	ParkingLot *map[int64]map[string]ParkVehicle // Parking lot resource mock
	ParkingCap *map[int64]int64 // Parking lot capacity resource mock
}

type ServiceResponse struct {
	errCode PLErrorCode
}

type ParkVehicle interface {
	GetRate(timestamp int64) float64
}

// Abstract Factory for ParkVehicle i-type
func GetParkVehicleFactory(vtype types.Vtype) (ParkVehicleFactory, error) {
	switch vtype {
	case types.VtypeCar:
		return ParkCarFactory{}, nil
	case types.VtypeMotorcycle:
		return ParkMotorcycleFactory{}, nil
	default:
		return nil, fmt.Errorf("invalid Vtype: received [%v]", vtype)
	}
}

type ParkVehicleFactory interface {
	CreateParkVehicle(vid string, vtype types.Vtype, timestamp int64) ParkVehicle
}

type ParkCarFactory struct {}

func (p ParkCarFactory) CreateParkVehicle(vid string, vtype types.Vtype, timestamp int64) ParkVehicle {
	return &ParkCar{
		Vehicle: &types.Car{
			Vid: vid,
			Vtype: vtype,
		},
		Timestamp: timestamp,
	}
}

type ParkMotorcycleFactory struct {}

func (p ParkMotorcycleFactory) CreateParkVehicle(vid string, vtype types.Vtype, timestamp int64) ParkVehicle {
	return &ParkMotorcycle{
		Vehicle: &types.MotorCycle{
			Vid: vid,
			Vtype: vtype,
		},
		Timestamp: timestamp,
	}
}

func ceilUnixTimeHourDiff (begin int64, end int64) int64 {
	diff := end - begin
	hourCeil := diff / 3600
	if diff % 3600 != 0 {
		hourCeil++
	}
	return hourCeil
}

type ParkCar struct {
	types.Vehicle
	Timestamp int64
}

// Flat Car park fee without any extra service(including valet) is 1 / hour
func (p *ParkCar) GetRate(timestamp int64) float64 {
	return 1 * float64(ceilUnixTimeHourDiff(p.Timestamp, timestamp))
}

type ParkMotorcycle struct {
	types.Vehicle
	Timestamp int64
}

// Flat MC park fee without any extra service(including valet) is 0.5 / hour
func (p *ParkMotorcycle) GetRate(timestamp int64) float64 {
	return 1 * float64(ceilUnixTimeHourDiff(p.Timestamp, timestamp))
}