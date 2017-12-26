package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regener/oic"
	"regener/sensor"
	"regener/util"
	"strconv"
	"sync"
	"time"
)

const (
	DIALTIMEOUT = 5
)

//RelayEnable relay socket enable
var RelayEnable = false

//RegenbEnable relay socket enable
var RegenbEnable = false

//RelayChan chan between socket conn and relay thread
var RelayChan = make(chan []byte, 100)
var logger = new(util.Logger)
var buffer []byte

//SQMap external data queue map
var SQMap map[uint16]*sensor.Queue
var maplock = new(sync.Mutex)
var trace = &util.Tracefile{}

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

	// Capture ctrl-c or other interruptions then clean up the global lock.
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go func(cc <-chan os.Signal) {
		s := <-cc
		RegenbEnable = false
		RelayEnable = false
		// Exiting with the expected exit codes when we can.
		if s == os.Interrupt {
			os.Exit(130)
		} else if s == os.Kill {
			os.Exit(137)
		} else {
			os.Exit(1)
		}
	}(ch)
	logger.Println("Start Server Thread......")
	SetupServer(config)
}

//Dispatcher dispatch the sensor and bsss data
func Dispatcher(buf []byte, queuelock *sync.Mutex) error {
	buffer = append(buffer, buf[:]...)
	for {
		if len(buffer) < 4 {
			break
		}
		if uint16(util.BytesToUIntLE(16, buffer)) == sensor.BsssId {
			if uint16(util.BytesToUIntLE(16, buffer[2:])) == sensor.BsssVersion {
				length := util.BytesToUIntLE(32, buffer[4:])
				if len(buffer) < int(length) {
					//no enough buffer
					break
				} else {
					data := make([]byte, length)
					copy(data, buffer[:length])
					buffer = append(buffer[:0], buffer[length:]...)
					return DispatchBsss(data, queuelock)
				}
			}
		}
		if uint16(util.BytesToUIntBE(16, buffer)) == sensor.SubbottomId {
			if uint16(util.BytesToUIntBE(16, buffer[2:])) == sensor.BsssVersion {
				length := util.BytesToUIntBE(32, buffer[4:])
				if len(buffer) < int(length) {
					//no enough buffer
					break
				} else {
					data := make([]byte, length)
					copy(data, buffer[:length])
					buffer = append(buffer[:0], buffer[length:]...)
					return DispatchSub(data, queuelock)
				}
			}
		}
		if uint16(util.BytesToUIntLE(16, buffer)) == sensor.SensorHeadId {
			if uint16(util.BytesToUIntLE(16, buffer[2:])) == sensor.SensorVersion {
				length := util.BytesToUIntLE(32, buffer[4:])
				if len(buffer) < int(length) {
					//no enough buffer
					break
				} else {
					data := make([]byte, length)
					copy(data, buffer[:length])
					buffer = append(buffer[:0], buffer[length:]...)
					return DispatchSensor(data, queuelock)
				}
			}
		}
		//shift 2 bytes
		buffer = append(buffer[:0], buffer[2:]...)
	}
	return nil
}

