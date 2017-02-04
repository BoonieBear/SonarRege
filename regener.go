// main
package main

import (
	"bathy"
	"define"
	"fmt"
	"io/ioutil"
	"os"
	"sonar"
	"strconv"
	"strings"
)

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
