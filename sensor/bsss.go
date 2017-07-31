package sensor

import (
	"fmt"
	"regener/util"
)

type Bsss struct {
	Header   Bsssheadr
	Wpara    Workpara
	Dpara    Datapara
	Payload  []ISensor
	Checksum uint16
}
type Subbottom struct {
	Header   Bsssheadr
	Wpara    Workpara
	Dpara    Datapara
	Sbdata   [16384]uint16
	Checksum uint16
}

const (
	ADID            uint16 = 0x01
	ReserveID       uint16 = 0x02
	PortByID        uint16 = 0x04
	StarboardByID   uint16 = 0x08
	PortSSID        uint16 = 0x10
	StarboardSSID   uint16 = 0x20
	PortVector      uint16 = 0x40
	StarboardVector uint16 = 0x80
)

type Bsssheadr struct {
	ID         uint16
	Version    uint16
	PackageLen uint32
}

type Workpara struct {
	Serial           uint16 //0~0xFFFF
	PulseLen         uint16 //x671/2^26 sec
	PortStartFq      uint32 //x2Hz
	StarboardFq      uint32 //x2Hz
	PortChirpFq      uint32 //x16Hz/s
	StarboardChirpFq uint32 //x16Hz/s
	RecvLatecy       uint16 //x671/2^26 sec
	Sampling         uint16
	EmitInterval     uint16 //×1/16384 sec
	RelativeGain     uint16
	StatusFlag       uint16
	TVGLatecy        uint16 //×67/2^26
	TVGRefRate       uint16
	TVGCtrl          uint16
	TVGFactor        uint16 //x0.1
	TVGAntenu        uint16 //×0.00001dB/m
	TVGGain          uint16 //×0.1dB
	RetFlag          uint16
	CMDFlag          uint32
	FixedTVG         uint32
	Reserved         uint32
}

type Datapara struct {
	DataID         uint16
	IsNewEmit      uint16
	EmitCount      uint32
	Reserved       uint32
	DataParaID     uint32 //share value as ID
	EmitShiftPoint uint32
	EmitTime1st    uint32
	EmitTime2nd    uint32
	Velocity       uint32
	DataCount      uint16
	Reserve1       [3]uint16
}

//side scan
type Ss struct {
	ID     uint32
	Length uint32
	Para   uint32
	Data   []float64 // count = Length/4
}

//single bathy scan
type SingelBathy struct {
	ID        uint32
	Length    uint32
	Para      uint32    //reserved
	DataAngle []float64 //rad,count = Length/4/2
	DataDelay []float64 //ms,count = Length/4/2
}

//ss struct for star/port ss
type DuSs struct {
	PortSs      *Ss
	StarboardSs *Ss
}

func (duss *DuSs) Parse(recvbuf []byte) error {
	return nil
}

func (duss *DuSs) Dump() {

}

//Bathy struct for star/port bathy
type DuBathy struct {
	PortBathy      *SingelBathy
	StarboardBathy *SingelBathy
}

func (duby *DuBathy) Dump() {

}

func (duby *DuBathy) Parse(recvbuf []byte) error {
	return nil
}
func (bsss *Bsss) Parse(recvbuf []byte) error {
	bsss.Header.Parse(recvbuf)
	bsss.Wpara.Parse(recvbuf[8:])
	bsss.Dpara.Parse(recvbuf[64:])

	index := 104 //shift index byte
	for {
		switch uint16(util.BytesToUIntBE(16, recvbuf[index:])) {
		case PortByID:
		case StarboardByID:
			by := &SingelBathy{}
			by.Parse(recvbuf[index:])
			bsss.Payload = append(bsss.Payload, by)
			index = index + 12 + int(by.Length)
		case PortSSID:
		case StarboardSSID:
			s := &Ss{}
			s.Parse(recvbuf[index:])
			bsss.Payload = append(bsss.Payload, s)
			index = index + 12 + int(s.Length)
		}
		if index+2 == int(bsss.Header.PackageLen) {
			break
		}
	}
	return nil
}

//Display all the item in Bsss
func (bsss *Bsss) Dump() {
	fmt.Println("=== display BSSS data ===")
	bsss.Header.Dump()
	bsss.Wpara.Dump()
	bsss.Dpara.Dump()
	for _, v := range bsss.Payload {
		v.Dump()
	}

}

func (sub *Subbottom) Parse(recvbuf []byte) error {
	sub.Header.Parse(recvbuf)
	sub.Wpara.Parse(recvbuf[8:])
	sub.Dpara.Parse(recvbuf[64:])
	for i := range sub.Sbdata {
		sub.Sbdata[i] = uint16(util.BytesToUIntBE(16, recvbuf[104+2*i:]))
	}
	return nil
}

