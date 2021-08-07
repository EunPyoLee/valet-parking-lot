package main

import (
	"bufio"
	"fmt"
	"os"
	"personal/valet-parking-lot/constant"
	"personal/valet-parking-lot/impls"
	servParking "personal/valet-parking-lot/services/parking"
	"personal/valet-parking-lot/types"
	"strconv"
	"strings"
)

var vehicleNames = map[string]bool {
	"motorcycle": true,
	"car": true,
}

type programArgs struct {
	fname string
}

// Return parking lot and capacity storing data structures
func buildParkingLot(strSlotNums []string) (map[int64]map[string]servParking.ParkVehicle, map[int64]int64, error) {
	parkingLot := make(map[int64]map[string]servParking.ParkVehicle)
	parkingCap := make(map[int64]int64)
	for i, v := range strSlotNums {
		parkingLot[int64(i)] = make(map[string]servParking.ParkVehicle)
		slotNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return parkingLot, parkingCap, err
		}
		parkingCap[int64(i)] = slotNum
	}
	return parkingLot, parkingCap, nil
}

func parseArgs() (*programArgs, error) {
	args := os.Args
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments: expected 2 but received %d\n", len(args))
	}
	return &programArgs{fname: args[1]}, nil
}

func isValidRequest(ops []string) bool {
	if len(ops) != 3 && len(ops) != 4 {
		return false
	}
	if ops[0] == "Enter" && len(ops) == 4 {
		if _, err := types.ConvertNameToVtype(ops[1]); err != nil {
			return false
		}
		return true
	}
	if ops[0] == "Exit" && len(ops) == 3{
		return true
	}
	return false
}

func main() {
	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	fname := args.fname
	file, err := os.Open(fmt.Sprintf("inputs/%s", fname))
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	firstLine := scanner.Text()
	strSlotNums := strings.Fields(firstLine)
	if len(strSlotNums) != constant.VTypeNum {
		panic(fmt.Errorf("invalid number of slot numbers for each vehicle type: expected %d but received %d\n",
			constant.VTypeNum, len(strSlotNums)))
	}
	parkingLot, parkingCap, err := buildParkingLot(strSlotNums)
	if err != nil {
		panic(err)
	}
	// Choose concrete behavior for the Parking Service
	var parkService servParking.Service
	parkService = &impls.ParkingServiceImpl{}
	for scanner.Scan() {
		var timestamp int64
		var vtype types.Vtype
		var vid string
		curLine := scanner.Text()
		ops := strings.Fields(curLine)
		if !isValidRequest(ops) {
			// Report invalid request and skip
			fmt.Printf("Warn: invalid client reuqest: %s\n", curLine)
			continue
		}
		opType, err := servParking.GetOpType(ops[0], true) // Current service uses valet always
		if err != nil {
			fmt.Printf("Warn: invalid client reuqest: %s\n", curLine)
			continue
		}
		if opType == servParking.OpTypeEnterValet || opType == servParking.OpTypeEnter {
			vtype, err = types.ConvertNameToVtype(ops[1])
			if err != nil {
				fmt.Printf("Warn: invalid client reuqest: %s\n", curLine)
				continue
			}
		}
		if opType == servParking.OpTypeEnterValet || opType == servParking.OpTypeEnter {
			vid = ops[2]
			timestamp, err = strconv.ParseInt(ops[3],10,64)
			if err != nil {
				fmt.Printf("Warn: invalid client reuqest: %s\n", curLine)
				continue
			}
		} else if opType == servParking.OpTypeExitValet || opType == servParking.OpTypeExit {
			vid = ops[1]
			timestamp, err = strconv.ParseInt(ops[2],10,64)
			if err != nil {
				fmt.Printf("Warn: invalid client reuqest: %s\n", curLine)
				continue
			}
		}

		parkServiceReq := servParking.ServiceRequest{
			OpType: opType,
			Vtype: vtype,
			Timestamp: timestamp,
			Vid: vid,
			ParkingLot: &parkingLot,
			ParkingCap: &parkingCap,
		}
		fmt.Printf("%+v\n",parkServiceReq)
		parkService.HandleRequest(parkServiceReq)
	}
	file.Close()
}
