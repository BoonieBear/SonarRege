package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regener/oic"
	"strconv"
	"strings"
	"unicode/utf16"
)

func LoadCfg(cfgfile string) *oic.Cfg {
	//file open ok?
	file, err := os.Open(cfgfile)
	if err != nil {
		fmt.Printf("Open config file failed:%s\n", err.Error())
		return nil
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	if reader == nil {
		fmt.Println("Read config file failed")
		return nil
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read RecvPort title failed:%s\n", err.Error())
		return nil
	}
	if strings.TrimSpace(line) != "[RecvPort]" {
		fmt.Println("RecvPort part invalid!")
		return nil
	}
	line, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read ServerPort failed:%s\n", err.Error())
		return nil
	}
	serverport := extractInt32(line, "ServerPort")
	if serverport == 0 {
		fmt.Println("ServerPort part invalid!")
		return nil
	}
	line, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read SensorPort failed:%s\n", err.Error())
		return nil
	}
	sensorport := extractInt32(line, "SensorPort")
	if sensorport == 0 {
		fmt.Println("SensorPort part invalid!")
		return nil
	}
	line, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read Relay title failed:%s\n", err.Error())
		return nil
	}
	if strings.TrimSpace(line) != "[Relay]" {
		fmt.Println("Relay part invalid!")
		return nil
	}
	line, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read RelayIP failed:%s\n", err.Error())
		return nil
	}
	relayIP := extractString(line, "RelayIP")
	if relayIP == "" {
		fmt.Println("RelayIP part invalid!")
		return nil
	}
	line, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read RelayServPort failed:%s\n", err.Error())
		return nil
	}
	relayservport := extractInt32(line, "RelayServPort")
	if relayservport == 0 {
		fmt.Println("RelayServPort part invalid!")
		return nil
	}
	line, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read RelaySenrPort failed:%s\n", err.Error())
		return nil
	}
	relaysenrport := extractInt32(line, "RelaySenrPort")
	if relaysenrport == 0 {
		fmt.Println("RelaySenrPort part invalid!")
		return nil
	}
	line, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("read Size title failed:%s\n", err.Error())
		return nil
	}
	if strings.TrimSpace(line) != "[Size](Mbyte)" {
		fmt.Println("[Size](Mbyte) part invalid!")
		return nil
	}
	line, err = reader.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		fmt.Printf("read MaxFileSize failed:%s\n", err.Error())
		return nil
	}
	maxFileSize := extractFloat64(line, "MaxFileSize")
	if maxFileSize == 0 {
		fmt.Println("MaxFileSize part invalid!")
		return nil
	}
	cfg := oic.Cfg{
		ServerPort:    serverport,
		SensorPort:    sensorport,
		RelayIP:       relayIP,
		RelayServPort: relayservport,
		RelaySenrPort: relaysenrport,
		MaxSize:       maxFileSize,
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
func extractInt32(line string, keyword string) int64 {
	v := strings.TrimSpace(line)
	sa := strings.Split(v, "=")
	if sa[0] == keyword {
		f, _ := strconv.ParseInt(sa[1], 0, 32)
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
