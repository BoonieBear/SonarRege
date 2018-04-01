// structure
package oic

import (
	//	"fmt"
	"regener/util"
)

//initialize sonar struct ,call OIC initialize
func (sonar *Sonar) Init() {
	header := &sonar.Header
	OICInit(header, false)
	header.ChanNum = 3
	header.ChanOffset[0] = 0
	header.ChanOffset[1] = 11650
	header.ChanOffset[2] = 23304
	header.Channel[0].Type = 0 //sonar type:0 = sidescan 1 = angle 2 = multibeam
	header.Channel[0].Side = 0 //sonar side:0 = port 1 = starboard
	//data sample type and size:0 = 1 byte integer 1 = 2 byte integer 2 = 4 byte integer 3 = 4 byte float 4 = 12 byte set of three floats - range, theta, amp */
	header.Channel[0].Size = 1

	header.Channel[1].Type = 0
	header.Channel[1].Side = 1
	header.Channel[1].Size = 1

	header.Channel[2].Type = 0
	header.Channel[2].Side = 2
	header.Channel[2].Size = 1
}

func (sonar *Sonar) Pack() []byte {
	datasize := (len(sonar.PortSidescan) + len(sonar.StarboardSidescan) + len(sonar.SubBottom)) * 2
	sonar.Header.DataSize = uint32(datasize)
	sonar.Header.ChanOffset[1] = uint32(len(sonar.PortSidescan) * 2)
	sonar.Header.ChanOffset[2] = uint32(len(sonar.PortSidescan)*2 + len(sonar.StarboardSidescan)*2)
	sonar.Header.Channel[0].Samples = uint32(len(sonar.PortSidescan))
	sonar.Header.Channel[1].Samples = uint32(len(sonar.StarboardSidescan))
	sonar.Header.Channel[2].Samples = uint32(len(sonar.SubBottom))

	buf := sonar.Header.Pack()
	for i := 0; i < len(sonar.PortSidescan); i++ {
		buf = append(buf, util.IntToBytesLE(16, int64(sonar.PortSidescan[i]))...)
	}
	for i := 0; i < len(sonar.StarboardSidescan); i++ {
		buf = append(buf, util.IntToBytesLE(16, int64(sonar.StarboardSidescan[i]))...)
	}
	for i := 0; i < len(sonar.SubBottom); i++ {
		buf = append(buf, util.IntToBytesLE(16, int64(sonar.SubBottom[i]))...)
	}

	return buf
}
