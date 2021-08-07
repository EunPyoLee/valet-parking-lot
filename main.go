package main

import (
	"bufio"
	"fmt"
	"os"
	"personal/parking/constant"
	"personal/parking/types"
	"strconv"
	"strings"

	//"time"
)

type programArgs struct {
	fname string
}

// Return parking lot and capacity storing data structures
func buildParkingLot(strSlotNums []string) (map[int64]map[string]types.ParkVehicle, map[int64]int64, error) {
	parkingLot := make(map[int64]map[string]types.ParkVehicle)
	capMap := make(map[int64]int64)
	for i, v := range strSlotNums {
		parkingLot[int64(i)] = make(map[string]types.ParkVehicle)
		slotNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return parkingLot, capMap, err
		}
		capMap[int64(i)] = slotNum
	}
	return parkingLot, capMap, nil
}

func parseArgs() (*programArgs, error) {
	args := os.Args
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments: expected 2 but received %d\n", len(args))
	}
	return &programArgs{fname: args[1]}, nil
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
	parkingLot, capMap, err := buildParkingLot(strSlotNums)
	if err != nil {
		panic(err)
	}
	fmt.Println(parkingLot, capMap)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	file.Close()
}
