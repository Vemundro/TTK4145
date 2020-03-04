package main

import(
	"./alive"
)

const elevators = 2
const elevatorID = 1



func main() {
	go alive.CheckElevatorStatus()
	go alive.StartStatusUpdate(elevators)
	go alive.BroadcastAlive(elevatorID)

	for {

	}
}
