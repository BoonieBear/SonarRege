package sensor

type bsss struct {
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

type datapara struct {
}

//side scan
type ss struct {
}

//subbottom
type sb struct {
}

//bathy scan
type bathy struct {
}
