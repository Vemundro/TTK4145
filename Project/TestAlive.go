package main

import(
	"./alive"
)

const elevators = 2



func main() {
	go alive.CheckElevatorStatus()
	go alive.StartStatusUpdate(elevators)

	alive.BroadcastAlive(1)
}
