package sensor

import (
	"regener/util"
)

type bsss struct {
	header   bsssheadr
	wpara    workpara
	dpara    datapara
	payload  []ISensor
	checksum uint16
}
type subbottom struct {
	header   bsssheadr
	wpara    workpara
	dpara    datapara
	sbdata   [16384]uint16
	checksum uint16
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

type bsssheadr struct {
	ID         uint16
	Version    uint16
	PackageLen uint32
}

type workpara struct {
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

type datapara struct {
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

func (header *bsssheadr) Parse(recvbuf []byte) {

}
func (w *workpara) Parse(recvbuf []byte) {

}
func (d *datapara) Parse(recvbuf []byte) {

}
func (s *ss) Parse(recvbuf []byte) error {
	s.ID = uint32(util.BytesToUIntBE(32, recvbuf))
	s.Length = uint32(util.BytesToUIntBE(32, recvbuf[4:]))
	s.Para = uint32(util.BytesToUIntBE(32, recvbuf[8:]))
	s.Data = make([]float64, s.Length/4)
	for i := 0; i < int(s.Length/4); i++ {
		s.Data[i] = float64(util.ByteToFloat32BE(recvbuf[12+i*4:]))
	}
	return nil
}
func (b *bathy) Parse(recvbuf []byte) error {
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

//side scan
type ss struct {
	ID     uint32
	Length uint32
	Para   uint32
	Data   []float64 // count = Length/4
}

//bathy scan
type bathy struct {
	ID        uint32
	Length    uint32
	Para      uint32    //reserved
	DataAngle []float64 //rad,count = Length/4/2
	DataDelay []float64 //ms,count = Length/4/2
}
