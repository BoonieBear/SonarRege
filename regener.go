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
var sby, pby, sss, pss, nby, nss int16
var duby = new(sensor.DuBathy)

var duss = new(sensor.DuSs)

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
		sensor.SSId:          sensor.NewQueue(1000),
		sensor.BathyId:       sensor.NewQueue(1000),
		sensor.SubbottomId:   sensor.NewQueue(1000),
		sensor.APHeader:      sensor.NewQueue(1000),
		sensor.CompassHeader: sensor.NewQueue(1000),
		sensor.CTD6000Header: sensor.NewQueue(1000),
		sensor.CTD4500Header: sensor.NewQueue(1000),
		sensor.PresureHeader: sensor.NewQueue(1000),
		sensor.OASHeader:     sensor.NewQueue(1000),
		sensor.DVLHeader:     sensor.NewQueue(1000),
		sensor.PHINSHeader:   sensor.NewQueue(1000),
	}
	duby.PortBathy = nil
	duby.StarboardBathy = nil
	duss.PortSs = nil
	duss.StarboardSs = nil
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
		trace.Close()
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
	//fmt.Printf("on buffer = %d\n", len(buffer))
	for {
		//fmt.Printf("for buffer = %d\n", len(buffer))
		if len(buffer) < 8 {
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
		if uint16(util.BytesToUIntLE(16, buffer)) == sensor.SubbottomId {
			if uint16(util.BytesToUIntLE(16, buffer[2:])) == sensor.BsssVersion {
				length := util.BytesToUIntLE(32, buffer[4:])
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
		if uint16(util.BytesToUIntLE(16, buffer)) == sensor.RawID {
			if uint16(util.BytesToUIntLE(16, buffer[2:])) == sensor.RawVersion {
				length := util.BytesToUIntLE(32, buffer[4:])

				if len(buffer) < int(length) {
					//no enough buffer
					break
				} else {
					//fmt.Printf(" ======= rawID %d\n", length)
					data := make([]byte, length)
					copy(data, buffer[:length])
					buffer = append(buffer[:0], buffer[length:]...)
					return nil
					//return DispatchSensor(data, queuelock)
				}
			}
		}
		if len(buffer) > 1 {
			//shift 1 bytes
			//fmt.Printf("shift byte %x\n", buffer[0])
			buffer = append(buffer[:0], buffer[1:]...)
		}

	}
	return nil
}

func DispatchBsss(recvbuf []byte, queuelock *sync.Mutex) error {
	bs := &sensor.Bsss{}
	bs.Parse(recvbuf)
	dby := &sensor.DuBathy{}
	dss := &sensor.DuSs{}

	var hasBy, hasSs bool

	for _, v := range bs.Payload {
		if value, ok := v.(*sensor.SingelBathy); ok {
			bathy := value
			if bathy.ID == uint32(sensor.PortByID) {
				duby.PortBathy = value
				pby++
			} else {
				duby.StarboardBathy = value
				sby++
			}
			if duby.PortBathy != nil && duby.StarboardBathy != nil {
				dby.StarboardBathy = duby.StarboardBathy
				dby.PortBathy = duby.PortBathy
				dby.Wpara = bs.Wpara
				hasBy = true
				duby.StarboardBathy = nil
				duby.PortBathy = nil
			}

		}

		if value, ok := v.(*sensor.Ss); ok {
			ss := value
			if ss.ID == uint32(sensor.PortSSID) {
				duss.PortSs = value
				pss++
			} else {
				duss.StarboardSs = value
				sss++
			}
			if duss.PortSs != nil && duss.StarboardSs != nil {
				dss.PortSs = duss.PortSs
				dss.StarboardSs = duss.StarboardSs
				dss.Wpara = bs.Wpara
				hasSs = true
				duss.PortSs = nil
				duss.StarboardSs = nil
			}

		}
	}
	if hasSs {
		node := &sensor.Node{
			Time: time.Unix(int64(bs.Dpara.EmitTime1st), int64(bs.Dpara.EmitTime2nd*1000)),
			Data: dss,
		}
		queuelock.Lock()
		if queue, ok := SQMap[sensor.SSId]; ok {

			queue.Push(node)
			nss++
		}
		queuelock.Unlock()
	}
	if hasBy {
		node := &sensor.Node{
			Time: time.Unix(int64(bs.Dpara.EmitTime1st), int64(bs.Dpara.EmitTime2nd*1000)),
			Data: dby,
		}
		queuelock.Lock()
		if queue, ok := SQMap[sensor.BathyId]; ok {

			queue.Push(node)
			nby++
		}
		queuelock.Unlock()
	}
	//fmt.Printf("%d %d %d %d %d %d\n", pby, sby, pss, sss, nby, nss)
	return nil
}
func DispatchSensor(recvbuf []byte, queuelock *sync.Mutex) error {
	totallength := util.BytesToUIntLE(32, recvbuf[4:])
	index := uint64(8)
	if string(recvbuf[8:14]) == "$GPZDA" {
		index += 38
	} else {
		index += 18
	}

	for {
		id := uint16(util.BytesToUIntBE(16, recvbuf[index:]))
		length := util.BytesToUIntBE(16, recvbuf[index+2:])
		//fmt.Printf("id %x index %d\n", id, index)
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
		case sensor.OASHeader:
			oas := &sensor.OAS{}
			err := oas.Parse(recvbuf[index:])
			if err == nil {
				node := &sensor.Node{
					Time: oas.Time,
					Data: oas,
				}
				queuelock.Lock()
				if queue, ok := SQMap[sensor.OASHeader]; ok {
					queue.Push(node)
				}
				queuelock.Unlock()

			} else {
				log.Println(err)
			}
		case sensor.DVLHeader:
			dvl := &sensor.DVL{}
			err := dvl.Parse(recvbuf[index:])
			if err == nil {
				node := &sensor.Node{
					Time: dvl.Time,
					Data: dvl,
				}
				queuelock.Lock()
				if queue, ok := SQMap[sensor.DVLHeader]; ok {
					queue.Push(node)
				}
				queuelock.Unlock()

			} else {
				log.Println(err)
			}
		case sensor.PHINSHeader:
			phi := &sensor.PHINS{}
			err := phi.Parse(recvbuf[index:])
			if err == nil {
				node := &sensor.Node{
					Time: phi.Time,
					Data: phi,
				}
				queuelock.Lock()
				if queue, ok := SQMap[sensor.PHINSHeader]; ok {
					queue.Push(node)
					//fmt.Printf("push phins %v\n", node)
				}
				queuelock.Unlock()

			} else {
				log.Println(err)
			}
		}
		index += 4
		index += length
		if index == totallength {
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
		sonar.Header.ShipX, sonar.Header.ShipY = util.Deg2utm(data.Ap.Lat, data.Ap.Lng)
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
	if data.Phins != nil {
		by.Header.VesselHeading = float32(data.Phins.Head)
		by.Header.Pitch = float32(data.Phins.Pitch)
		by.Header.Roll = float32(data.Phins.Roll)
		by.Header.FishPitch = float32(data.Phins.Pitch)
		by.Header.FishRoll = float32(data.Phins.Roll)
		by.Header.ShipCourse = float32(data.Phins.Head)
		by.Header.FishHeading = float32(data.Phins.Head)
		//fmt.Printf("Phins write %f\n", sonar.Header.VesselHeading)
	}
	if data.Dvl != nil {
		by.Header.FishAltitude = float32(data.Dvl.Botrange)

	}

	return by.Pack()
}

//MergeOICSonar merge sensor data to sonar structure
func MergeOICSonar(sonar *oic.Sonar, data *sensor.MixData) []byte {

	if data.Ap != nil {
		sonar.Header.NavFixLatitude = data.Ap.Lat
		sonar.Header.NavFixLongtitude = data.Ap.Lng
		sonar.Header.ShipX, sonar.Header.ShipY = util.Deg2utm(data.Ap.Lat, data.Ap.Lng)
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
	if data.Phins != nil {
		sonar.Header.VesselHeading = float32(data.Phins.Head)
		sonar.Header.Pitch = float32(data.Phins.Pitch)
		sonar.Header.Roll = float32(data.Phins.Roll)
		sonar.Header.FishPitch = float32(data.Phins.Pitch)
		sonar.Header.FishRoll = float32(data.Phins.Roll)
		sonar.Header.ShipCourse = float32(data.Phins.Head)
		sonar.Header.FishHeading = float32(data.Phins.Head)
		//fmt.Printf("Phins write %f\n", sonar.Header.VesselHeading)
	}
	if data.Dvl != nil {
		sonar.Header.FishAltitude = float32(data.Dvl.Botrange)
		//sonar.Header.Reserved2[0] = float32(data.Dvl.Botrange)
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
		fmt.Println("======================")
		for id, queue := range SQMap {
			ncount := queue.Count()
			fmt.Printf("queue 0x%x has %d items:\n", id, ncount)
		}
		if RegenbEnable == false {
			logger.Println("RegenThread - RegenbEnable set to false,exit thread")
		}
		if byfound == false {
			if queueBy, ok := SQMap[sensor.BathyId]; ok {
				maplock.Lock()
				if by = queueBy.Pop(); by != nil {
					fmt.Println("bathy found!")
					byfound = true
				}
				maplock.Unlock()
			}
		}
		if ssfound == false {
			if queueSs, ok := SQMap[sensor.SSId]; ok {
				maplock.Lock()
				if ss = queueSs.Pop(); ss != nil {
					fmt.Println("ss found!")
					ssfound = true

				}
				maplock.Unlock()
			}
		}
		if subfound == false {
			if queueSub, ok := SQMap[sensor.SubbottomId]; ok {
				maplock.Lock()
				if sub = queueSub.Pop(); sub != nil {
					fmt.Println("sub found!")
					subfound = true

				}
				maplock.Unlock()
			}
		}
		//fmt.Printf("merge time = %v\n", by.Time)
		time.Sleep(time.Millisecond * 1000)
		if ssfound && byfound && subfound {

			sensordata := &sensor.MixData{}
			maplock.Lock()
			if SQMap[sensor.APHeader].Count() > 0 {
				sensordata.Ap = SQMap[sensor.APHeader].FetchData(by.Time).(*sensor.AP)
			}
			if SQMap[sensor.CompassHeader].Count() > 0 {
				sensordata.Comp = SQMap[sensor.CompassHeader].FetchData(by.Time).(*sensor.Compass)
			}
			if SQMap[sensor.CTD4500Header].Count() > 0 {
				sensordata.Ctd45 = SQMap[sensor.CTD4500Header].FetchData(by.Time).(*sensor.Ctd4500)
			}
			if SQMap[sensor.CTD6000Header].Count() > 0 {
				sensordata.Ctd60 = SQMap[sensor.CTD6000Header].FetchData(by.Time).(*sensor.Ctd6000)
			}
			if SQMap[sensor.PresureHeader].Count() > 0 {
				sensordata.Pre = SQMap[sensor.PresureHeader].FetchData(by.Time).(*sensor.Presure)
			}
			if SQMap[sensor.OASHeader].Count() > 0 {
				sensordata.Oas = SQMap[sensor.OASHeader].FetchData(by.Time).(*sensor.OAS)
			}
			if SQMap[sensor.DVLHeader].Count() > 0 {
				sensordata.Dvl = SQMap[sensor.DVLHeader].FetchData(by.Time).(*sensor.DVL)
				fmt.Printf("fish pitch = %f\n", sensordata.Dvl.Pitch)
			}
			if SQMap[sensor.PHINSHeader].Count() > 0 {
				sensordata.Phins = SQMap[sensor.PHINSHeader].FetchData(by.Time).(*sensor.PHINS)
			}
			maplock.Unlock()
			if trace.File == nil {
				err := trace.New("BSSS", uint32(cfg.MaxSize)*1024*1024)
				fmt.Printf("trace file size = %d\n", trace.MaxSize)
				if err != nil {
					logger.Fatal("Create bsss data file failed!")
				}
			}
			duby := new(sensor.DuBathy)
			duss := new(sensor.DuSs)
			subbottom := new(sensor.Subbottom)

			if byfound {
				duby = by.Data.(*sensor.DuBathy)
			}
			if ssfound {
				duss = ss.Data.(*sensor.DuSs)
			}

			if subfound {
				subbottom = sub.Data.(*sensor.Subbottom)
			}

			sonar := formatSonar(duss, subbottom)
			sonar.Header.Sec = uint32(by.Time.Unix())
			sonar.Header.USec = uint32(by.Time.Nanosecond() / 1000)
			sonar.Header.PingTime = float64(by.Time.Unix())
			trace.Write(MergeOICSonar(sonar, sensordata), false)
			bathy := formatBathy(duby)
			bathy.Header.Sec = uint32(by.Time.Unix())
			bathy.Header.USec = uint32(by.Time.Nanosecond() / 1000)
			bathy.Header.PingTime = float64(by.Time.Unix())
			trace.Write(MergeOICBathy(bathy, sensordata), false)
		}
		byfound = false
		ssfound = false
		subfound = false

	}
}
func formatBathy(duby *sensor.DuBathy) *oic.Bathy {
	bathy := &oic.Bathy{}
	bathy.Init()
	length := len(duby.PortBathy.DataAngle)
	bathy.Header.Channel[0].Frequency = duby.Wpara.PortStartFq
	bathy.Header.Channel[0].Empt = 1
	bathy.Header.Channel[1].Frequency = duby.Wpara.StarboardFq
	bathy.Header.Channel[1].Empt = 1
	bathy.Header.Channel[2].Frequency = duby.Wpara.PortStartFq
	bathy.Header.Channel[2].Empt = 1
	if length > 0 {
		bathy.PortAngle = make([]float32, length)
		bathy.PortR = make([]float32, length)
	}
	for i := 0; i < length; i++ {
		bathy.PortAngle[i] = (float32)(duby.PortBathy.DataAngle[i])
		bathy.PortR[i] = (float32)(duby.PortBathy.DataDelay[i]) //should be same length

	}
	length = len(duby.StarboardBathy.DataAngle)

	if length > 0 {
		bathy.StarboardAngle = make([]float32, length)
		bathy.StarboardR = make([]float32, length)
	}
	for i := 0; i < length; i++ {
		bathy.StarboardAngle[i] = (float32)(duby.StarboardBathy.DataAngle[i])
		bathy.StarboardR[i] = (float32)(duby.StarboardBathy.DataDelay[i]) //should be same length
	}
	return bathy
}
func formatSonar(duss *sensor.DuSs, sub *sensor.Subbottom) *oic.Sonar {

	sonar := &oic.Sonar{}
	sonar.Init()
	sonar.Header.Channel[0].Frequency = duss.Wpara.PortStartFq
	sonar.Header.Channel[0].Empt = 1
	sonar.Header.Channel[1].Frequency = duss.Wpara.StarboardFq
	sonar.Header.Channel[1].Empt = 1
	sonar.Header.Channel[2].Frequency = sub.Wpara.PortStartFq
	sonar.Header.Channel[2].Empt = 1
	length := len(duss.PortSs.Data)
	if length > 0 {
		sonar.PortSidescan = make([]int16, length)
	}
	for i := 0; i < length; i++ {
		sonar.PortSidescan[i] = (int16)(duss.PortSs.Data[i])
	}
	length = len(duss.StarboardSs.Data)
	if length > 0 {
		sonar.StarboardSidescan = make([]int16, length)
	}
	for i := 0; i < length; i++ {
		sonar.StarboardSidescan[i] = (int16)(duss.StarboardSs.Data[i])
	}
	if sub != nil {
		length = len(sub.Sbdata)
		if length > 0 {
			sonar.SubBottom = make([]int16, length)
		}
		for i := 0; i < length; i++ {
			sonar.SubBottom[i] = (int16)(sub.Sbdata[i])
		}
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
				//logger.Println(fmt.Sprintf("Dispatcher error: ", err))
			}
			if RelayEnable {
				RelayChan <- buf[:n]
			}

		}

	}

}
