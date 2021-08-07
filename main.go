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
	"sync"
)

var vehicleNames = map[string]bool {
	"motorcycle": true,
	"car": true,
}

type programArgs struct {
	fname string
}

// Return parking lot and capacity storing data structures
func buildParkingLot(strSlotNums []string) (*map[int64]map[string]servParking.ParkVehicle, *map[int64]map[int64]bool,
	*map[string]types.Vtype, error) {
	parkingLot := make(map[int64]map[string]servParking.ParkVehicle)
	pids := make(map[int64]map[int64]bool)
	vidToVtype := make(map[string]types.Vtype)
	for i, v := range strSlotNums {
		parkingLot[int64(i)] = make(map[string]servParking.ParkVehicle)
		pids[int64(i)] = make(map[int64]bool)
		slotNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, nil, nil, err
		}
		for j := int64(0); j < slotNum; j++ {
			pids[int64(i)][j] = false
		}
	}
	return &parkingLot, &pids, &vidToVtype, nil
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
	parkingLot, pids,vidToVtype, err := buildParkingLot(strSlotNums)
	if err != nil {
		panic(err)
	}
	// Choose concrete behavior for the Parking Service
	var parkService servParking.Service
	// Sync primitive for exclusive locks on parking lot resources - mimicks DB transaction & lock
	mutex := &sync.Mutex{}
	parkService = &impls.ParkingServiceImpl{}
	for scanner.Scan() {
		var timestamp int64
		var vtype types.Vtype
		var vid string
		curLine := scanner.Text()
		ops := strings.Fields(curLine)
		if !isValidRequest(ops) {
			// Report invalid request and skip
			// Todo: Warn log
			continue
		}
		opType, err := servParking.GetOpType(ops[0], true) // Current service uses valet always
		if err != nil {
			// Todo: Warn log
			continue
		}
		if opType == servParking.OpTypeEnterValet || opType == servParking.OpTypeEnter {
			vtype, err = types.ConvertNameToVtype(ops[1])
			if err != nil {
				fmt.Printf("[Warn] invalid client reuqest: %s\n", curLine)
				continue
			}
		}
		if opType == servParking.OpTypeEnterValet || opType == servParking.OpTypeEnter {
			vid = ops[2]
			timestamp, err = strconv.ParseInt(ops[3],10,64)
			if err != nil {
				// Todo: Warn log
				continue
			}
		} else if opType == servParking.OpTypeExitValet || opType == servParking.OpTypeExit {
			vid = ops[1]
			timestamp, err = strconv.ParseInt(ops[2],10,64)
			if err != nil {
				// Todo: Warn log
				continue
			}
		}

		parkServiceReq := servParking.ServiceRequest{
			OpType: opType,
			Vtype: vtype,
			Timestamp: timestamp,
			Vid: vid,
			ParkingLot: parkingLot,
			Pids: pids,
			VidToVtype: vidToVtype,
			Mutex: mutex,
		}
		_, err = parkService.HandleRequest(parkServiceReq)
		if err != nil {
			// Todo: Error log + Metric
		}
	}
	file.Close()
}