func DispatchBsss(recvbuf []byte, queuelock *sync.Mutex) error {
	bs := &sensor.Bsss{}
	bs.Parse(recvbuf)
	duby := &sensor.DuBathy{}
	duss := &sensor.DuSs{}
	var hasBy, hasSs bool
	for _, v := range bs.Payload {
		if value, ok := v.(*sensor.SingelBathy); ok {
			bathy := value
			if bathy.ID == uint32(sensor.PortByID) {
				duby.PortBathy = value
			} else {
				duby.StarboardBathy = value
			}
			hasBy = true
		}

		if value, ok := v.(*sensor.Ss); ok {
			ss := value
			if ss.ID == uint32(sensor.PortSSID) {
				duss.PortSs = value
			} else {
				duss.StarboardSs = value
			}
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
func DispatchSensor(recvbuf []byte, queuelock *sync.Mutex) error {
	totallength := util.BytesToUIntLE(32, recvbuf[4:])
	fmt.Printf("totallength = %d\n", totallength)
	index := uint64(8)
	if string(recvbuf[8:14]) == "$GPZDA" {
		index += 38
	} else {
		index += 18
	}

	for {
		id := uint16(util.BytesToUIntBE(16, recvbuf[index:]))
		fmt.Printf("id = %x\n", id)
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
					fmt.Println("push AP")
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
					fmt.Println("push Compass")
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
					fmt.Println("push Ctd4500")
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
					fmt.Println("push Ctd6000")
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
					fmt.Println("push Presure")
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
func DispatchSub(recvbuf []byte, queuelock *sync.Mutex) error {
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

//MergeOICBathy add sensor data to bathy structure
func MergeOICBathy(by *oic.Bathy, data *sensor.MixData) []byte {
	if data.Ap != nil {
		by.Header.NavFixLatitude = data.Ap.Lat
		by.Header.NavFixLongtitude = data.Ap.Lng
	}
	if data.Comp != nil {
		by.Header.VesselHeading = float32(data.Comp.Head)
		by.Header.Pitch = float32(data.Comp.Pitch)
		by.Header.Roll = float32(data.Comp.Roll)
	}
	if data.Ctd45 != nil {
		by.Header.Temperature = float32(data.Ctd45.Temp)
		by.Header.SoundVelocity = float32(data.Ctd45.Vel)
	}
	if data.Ctd60 != nil {
		by.Header.Temperature = float32(data.Ctd60.Temp)
		by.Header.SoundVelocity = float32(data.Ctd60.Vel)
	}
	if data.Pre != nil {
		by.Header.Pressure = float32(data.Pre.P)
	}
	//by should be merged already

	return by.Pack()
}

//MergeOICSonar merge sensor data to sonar structure
func MergeOICSonar(sonar *oic.Sonar, data *sensor.MixData) []byte {
	if data.Ap != nil {
		sonar.Header.NavFixLatitude = data.Ap.Lat
		sonar.Header.NavFixLongtitude = data.Ap.Lng
	}
	if data.Comp != nil {
		sonar.Header.VesselHeading = float32(data.Comp.Head)
		sonar.Header.Pitch = float32(data.Comp.Pitch)
		sonar.Header.Roll = float32(data.Comp.Roll)
	}
	if data.Ctd45 != nil {
		sonar.Header.Temperature = float32(data.Ctd45.Temp)
		sonar.Header.SoundVelocity = float32(data.Ctd45.Vel)
	}
	if data.Ctd60 != nil {
		sonar.Header.Temperature = float32(data.Ctd60.Temp)
		sonar.Header.SoundVelocity = float32(data.Ctd60.Vel)
	}
	if data.Pre != nil {
		sonar.Header.Pressure = float32(data.Pre.P)
	}
	//by should be merged already

	return sonar.Pack()
}

//RegenThread thread of merge data to OIC structure
func RegenThread(cfg *util.Cfg) {
	byfound := false
	ssfound := false
	subfound := false
	RegenbEnable = true
	var by *sensor.Node
	var ss *sensor.Node
	var sub *sensor.Node
	for {
		if RegenbEnable == false {
			logger.Println("RegenThread - RegenbEnable set to false,exit thread")
		}
		if byfound == false {
			if queueBy, ok := SQMap[sensor.BathyId]; ok {
				maplock.Lock()
				if by := queueBy.Pop(); by != nil {
					//bathy found!
					byfound = true
				}
				maplock.Unlock()
			}
		}
		if ssfound == false {
			if queueSs, ok := SQMap[sensor.SSId]; ok {
				maplock.Lock()
				if ss := queueSs.Pop(); ss != nil {
					//ss found!
					ssfound = true

				}
				maplock.Unlock()
			}
		}
		if subfound == false {
			if queueSub, ok := SQMap[sensor.SubbottomId]; ok {
				maplock.Lock()
				if sub := queueSub.Pop(); sub != nil {
					//sub found!
					subfound = true

				}
				maplock.Unlock()
			}
		}
		time.Sleep(time.Second * 1)
		if ssfound && subfound && byfound {
			if (ss.Time == by.Time) && (by.Time == sub.Time) {
				sensordata := &sensor.MixData{}
				maplock.Lock()
				sensordata.Ap = SQMap[sensor.ADID].FetchData(by.Time).(*sensor.AP)
				sensordata.Comp = SQMap[sensor.CompassHeader].FetchData(by.Time).(*sensor.Compass)
				sensordata.Ctd45 = SQMap[sensor.CTD4500Header].FetchData(by.Time).(*sensor.Ctd4500)
				sensordata.Ctd60 = SQMap[sensor.CTD6000Header].FetchData(by.Time).(*sensor.Ctd6000)
				sensordata.Pre = SQMap[sensor.PresureHeader].FetchData(by.Time).(*sensor.Presure)
				maplock.Unlock()
				if trace.File == nil {
					err := trace.New("BSSS", uint32(cfg.MaxSize)*1024*1024)
					if err != nil {
						logger.Fatal("Create bsss data file failed!")
					}
				}
				duby := by.Data.(*sensor.DuBathy)
				duss := ss.Data.(*sensor.DuSs)
				subbottom := sub.Data.(*sensor.Subbottom)

				bathy := formatBathy(duby)
				sonar := formatSonar(duss, subbottom)

				trace.Write(MergeOICBathy(bathy, sensordata), false)
				trace.Write(MergeOICSonar(sonar, sensordata), true)
			}
		}

	}
}
func formatBathy(duby *sensor.DuBathy) *oic.Bathy {
	bathy := &oic.Bathy{}
	bathy.Init()
	length := len(duby.PortBathy.DataAngle)
	for i := 0; i < length; i++ {
		bathy.PortAngle[i] = (float32)(duby.PortBathy.DataAngle[i])
		bathy.PortR[i] = (float32)(duby.PortBathy.DataDelay[i]) //should be same length

	}
	length = len(duby.StarboardBathy.DataAngle)
	for i := 0; i < length; i++ {
		bathy.StarboardAngle[i] = (float32)(duby.StarboardBathy.DataAngle[i])
		bathy.StarboardR[i] = (float32)(duby.StarboardBathy.DataDelay[i]) //should be same length
	}
	return bathy
}
func formatSonar(duss *sensor.DuSs, sub *sensor.Subbottom) *oic.Sonar {

	sonar := &oic.Sonar{}
	sonar.Init()
	length := len(duss.PortSs.Data)

	for i := 0; i < length; i++ {
		sonar.PortSidescan[i] = (int16)(duss.PortSs.Data[i])
	}
	length = len(duss.StarboardSs.Data)
	for i := 0; i < length; i++ {
		sonar.StarboardSidescan[i] = (int16)(duss.StarboardSs.Data[i])
	}
	length = len(sub.Sbdata)
	for i := 0; i < length; i++ {
		sonar.SubBottom[i] = (int16)(sub.Sbdata[i])
	}
	return sonar
}

//RelayThread wait for incoming data and relay to dest addr
func RelayThread(cfg *util.Cfg) {
	server := cfg.RelayIP + ":" + strconv.FormatInt(int64(cfg.RelaySenrPort), 10)
	_, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		logger.Fatal(fmt.Sprintf("RelayThread - Fatal error: %s", err.Error()))
	}
	for {
		RelayEnable = false
		conn, err := net.DialTimeout("tcp", server, time.Second*DIALTIMEOUT)
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
			if RelayEnable == false {
				logger.Println("RelayThread - RelayEnable set to false,exit thread")
				break
			}
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

//SetupServer local net listener
func SetupServer(cfg *util.Cfg) {
	listenaddr := "127.0.0.1:" + strconv.FormatInt(int64(cfg.SensorPort), 10)
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
	buf := make([]byte, 4096)
	for {

		n, err := conn.Read(buf)

		if err != nil {
			logger.Println(fmt.Sprintf("Connection - "+conn.RemoteAddr().String(), " connection error: ", err))
			return
		}
		if n > 0 {
			err = Dispatcher(buf[:n], maplock)
			if err != nil {
				logger.Println(fmt.Sprintf("Dispatcher error- ", err))
			}
			if RelayEnable {
				RelayChan <- buf[:n]
			}

		}

	}

}
