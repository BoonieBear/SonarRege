package oic

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

type Chan struct {
	Type      uint8
	Side      uint8
	Size      uint8
	Empt      uint8
	Frequency uint32
	Samples   uint32
}
type Cfg struct {
	ServerPort    int64
	SensorPort    int64
	RelayIP       string
	RelayServPort int64
	RelaySenrPort int64
	MaxSize       float64
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
	Channel             [8]Chan
	Reserved2           [5]uint32
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

type Sonar struct {
	Header            OIC_Header
	PortSidescan      [5825]int16
	StarboardSidescan [5825]int16
	SubBottom         [6000]int16
}
type Bathy struct {
	Header         OIC_Header
	PortAngle      [155]float32
	PortR          [155]float32
	StarboardAngle [300]float32
	StarboardR     [300]float32
}
