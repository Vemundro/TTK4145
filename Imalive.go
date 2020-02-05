package main

import(
	"fmt"
	"net"
	"os"
    "time"
)

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
        os.Exit(1)
    }
}

func main(){
	if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
        os.Exit(1)
    }
    service := os.Args[1]
    udpAddr, err := net.ResolveUDPAddr("udp4", service)
    checkError(err)
    conn, err := net.DialUDP("udp", nil, udpAddr)
    checkError(err)
    lisconn, err := net.ListenUDP("udp", os.Args[2])
    go for {
        _, err = conn.Write([]byte("anything"))
        time.sleep(5 * time.Second())
    }
    go for {
        var buf [512]byte
        n, err := lisconn.Read(buf[0:])
        checkError(err)
        fmt.Println(string(buf[0:n]))
    }
}