package sensor

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
	sbdata   []sb
	checksum uint16
}

//subbottom
type sb struct {
	ID uint32
}
type bsssheadr struct {
	ID         int16
	Version    int16
	PackageLen int32
}

type workpara struct {
	Serial           int16 //0~0xFFFF
	PulseLen         int16 //x671/2^26 sec
	PortStartFq      int32 //x2Hz
	StarboardFq      int32 //x2Hz
	PortChirpFq      int32 //x16Hz/s
	StarboardChirpFq int32 //x16Hz/s
	RecvLatecy       int16 //x671/2^26 sec
	Sampling         int16
	EmitInterval     int16 //×1/16384 sec
	RelativeGain     int16
	StatusFlag       int16
	TVGLatecy        int16 //×67/2^26
	TVGRefRate       int16
	TVGCtrl          int16
	TVGFactor        int16 //x0.1
	TVGAntenu        int16 //×0.00001dB/m
	TVGGain          int16 //×0.1dB
	RetFlag          int16
	CMDFlag          int32
	FixedTVG         int32
	Reserved         int32
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

func (s *ss) Parse(recvbuf []byte) error {
	return nil
}
func (b *bathy) Parse(recvbuf []byte) error {
	return nil
}

//side scan
type ss struct {
	ID    uint32
	Count uint32
	Para  uint32
	Data  []float64 // length = Count/4
}

//bathy scan
type bathy struct {
	ID        uint32
	Count     uint32
	Para      uint32    //reserved
	DataAngle []float64 //rad,length = Count/4/2
	DataDelay []float64 //ms,length = Count/4/2
}
