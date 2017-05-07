package oic

import (
	"regener/util"
)

type STATUS uint8

const (
	FOCUSAUTOMANUAL          STATUS = 0
	FOCUSMANUALDISABLEENABLE STATUS = 1
	PINGRATEAUTOMANUAL       STATUS = 2
	TVGAUTOMANUAL            STATUS = 3
	CALIBOFFON               STATUS = 4
	OUTPUTMODEPROCRAW        STATUS = 5
	SHADOWMASK               STATUS = 6
	QUAlITYBIT               STATUS = 7
)

type Channel struct {
	Type      uint8
	Side      uint8
	Size      uint8
	Empt      uint8
	Frequency uint32
	Samples   uint32
}

func (channel *Channel) Pack() []byte {
	var buf []byte
	buf = append(buf, util.UIntToBytesBE(8, uint64(channel.Type))...)
	buf = append(buf, util.UIntToBytesBE(8, uint64(channel.Side))...)
	buf = append(buf, util.UIntToBytesBE(8, uint64(channel.Size))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(channel.Frequency))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(channel.Samples))...)
	return buf
}

type OIC_Header struct {
	Kind                uint32
	Type                uint32
	DataSize            uint32
	ProcStatus          uint32
	ClientSize          uint32
	FishStatus          STATUS
	NavUsed             uint8
	NavType             uint8
	UTMZone             uint32
	ShipX               float64
	ShipY               float64
	ShipCourse          float32
	ShipSpeed           float32
	Sec                 uint32
	uSec                uint32
	SpareGain           float32
	FishHeading         float32
	FishDepth           float32
	FishRange           float32
	FishPulseWidth      float32
	GainC0              float32
	GainC1              float32
	GainC2              float32
	FishPitch           float32
	FishRoll            float32
	FishYaw             float32
	Temperature         float32
	FishX               float64
	FishY               float64
	FishLayback         float32
	FishAltitude        float32
	FishAltitudeSamples uint32
	FishPingPeriod      float32
	SoundVelocity       float32
	Reserved1           uint32
	ChanNum             uint32
	ChanOffset          [8]uint32
	Channel             [8]Channel
	Reserved2           [5]uint32
	NavFixLatitude      float64
	NavFixLongtitude    float64
	HDOP                float32
	EllipsoidElevation  float32
	VesselHeading       float32
	Pitch               float32
	Roll                float32
	Heave               float32
	Draft               float32
	Tide                float32
	Reserved3           uint32
	Pressure            float32
	Reserved4           [13]uint32
	AuxFloat4           float32
	Reserved5           [4]uint32
	Aux3                float64
	Aux4                float64
	Reserved6           [2]uint32
	PingTime            float64
	Reserved7           [18]uint32
}

//initial OIC header
func OICInit(header *OIC_Header, isbathy bool) {
	header.Kind = 0x4F49432F
	header.Type = 26
	if isbathy {
		header.DataSize = 3752
	} else {
		header.DataSize = 35304
	}
	header.ClientSize = 248
	header.FishStatus = FOCUSAUTOMANUAL
	header.NavUsed = 6
	header.NavType = 1
	header.UTMZone = 0x3200DA02
	//field have not initilized use the default value if other not assign value to them

}

