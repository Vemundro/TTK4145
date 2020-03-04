package alive

import(
	"time"
	"net"
	"fmt"
	"os"
)

//Simple error funtion that will be reused alot for every UDP action
func checkErrorAlive(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
    }
}

//Broadcasts to the other elevators that it's alive
func BroadcastAlive(elevatorID int){
	for{
		msg := []byte{byte(elevatorID)}
		fmt.Println(msg)
		udpAddr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:20345")
	    checkErrorAlive(err)
	    conn, err := net.DialUDP("udp", nil, udpAddr)
	    checkErrorAlive(err)
		_, err = conn.Write([]byte("Hei"))
		checkErrorAlive(err)
		time.Sleep(time.Second)
	}
}

//Creates a UDP connection and sends it to updateStatus
func CheckElevatorStatus(){
	lisAddr, err := net.ResolveUDPAddr("udp4", ":20345")
	checkErrorAlive(err)
	lisconn, err := net.ListenUDP("udp", lisAddr)
	checkErrorAlive(err)

	for {
		updateStatus(lisconn)
	}

}

//Reads the UDP messages and fills channel c with it
func updateStatus(conn *net.UDPConn) {
	c := make(chan int)
	var message []byte

	_, _, err := conn.ReadFromUDP(message)
	checkErrorAlive(err)
	if message != nil {
		intMessage := int(message[0])
		fmt.Println(message)
		c <- intMessage
	}
}

//Initializes the continous status check function and keeps tabs on incoming messages
func StartStatusUpdate(elevators int){
	status := make([]int, elevators-1)
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

//Checks every 10 seconds wether or not an elevator is dead and alerts if one of them is
func statusUpdate(elevators int, rChan chan []int, sChan chan []int){
	status := make([]int, elevators)
	alert := make([]int, elevators)
	alertChan := make(chan []int)

	for index, _ := range alert{
		alert[index] = 0
	}

	for {
		time.Sleep(10 *time.Second)
		status = <- sChan
		fmt.Println(status)
		for index, _ := range status {
			if status[index] == 0 {
				alert[index] += 1
				fmt.Println(alert)
			}
			if status[index] == 1{
				alert[index] == 0
			}
			if alert[index] == 2{
				alertChan <- alert
				fmt.Println(alert)
				time.Sleep(time.Second)
				alert = <- alertChan
			}
		}

	}
}
