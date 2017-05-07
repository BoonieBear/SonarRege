package main

import (
	"fmt"
	"log"
	"net"
	"regener/sensor"
	"regener/util"
	"strconv"
	"sync"
	"time"
)

var RelayEnable bool = false
var RelayChan chan []byte = make(chan []byte, 100)
var logger = new(util.Logger)
var SQMap map[uint16]*sensor.Queue
var maplock = new(sync.Mutex)

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
	SQMap = map[uint16]*sensor.Queue{
		sensor.SSId:          sensor.NewQueue(100),
		sensor.BathyId:       sensor.NewQueue(100),
		sensor.SubbottomId:   sensor.NewQueue(100),
		sensor.APHeader:      sensor.NewQueue(100),
		sensor.CompassHeader: sensor.NewQueue(100),
		sensor.CTD6000Header: sensor.NewQueue(100),
		sensor.CTD4500Header: sensor.NewQueue(100),
		sensor.PresureHeader: sensor.NewQueue(100),
	}
	logger.Println("Start Regenarator Thread.....")
	go RegenThread(config)
	logger.Println("Start Relay Thread.....")
	go RelayThread(config)
	logger.Println("Start Server Thread......")
	SetupServer(config)
}

//dispatch the sensor and bsss data
func Dispatcher(recvbuf []byte, queuelock *sync.Mutex) error {
	for {
		if len(recvbuf) < 4 {
			break
		}
		if uint16(util.BytesToUIntBE(16, recvbuf)) == sensor.BsssId {
			if uint16(util.BytesToUIntBE(16, recvbuf[2:])) == sensor.BsssVersion {
				length := util.BytesToUIntBE(32, recvbuf[4:])
				if len(recvbuf) < int(length) {
					//no enough buffer
					break
				} else {
					data := recvbuf[:length]
					recvbuf = append(recvbuf, recvbuf[length:]...)
					return dispatchBsss(data, queuelock)
				}
			}
		}
		if uint16(util.BytesToUIntBE(16, recvbuf)) == sensor.SubbottomId {
			if uint16(util.BytesToUIntBE(16, recvbuf[2:])) == sensor.BsssVersion {
				length := util.BytesToUIntBE(32, recvbuf[4:])
				if len(recvbuf) < int(length) {
					//no enough buffer
					break
				} else {
					data := recvbuf[:length]
					recvbuf = append(recvbuf, recvbuf[length:]...)
					return dispatchSub(data, queuelock)
				}
			}
		}
		if uint16(util.BytesToUIntBE(16, recvbuf)) == sensor.SensorHeadId {
			if uint16(util.BytesToUIntBE(16, recvbuf[2:])) == sensor.SensorVersion {
				length := util.BytesToUIntBE(32, recvbuf[4:])
				if len(recvbuf) < int(length) {
					//no enough buffer
					break
				} else {
					data := recvbuf[:length]
					recvbuf = append(recvbuf, recvbuf[length:]...)
					return dispatchSensor(data, queuelock)
				}
			}
		}
		//shift 2 bytes
		recvbuf = append(recvbuf, recvbuf[2:]...)
	}
	return nil
}

