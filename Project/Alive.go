package alive

import(
	"time"
	"net"
	"fmt"
	"os"
)

/*
To use this package, start two goroutines in your main function. 
One for CheckElevatorStatus() and one for StartStatusUpdate() which needs your timeout delay and number of elevators as arguments
CheckElevatorStatus() is constantly listening for other nodes broadcast, and when it recieves one it will update the Status which is maintained in StartStatusUpdate()
StartStatusUpdate() checks the status array every n seconds and notifies you if one elevator times out
*/

func checkErrorAlive(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
    }
}

func BroadcastAlive(elevatorID int){
	udpAddr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:20345")
    checkErrorAlive(err)
    conn, err := net.DialUDP("udp", nil, udpAddr)
    checkErrorAlive(err)
	_, err = conn.Write([]byte(elevatorID))
	checkErrorAlive(err)
}

func CheckElevatorStatus(){
	lisAddr, err := net.ResolveUDPAddr("udp4", ":20345")
	checkErrorAlive(err)
	lisconn, err := net.ListenUDP("udp", lisAddr)
	checkErrorAlive(err)

	for {
		updateStatus(lisconn)
	}

}

func updateStatus(conn *net.UDPConn) {
	c := make(chan int)
	var message int

	_, addr, err := conn.ReadFromUDP(message)
	checkErrorAlive(err)
	c <- message
}

func StartStatusUpdate(elevators int){
	var status [elevators-1]int
	statusChan := make(chan []int)
	recieveChan := make(chan []int)
	c := make(chan int)
	go statusUpdate(elevators, recieveChan, statusChan)
	for{
		status = <- recieveChan
		select{
		case v:= <- c:
			status[v]=1
		}
		statusChan <- status
	}

}

func statusUpdate(elevators int, rChan chan []int, sChan chan []int){
	var status [elevators-1]int
	var alert [elevators-1]int
	alertChan := make(chan []int)

	for index, element := range alert{
		alert[index] = 0
	}

	for {
		time.Sleep(10 *time.Second)
		status <- sChan
		for index, element := range status {
			if status[index] == 0 {
				alert[index] += 1
			}
			if alert[index] == 2{
				alertChan <- alert
				time.Sleep(time.Second)
				alert <- alertChan
			}
		}
		for index, element := range status {
			status[index] = 0
		}
		rChan <- status
	}
}