func (sub *Subbottom) Dump() {
	fmt.Println("=== display sub bottom data ===")
	sub.Header.Dump()
	sub.Wpara.Dump()
	sub.Dpara.Dump()
	fmt.Println("print top 10 sub bottom data: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d %d |", i, sub.Sbdata[i])
	}
}

func (header *Bsssheadr) Parse(recvbuf []byte) {
	header.ID = uint16(util.BytesToUIntBE(16, recvbuf))
	header.Version = uint16(util.BytesToUIntBE(16, recvbuf[2:]))
	header.PackageLen = uint32(util.BytesToUIntBE(32, recvbuf[4:]))
}

func (header *Bsssheadr) Dump() {
	fmt.Println("--- display bsss header data ---")
	fmt.Printf("- headr id = %d", header.ID)
	fmt.Printf("- headr version = %d", header.Version)
	fmt.Printf("- headr package len = %d", header.PackageLen)
}

func (w *Workpara) Parse(recvbuf []byte) {
	w.Serial = uint16(util.BytesToUIntBE(16, recvbuf))
	w.PulseLen = uint16(util.BytesToUIntBE(16, recvbuf[2:]))
	w.PortStartFq = uint32(util.BytesToUIntBE(32, recvbuf[4:]))       //x2Hz
	w.StarboardFq = uint32(util.BytesToUIntBE(32, recvbuf[8:]))       //x2Hz
	w.PortChirpFq = uint32(util.BytesToUIntBE(32, recvbuf[12:]))      //x16Hz/s
	w.StarboardChirpFq = uint32(util.BytesToUIntBE(32, recvbuf[16:])) //x16Hz/s
	w.RecvLatecy = uint16(util.BytesToUIntBE(16, recvbuf[20:]))       //x671/2^26 sec
	w.Sampling = uint16(util.BytesToUIntBE(16, recvbuf[22:]))
	w.EmitInterval = uint16(util.BytesToUIntBE(16, recvbuf[24:])) //×1/16384 sec
	w.RelativeGain = uint16(util.BytesToUIntBE(16, recvbuf[26:]))
	w.StatusFlag = uint16(util.BytesToUIntBE(16, recvbuf[28:]))
	w.TVGLatecy = uint16(util.BytesToUIntBE(16, recvbuf[30:])) //×67/2^26
	w.TVGRefRate = uint16(util.BytesToUIntBE(16, recvbuf[32:]))
	w.TVGCtrl = uint16(util.BytesToUIntBE(16, recvbuf[34:]))
	w.TVGFactor = uint16(util.BytesToUIntBE(16, recvbuf[36:])) //x0.1
	w.TVGAntenu = uint16(util.BytesToUIntBE(16, recvbuf[38:])) //×0.00001dB/m
	w.TVGGain = uint16(util.BytesToUIntBE(16, recvbuf[40:]))   //×0.1dB
	w.RetFlag = uint16(util.BytesToUIntBE(16, recvbuf[42:]))
	w.CMDFlag = uint32(util.BytesToUIntBE(32, recvbuf[44:]))
	w.FixedTVG = uint32(util.BytesToUIntBE(32, recvbuf[48:]))
	w.Reserved = uint32(util.BytesToUIntBE(32, recvbuf[52:]))
}

func (w *Workpara) Dump() {
	fmt.Println("--- display work para data ---")
	fmt.Printf("- Serial = %d", w.Serial)
	fmt.Printf("- PulseLen = %d", w.PulseLen)
	fmt.Printf("- PortStartFq = %d", w.PortStartFq)
	fmt.Printf("- StarboardFq = %d", w.StarboardFq)
	fmt.Printf("- PortChirpFq = %d", w.PortChirpFq)
	fmt.Printf("- StarboardChirpFq = %d", w.StarboardChirpFq)
	fmt.Printf("- RecvLatecy = %d", w.RecvLatecy)
	fmt.Printf("- Sampling = %d", w.Sampling)
	fmt.Printf("- EmitInterval = %d", w.EmitInterval)
	fmt.Printf("- RelativeGain = %d", w.RelativeGain)
	fmt.Printf("- StatusFlag = %d", w.StatusFlag)
	fmt.Printf("- TVGLatecy = %d", w.TVGLatecy)
	fmt.Printf("- TVGRefRate = %d", w.TVGRefRate)
	fmt.Printf("- TVGCtrl = %d", w.TVGCtrl)
	fmt.Printf("- TVGFactor = %d", w.TVGFactor)
	fmt.Printf("- TVGAntenu = %d", w.TVGAntenu)
	fmt.Printf("- TVGGain = %d", w.TVGGain)
	fmt.Printf("- RetFlag = %d", w.RetFlag)
	fmt.Printf("- CMDFlag = %d", w.CMDFlag)
	fmt.Printf("- FixedTVG = %d", w.FixedTVG)

}

