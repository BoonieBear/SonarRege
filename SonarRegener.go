// SonarRegener
package SonarRegener

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

type Chan struct {
	Type      uint8
	Side      uint8
	Size      uint8
	Empt      uint8
	Frequency uint32
	Samples   uint32
}
type Cfg struct {
	ServerPort    int32
	SensorPort    int32
	RelayServPort int32
	RelaySenrPort int32
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
func OICInit(header *OIC_Header) {
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
func LoadCfg(string cfgfile) *Cfg {
	cfg := Cfg{}
	//file open ok?
	file, err := os.Open(cfgfile)
	if err != nil {
		return ant, err
	}
}
func extractString(line string, keyword string) string {
	v := strings.TrimSpace(line)
	sa := strings.Split(v, "=")
	if sa[0] == keyword {
		return sa[1]
	}
	return ""
}
func extractFloat64(line string, keyword string) float64 {
	v := strings.TrimSpace(line)
	sa := strings.Split(v, "=")
	if sa[0] == keyword {
		f, _ := strconv.ParseFloat(sa[1], 64)
		return f
	}
	return 0
}
func main() {
	fmt.Println("Start SonarGenerator...")
	fmt.Println("Load Configuration from .ini file...")
	err := LoadCfg("Cfg.ini")
	if err != nil {

	}
}