//pack the oic header struct
func (header *OIC_Header) Pack() []byte {
	var buf []byte
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.Kind))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.Type))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.DataSize))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.ProcStatus))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.ClientSize))...)
	buf = append(buf, util.UIntToBytesBE(8, uint64(header.FishStatus))...)
	buf = append(buf, util.UIntToBytesBE(8, uint64(header.NavUsed))...)
	buf = append(buf, util.UIntToBytesBE(8, uint64(header.NavType))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.UTMZone))...)
	buf = append(buf, util.Float64ToByteBE(header.ShipX)...)
	buf = append(buf, util.Float64ToByteBE(header.ShipY)...)
	buf = append(buf, util.Float32ToByteBE(header.ShipCourse)...)
	buf = append(buf, util.Float32ToByteBE(header.ShipSpeed)...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.Sec))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.uSec))...)
	buf = append(buf, util.Float32ToByteBE(header.SpareGain)...)
	buf = append(buf, util.Float32ToByteBE(header.FishHeading)...)
	buf = append(buf, util.Float32ToByteBE(header.FishDepth)...)
	buf = append(buf, util.Float32ToByteBE(header.FishRange)...)
	buf = append(buf, util.Float32ToByteBE(header.FishPulseWidth)...)
	buf = append(buf, util.Float32ToByteBE(header.GainC0)...)
	buf = append(buf, util.Float32ToByteBE(header.GainC1)...)
	buf = append(buf, util.Float32ToByteBE(header.GainC2)...)
	buf = append(buf, util.Float32ToByteBE(header.FishPitch)...)
	buf = append(buf, util.Float32ToByteBE(header.FishRoll)...)
	buf = append(buf, util.Float32ToByteBE(header.FishYaw)...)
	buf = append(buf, util.Float32ToByteBE(header.Temperature)...)
	buf = append(buf, util.Float64ToByteBE(header.FishX)...)
	buf = append(buf, util.Float64ToByteBE(header.FishY)...)
	buf = append(buf, util.Float32ToByteBE(header.FishLayback)...)
	buf = append(buf, util.Float32ToByteBE(header.FishAltitude)...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.FishAltitudeSamples))...)
	buf = append(buf, util.Float32ToByteBE(header.FishPingPeriod)...)
	buf = append(buf, util.Float32ToByteBE(header.SoundVelocity)...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.Reserved1))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.ChanNum))...)
	for i := 0; i < 8; i++ {
		buf = append(buf, util.UIntToBytesBE(32, uint64(header.ChanOffset[i]))...)
	}
	for i := 0; i < 8; i++ {
		buf = append(buf, header.Channel[i].Pack()...)
	}
	for i := 0; i < 5; i++ {
		buf = append(buf, util.UIntToBytesBE(32, uint64(header.Reserved2[i]))...)
	}
	buf = append(buf, util.Float64ToByteBE(header.NavFixLatitude)...)
	buf = append(buf, util.Float64ToByteBE(header.NavFixLongtitude)...)
	buf = append(buf, util.Float32ToByteBE(header.HDOP)...)
	buf = append(buf, util.Float32ToByteBE(header.EllipsoidElevation)...)
	buf = append(buf, util.Float32ToByteBE(header.VesselHeading)...)
	buf = append(buf, util.Float32ToByteBE(header.Pitch)...)
	buf = append(buf, util.Float32ToByteBE(header.Roll)...)
	buf = append(buf, util.Float32ToByteBE(header.Heave)...)
	buf = append(buf, util.Float32ToByteBE(header.Draft)...)
	buf = append(buf, util.Float32ToByteBE(header.Tide)...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.Reserved3))...)
	buf = append(buf, util.Float32ToByteBE(header.Pressure)...)
	for i := 0; i < 13; i++ {
		buf = append(buf, util.UIntToBytesBE(32, uint64(header.Reserved4[i]))...)
	}
	buf = append(buf, util.Float32ToByteBE(header.AuxFloat4)...)
	for i := 0; i < 4; i++ {
		buf = append(buf, util.UIntToBytesBE(32, uint64(header.Reserved5[i]))...)
	}
	buf = append(buf, util.Float64ToByteBE(header.Aux3)...)
	buf = append(buf, util.Float64ToByteBE(header.Aux4)...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.Reserved6[0]))...)
	buf = append(buf, util.UIntToBytesBE(32, uint64(header.Reserved6[1]))...)
	buf = append(buf, util.Float64ToByteBE(header.PingTime)...)
	for i := 0; i < 18; i++ {
		buf = append(buf, util.UIntToBytesBE(32, uint64(header.Reserved7[i]))...)
	}
	return buf
}

type Sonar struct {
	Header            OIC_Header
	PortSidescan      []int16
	StarboardSidescan []int16
	SubBottom         []int16
}
type Bathy struct {
	Header         OIC_Header
	PortAngle      []float32
	PortR          []float32
	StarboardAngle []float32
	StarboardR     []float32
	Unknow         [56]int16
}