func (d *Datapara) Parse(recvbuf []byte) {
	d.DataID = uint16(util.BytesToUIntBE(16, recvbuf))
	d.IsNewEmit = uint16(util.BytesToUIntBE(16, recvbuf[2:]))
	d.EmitCount = uint32(util.BytesToUIntBE(32, recvbuf[4:]))
	d.Reserved = uint32(util.BytesToUIntBE(32, recvbuf[8:]))
	d.DataParaID = uint32(util.BytesToUIntBE(32, recvbuf[12:])) //share value as ID
	d.EmitShiftPoint = uint32(util.BytesToUIntBE(32, recvbuf[16:]))
	d.EmitTime1st = uint32(util.BytesToUIntBE(32, recvbuf[20:]))
	d.EmitTime2nd = uint32(util.BytesToUIntBE(32, recvbuf[24:]))
	d.Velocity = uint32(util.BytesToUIntBE(32, recvbuf[28:]))
	d.DataCount = uint16(util.BytesToUIntBE(16, recvbuf[32:]))
	d.Reserve1[0] = uint16(util.BytesToUIntBE(16, recvbuf[34:]))
	d.Reserve1[1] = uint16(util.BytesToUIntBE(16, recvbuf[36:]))
	d.Reserve1[2] = uint16(util.BytesToUIntBE(16, recvbuf[38:]))
}

func (d *Datapara) Dump() {
	fmt.Println("--- display data para data ---")
	fmt.Printf("- DataID = %d", d.DataID)
	fmt.Printf("- IsNewEmit = %d", d.IsNewEmit)
	fmt.Printf("- EmitCount = %d", d.EmitCount)
	fmt.Printf("- Reserved = %d", d.Reserved)
	fmt.Printf("- DataParaID = %d", d.DataParaID)
	fmt.Printf("- EmitShiftPoint = %d", d.EmitShiftPoint)
	fmt.Printf("- EmitTime1st = %d", d.EmitTime1st)
	fmt.Printf("- EmitTime2nd = %d", d.EmitTime2nd)
	fmt.Printf("- Velocity = %d", d.Velocity)
	fmt.Printf("- DataCount = %d", d.DataCount)
}

func (s *Ss) Parse(recvbuf []byte) error {
	s.ID = uint32(util.BytesToUIntBE(32, recvbuf))
	s.Length = uint32(util.BytesToUIntBE(32, recvbuf[4:]))
	s.Para = uint32(util.BytesToUIntBE(32, recvbuf[8:]))
	s.Data = make([]float64, s.Length/4)
	for i := 0; i < int(s.Length/4); i++ {
		s.Data[i] = float64(util.ByteToFloat32BE(recvbuf[12+i*4:]))
	}
	return nil
}

func (s *Ss) Dump() {
	fmt.Println("=== display side scan data ===")
	fmt.Printf("- ID = %d", s.ID)
	fmt.Printf("- Length = %d", s.Length)
	fmt.Printf("- Para = %d", s.Para)
	fmt.Println("print top 10 side scan data: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d %f |", i, s.Data[i])
	}
}

func (b *SingelBathy) Parse(recvbuf []byte) error {
	b.ID = uint32(util.BytesToUIntBE(32, recvbuf))
	b.Length = uint32(util.BytesToUIntBE(32, recvbuf[4:]))
	b.Para = uint32(util.BytesToUIntBE(32, recvbuf[8:]))
	b.DataAngle = make([]float64, b.Length/4/2)
	for i := 0; i < int(b.Length/4/2); i++ {
		b.DataAngle[i] = float64(util.ByteToFloat32BE(recvbuf[12+i*4:]))
	}
	b.DataDelay = make([]float64, b.Length/4/2)
	for i := int(b.Length / 4 / 2); i < int(b.Length/4); i++ {
		b.DataDelay[i] = float64(util.ByteToFloat32BE(recvbuf[12+i*4:]))
	}
	return nil
}

func (b *SingelBathy) Dump() {
	fmt.Println("=== display single bathy data ===")
	fmt.Printf("- ID = %d", b.ID)
	fmt.Printf("- Length = %d", b.Length)
	fmt.Printf("- Para = %d", b.Para)
	fmt.Println("print top 10 angle data: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d %f |", i, b.DataAngle[i])
	}
	fmt.Println("print top 10 delay data: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d %f |", i, b.DataDelay[i])
	}
}
