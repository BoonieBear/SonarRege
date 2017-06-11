// structure
package oic

import (
	"regener/util"
)

//initialize bathy struct ,call OIC initialize
func (bathy *Bathy) Init() {
	header := &bathy.Header
	OICInit(header, true)
	bathy.Header.ChanNum = 5
	bathy.Header.ChanOffset[0] = 0
	bathy.Header.ChanOffset[1] = 620
	bathy.Header.ChanOffset[2] = 1240
	bathy.Header.ChanOffset[3] = 2440
	bathy.Header.ChanOffset[4] = 3640
	bathy.Header.Channel[0].Type = 7
	bathy.Header.Channel[0].Side = 0
	bathy.Header.Channel[0].Size = 3 //float32
	bathy.Header.Channel[0].Empt = 0
	bathy.Header.Channel[0].Frequency = 12
	bathy.Header.Channel[0].Samples = 155
	bathy.Header.Channel[1].Type = 7 //8?
	bathy.Header.Channel[1].Side = 0
	bathy.Header.Channel[1].Size = 3 //float32
	bathy.Header.Channel[1].Empt = 0
	bathy.Header.Channel[1].Frequency = 12
	bathy.Header.Channel[1].Samples = 155
	bathy.Header.Channel[2].Type = 7
	bathy.Header.Channel[2].Side = 1
	bathy.Header.Channel[2].Size = 3 //float32
	bathy.Header.Channel[2].Empt = 0
	bathy.Header.Channel[2].Frequency = 12
	bathy.Header.Channel[2].Samples = 300
	bathy.Header.Channel[3].Type = 7
	bathy.Header.Channel[3].Side = 1
	bathy.Header.Channel[3].Size = 3 //float32
	bathy.Header.Channel[3].Empt = 0
	bathy.Header.Channel[3].Frequency = 12
	bathy.Header.Channel[3].Samples = 155
	bathy.Header.Channel[4].Type = 9
	bathy.Header.Channel[4].Side = 0
	bathy.Header.Channel[4].Size = 1 //int16
	bathy.Header.Channel[4].Empt = 0
	bathy.Header.Channel[4].Frequency = 12
	bathy.Header.Channel[4].Samples = 56

}
func (bathy *Bathy) Pack() []byte {
	buf := bathy.Header.Pack()
	for i := 0; i < len(bathy.PortAngle); i++ {
		buf = append(buf, util.Float32ToByteBE(bathy.PortAngle[i])...)
	}
	for i := 0; i < len(bathy.PortR); i++ {
		buf = append(buf, util.Float32ToByteBE(bathy.PortR[i])...)
	}
	for i := 0; i < len(bathy.StarboardAngle); i++ {
		buf = append(buf, util.Float32ToByteBE(bathy.StarboardAngle[i])...)
	}
	for i := 0; i < len(bathy.StarboardR); i++ {
		buf = append(buf, util.Float32ToByteBE(bathy.StarboardR[i])...)
	}
	for i := 0; i < 56; i++ {
		buf = append(buf, util.IntToBytesBE(16, int64(bathy.Unknow[i]))...)
	}
	return buf
}
