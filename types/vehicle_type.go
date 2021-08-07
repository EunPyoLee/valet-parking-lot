package types

import "fmt"

type Vtype int64
const (
	VtypeCar Vtype = 0
	VtypeMotorcycle Vtype = 1
)

func GetVtype (vtype int64) (Vtype, error) {
	if vtype < 0 || vtype > 1 {
		return Vtype(-1),fmt.Errorf("invalid Vtype: received [%d]", vtype)
	}
	return Vtype(vtype), nil
}

func ConvertNameToVtype (vname string) (Vtype, error) {
	switch vname {
	case "car":
		return VtypeCar, nil
	case "motorcycle":
		return VtypeMotorcycle, nil
	default:
		return Vtype(-1), fmt.Errorf("invalid vehicle type name: received [%s]", vname)
	}
}


func (v Vtype) GetInt64Val() int64 {
	return int64(v)
}

type Vehicle interface {
	GetVid() string
	GetVtype() Vtype
}

type MotorCycle struct {
	Vid string
	Vtype
}

func (m MotorCycle) GetVid() string {
	return m.Vid
}

func (m MotorCycle) GetVtype() Vtype {
	return m.Vtype
}

type Car struct {
	Vid string
	Vtype
}

func (c Car) GetVid() string {
	return c.Vid
}

func (c Car) GetVtype() Vtype {
	return c.Vtype
}