func dispatchBsss(recvbuf []byte, queuelock *sync.Mutex) error {
	bs := &sensor.Bsss{}
	bs.Parse(recvbuf)
	duby := &sensor.DuBathy{}
	duss := &sensor.DuSs{}
	var hasBy, hasSs bool
	for _, v := range bs.Payload {
		if value, ok := v.(sensor.PortByID); ok {
			duby.PortBathy = value
			hasBy = true
		}
		if value, ok := v.(sensor.StarboardByID); ok {
			duby.StarboardBathy = value
			hasBy = true
		}
		if value, ok := v.(sensor.PortSSID); ok {
			duss.PortSs = value
			hasSs = true
		}
		if value, ok := v.(sensor.StarboardSSID); ok {
			duss.StarboardSs = value
			hasSs = true
		}
	}
	if hasSs {
		node := &sensor.Node{
			Time: time.Unix(int64(bs.Dpara.EmitTime1st), int64(bs.Dpara.EmitTime2nd*1000)),
			Data: duss,
		}
		queuelock.Lock()
		if queue, ok := SQMap[sensor.SSId]; ok {

			queue.Push(node)
		}
		queuelock.Unlock()
	}
	if hasBy {
		node := &sensor.Node{
			Time: time.Unix(int64(bs.Dpara.EmitTime1st), int64(bs.Dpara.EmitTime2nd*1000)),
			Data: duby,
		}
		queuelock.Lock()
		if queue, ok := SQMap[sensor.BathyId]; ok {

			queue.Push(node)
		}
		queuelock.Unlock()
	}

	return nil
}
func dispatchSensor(recvbuf []byte, queuelock *sync.Mutex) error {
	totallength := util.BytesToUIntBE(32, recvbuf[4:])
	index := uint64(8)
	if string(recvbuf[8:14]) == "$GPZDA" {
		index += 38
	} else {
		index += 18
	}

	for {
		id := uint16(util.BytesToUIntBE(16, recvbuf[index:]))
		length := util.BytesToUIntBE(16, recvbuf[index+2:])
		switch id {
		case sensor.APHeader:
			ap := &sensor.AP{}
			err := ap.Parse(recvbuf[index:])
			if err == nil {
				node := &sensor.Node{
					Time: ap.Time,
					Data: ap,
				}
				queuelock.Lock()
				if queue, ok := SQMap[sensor.APHeader]; ok {

					queue.Push(node)
				}
				queuelock.Unlock()

			} else {
				log.Println(err)
			}
		case sensor.CompassHeader:
			comp := &sensor.Compass{}
			err := comp.Parse(recvbuf[index:])
			if err == nil {
				node := &sensor.Node{
					Time: comp.Time,
					Data: comp,
				}
				queuelock.Lock()
				if queue, ok := SQMap[sensor.CompassHeader]; ok {

					queue.Push(node)
				}
				queuelock.Unlock()

			} else {
				log.Println(err)
			}
		case sensor.CTD4500Header:
			ctd := &sensor.Ctd4500{}
			err := ctd.Parse(recvbuf[index:])
			if err == nil {
				node := &sensor.Node{
					Time: ctd.Time,
					Data: ctd,
				}
				queuelock.Lock()
				if queue, ok := SQMap[sensor.CTD4500Header]; ok {

					queue.Push(node)
				}
				queuelock.Unlock()

			} else {
				log.Println(err)
			}
		case sensor.CTD6000Header:
			ctd := &sensor.Ctd6000{}
			err := ctd.Parse(recvbuf[index:])
			if err == nil {
				node := &sensor.Node{
					Time: ctd.Time,
					Data: ctd,
				}
				queuelock.Lock()
				if queue, ok := SQMap[sensor.CTD6000Header]; ok {

					queue.Push(node)
				}
				queuelock.Unlock()

			} else {
				log.Println(err)
			}
		case sensor.PresureHeader:
			pre := &sensor.Presure{}
			err := pre.Parse(recvbuf[index:])
			if err == nil {
				node := &sensor.Node{
					Time: pre.Time,
					Data: pre,
				}
				queuelock.Lock()
				if queue, ok := SQMap[sensor.PresureHeader]; ok {

					queue.Push(node)
				}
				queuelock.Unlock()

			} else {
				log.Println(err)
			}
		}
		index += 4
		index += length
		if index == totallength-1 {
			break
		}

	}
	return nil
}
func dispatchSub(recvbuf []byte, queuelock *sync.Mutex) error {
	sub := &sensor.Subbottom{}
	sub.Parse(recvbuf)
	node := &sensor.Node{
		Time: time.Unix(int64(sub.Dpara.EmitTime1st), int64(sub.Dpara.EmitTime2nd*1000)),
		Data: sub,
	}
	queuelock.Lock()
	if queue, ok := SQMap[sensor.SubbottomId]; ok {

		queue.Push(node)
	}
	queuelock.Unlock()
	return nil
}
func 
func RegenThread(cfg *util.Cfg) {

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
	//incoming buffer
	var buffer []byte
	buf := make([]byte, 4096)
	for {

		n, err := conn.Read(buf)

		if err != nil {
			logger.Println(fmt.Sprintf("Connection - "+conn.RemoteAddr().String(), " connection error: ", err))
			return
		}
		if n > 0 {
			buffer = append(buffer, buf[:n]...)
			err = Dispatcher(buffer, maplock)
			if err != nil {
				logger.Println(fmt.Sprintf("Dispatcher error- ", err))
			}
			if RelayEnable {
				RelayChan <- buf[:n]
			}

		}

	}

}
