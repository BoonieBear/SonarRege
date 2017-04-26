package main

import (
	"fmt"
	//"regener/oic"
	"net"
	"regener/sensor"
)

var (
	sGPS  = "gps"
	sPOSE = "pose"
)

var RelaySender net.Conn = nil

func main() {
	fmt.Println("Start SonarGenerator...")
	fmt.Println("Load Configuration from cfg.ini ......")
	config := LoadCfg("cfg.ini")
	if config == nil {
		fmt.Println("No valid configuration, exit......")
		return
	}
	//dump all items in config
	config.Dump()

	fmt.Println("Create queue map for sensor data......")
	SQMap := map[string]*sensor.Queue{
		sGPS:  sensor.NewQueue(100),
		sPOSE: sensor.NewQueue(100),
	}
	fmt.Println("Start listen incoming sensor and bsss data......")
	go SetupServer()
}

func SetupServer(*Cfg) {
	netListen, err := net.Listen("tcp", "localhost:1024")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer netListen.Close()

	//Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		//Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)

	for {

		n, err := conn.Read(buffer)

		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		//Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))

	}

}
