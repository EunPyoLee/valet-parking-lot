package parking

import (
	"fmt"
	"time"
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

type Vtype int64

func GetVtype (vtype int64) (Vtype, error) {
	if vtype < 0 || vtype > 1 {
		return Vtype(-1),fmt.Errorf("invalid Vtype: received [%d]", vtype)
	}
	return Vtype(vtype), nil
}

func (v Vtype) GetInt64Val() int64 {
	return int64(v)
}

type ParkingLotService interface {
	HandleRequest(req ParkinglotServiceRequest) (ParkinglotServiceResponse, error)
}

type ParkinglotServiceRequest struct {
	OpType OpType // 0 - Enter , 1 - Exit, 2 - ValetEnter, 3 - ValetExit
	Vtype Vtype // 0 - Car, 1- Motorcycle
	Timestamp time.Time
	Vid string // vehicle number
}

type ParkinglotServiceResponse struct {
	errCode PLErrorCode
}

type ParkVehicle interface {
	GetRate() float64
}

type Vehicle interface {
	GetVid() string
	GetVtype() Vtype
}