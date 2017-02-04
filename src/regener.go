package main

import (
	"fmt"
	"io/ioutil"
	"oic"
	"os"
	"strconv"
	"strings"
)

func LoadCfg(string cfgfile) *Cfg {
	//file open ok?
	file, err := os.Open(cfgfile)
	if err != nil {
		fmt.Printf("Open config file failed:%s\n", err.Error())
		return nil
	}
	buff, err := ioutil.ReadFile(filename)
	s, err := utf16toString(buff[:])
	lines := strings.Split(s, "\n")
	if lines[0] != "[RecvPort]" {
		fmt.Println("RecvPort part invalid!")
		return nil
	}
	serverport := extractString(lines[1], "ServerPort")
	if serverport == "" {
		fmt.Println("ServerPort part invalid!")
		return nil
	}
	sensorport := extractString(lines[2], "SensorPort")
	if sensorport == "" {
		fmt.Println("SensorPort part invalid!")
		return nil
	}
	if lines[3] != "[Relay]" {
		fmt.Println("Relay part invalid!")
		return nil
	}
	relayIP := extractString(lines[4], "RelayIP")
	if relayIP == "" {
		fmt.Println("RelayIP part invalid!")
		return nil
	}
	relayservport := extractString(lines[5], "RelayServPort")
	if relayservport == "" {
		fmt.Println("RelayServPort part invalid!")
		return nil
	}
	relaysenrport := extractString(lines[6], "RelaySenrPort")
	if relaysenrport == "" {
		fmt.Println("RelaySenrPort part invalid!")
		return nil
	}
	if lines[7] != "[Size](Mbyte)" {
		fmt.Println("[Size](Mbyte) part invalid!")
		return nil
	}
	maxFileSize := extractString(lines[8], "MaxFileSize")
	if maxFileSize == "" {
		fmt.Println("MaxFileSize part invalid!")
		return nil
	}
	cfg := Cfg{
		ServerPort:    int32(serverport),
		SensorPort:    int32(sensorport),
		RelayIP:       relayIP,
		RelayServPort: int32(relayservport),
		RelaySenrPort: int32(relaysenrport),
	}
	return &cfg
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
func utf16toString(b []uint8) (string, error) {
	if len(b)&1 != 0 {
		return "", errors.New("len(b) must be even")
	}

	// Check BOM
	var bom int
	if len(b) >= 2 {
		switch n := int(b[0])<<8 | int(b[1]); n {
		case 0xfffe:
			bom = 1
			fallthrough
		case 0xfeff:
			b = b[2:]
		}
	}

	w := make([]uint16, len(b)/2)
	for i := range w {
		w[i] = uint16(b[2*i+bom&1])<<8 | uint16(b[2*i+(bom+1)&1])
	}
	return string(utf16.Decode(w)), nil
}
func main() {
	fmt.Println("Start SonarGenerator...")
	fmt.Println("Load Configuration from cfg.ini ...")
	config := LoadCfg("cfg.ini")
	if config == nil {
		fmt.Println("No valid configuration, exit...")
		return
	}

}
