package main

import (
	"fmt"
	"log"
	"net"
	"regener/sensor"
	"regener/util"
	"strconv"
	"time"
)

var (
	sAP      = "ap"
	sPOSE    = "pose"
	sCTD     = "ctd"
	sPresure = "presure"
)

var RelayEnable bool = false
var RelayChan chan []byte = make(chan []byte, 100)
var logger = new(util.Logger)
var SQMap map[string]*sensor.Queue

func main() {
	logger.New("regener.log")
	logger.Println("Start SonarGenerator...")
	logger.Println("Load Configuration from cfg.ini ......")
	config := util.LoadCfg("cfg.ini")
	if config == nil {
		logger.Fatal("No valid configuration, exit......")
	}
	//dump all items in config
	config.Dump()

	//logger.Println("Create queue map for sensor data......")
	SQMap = map[string]*sensor.Queue{
		sAP:      sensor.NewQueue(100),
		sPOSE:    sensor.NewQueue(100),
		sCTD:     sensor.NewQueue(100),
		sPresure: sensor.NewQueue(100),
	}
	logger.Println("Start RelayThread.....")
	go RelayThread(config)
	logger.Println("Start ServerThread......")
	SetupServer(config)
}

//relay thread: wait for incoming data and relay to dest addr
func RelayThread(cfg *util.Cfg) {
	server := cfg.RelayIP + ":" + strconv.FormatInt(int64(cfg.RelaySenrPort), 10)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		logger.Fatal(fmt.Sprintf("RelayThread - Fatal error: %s", err.Error()))
	}
	for {
		RelayEnable = false
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Println(fmt.Sprintf("RelayThread - DialTCP error: %s Thread will try connect every 10 seconds", err.Error()))
			select {
			case <-time.After(time.Second * 10):
				continue
			}
		}
		//connected!
		logger.Println("RelayThread - Connect " + server + " successfull")
		RelayEnable = true
		for {
			select {
			case data := <-RelayChan:
				_, err := conn.Write(data)
				if err != nil {
					logger.Println(fmt.Sprintf("RelayThread - relay data error: %s Thread will try relay again", err.Error()))
					break
				}
			}
		}
	}

}
func SetupServer(cfg *util.Cfg) {
	listenaddr := "129.196.34.13:" + strconv.FormatInt(int64(cfg.SensorPort), 10)
	netListen, err := net.Listen("tcp", listenaddr)
	if err != nil {
		logger.Fatal(fmt.Sprintf("ServerThread - Fatal error: %s", err.Error()))

	}
	defer netListen.Close()

	logger.Println("ServerThread - setup " + listenaddr + " successfull- Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		log.Println("ServerThread - "+conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)

	for {

		n, err := conn.Read(buffer)

		if err != nil {
			logger.Println(fmt.Sprintf("Connection - "+conn.RemoteAddr().String(), " connection error: ", err))
			return
		}
		if n > 8 {
			logger.Println(fmt.Sprintf("Connection - %s"+conn.RemoteAddr().String(), string(buffer[:])))
		}

	}

}
