package util

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	//	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

type Cfg struct {
	ServerPort    int32
	SensorPort    int32
	RelayIP       string
	RelayServPort int32
	RelaySenrPort int32
	MaxSize       float64
}

type Logger struct {
	logger  *log.Logger
	logfile *os.File
}

func (l *Logger) New(file string) {
	logfile, err := os.Create(file)
	if err != nil {
		log.Fatalln("Create log file failed.")

	}
	l.logger = log.New(logfile, "", log.LstdFlags)

}

func (l *Logger) Println(msg string) {
	log.Println(msg)
	l.logger.Println(msg)
}

func (l *Logger) Fatal(msg string) {
	l.logger.Println(msg)
	log.Fatalln(msg)
}

func LoadCfg(cfgfile string) *Cfg {
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
	cfg := Cfg{
		ServerPort:    serverport,
		SensorPort:    sensorport,
		RelayIP:       relayIP,
		RelayServPort: relayservport,
		RelaySenrPort: relaysenrport,
		MaxSize:       maxFileSize,
	}
	return &cfg
}
func (cfg *Cfg) Dump() {
	if cfg != nil {
		fmt.Println("====List all items in Config file=====")
		fmt.Printf("== Server Port: %d \n", cfg.ServerPort)
		fmt.Printf("== Sensor Port: %d \n", cfg.SensorPort)
		fmt.Printf("== Relay IP: %s \n", cfg.RelayIP)
		fmt.Printf("== Relay Server Port: %d \n", cfg.RelayServPort)
		fmt.Printf("== Relay Sensor Port: %d \n", cfg.RelaySenrPort)
		fmt.Printf("== Max File Size(M): %f \n", cfg.MaxSize)

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
func extractInt32(line string, keyword string) int32 {
	v := strings.TrimSpace(line)
	sa := strings.Split(v, "=")
	if sa[0] == keyword {
		f, _ := strconv.ParseInt(sa[1], 0, 32)
		return int32(f)
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

//////LittleEndian//////
func IntToBytesLE(bits int32, i int64) []byte {
	buf := make([]byte, bits/8)
	switch bits / 8 {
	case 1:
		buf[0] = byte(i)
		return buf
	case 2:
		binary.LittleEndian.PutUint16(buf, uint16(i))
		return buf
	case 4:
		binary.LittleEndian.PutUint32(buf, uint32(i))
		return buf
	case 8:
		binary.LittleEndian.PutUint64(buf, uint64(i))
		return buf
	default:
		return nil

	}
}
func UIntToBytesLE(bits int32, i uint64) []byte {
	buf := make([]byte, bits/8)
	switch bits / 8 {
	case 1:
		buf[0] = byte(i)
		return buf
	case 2:
		binary.LittleEndian.PutUint16(buf, uint16(i))
		return buf
	case 4:
		binary.LittleEndian.PutUint32(buf, uint32(i))
		return buf
	case 8:
		binary.LittleEndian.PutUint64(buf, uint64(i))
		return buf
	default:
		return nil

	}
}
func BytesToUIntLE(bits int32, buf []byte) uint64 {
	switch bits / 8 {
	case 1:
		return uint64(buf[0])
	case 2:
		return uint64(binary.LittleEndian.Uint16(buf))
	case 4:
		return uint64(binary.LittleEndian.Uint32(buf))
	case 8:
		return uint64(binary.LittleEndian.Uint64(buf))
	default:
		return 0
	}
}

func BytesToIntLE(bits int32, buf []byte) int64 {
	return int64(BytesToUIntLE(bits, buf))
}
func Float32ToByteLE(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32LE(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToByteLE(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64LE(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

//////BigEndian//////
func IntToBytesBE(bits int32, i int64) []byte {
	buf := make([]byte, bits/8)
	switch bits / 8 {
	case 1:
		buf[0] = byte(i)
		return buf
	case 2:
		binary.BigEndian.PutUint16(buf, uint16(i))
		return buf
	case 4:
		binary.BigEndian.PutUint32(buf, uint32(i))
		return buf
	case 8:
		binary.BigEndian.PutUint64(buf, uint64(i))
		return buf
	default:
		return nil

	}
}
func UIntToBytesBE(bits int32, i uint64) []byte {
	buf := make([]byte, bits/8)
	switch bits / 8 {
	case 1:
		buf[0] = byte(i)
		return buf
	case 2:
		binary.BigEndian.PutUint16(buf, uint16(i))
		return buf
	case 4:
		binary.BigEndian.PutUint32(buf, uint32(i))
		return buf
	case 8:
		binary.BigEndian.PutUint64(buf, uint64(i))
		return buf
	default:
		return nil

	}
}
func BytesToUIntBE(bits int32, buf []byte) uint64 {
	switch bits / 8 {
	case 1:
		return uint64(buf[0])
	case 2:
		return uint64(binary.BigEndian.Uint16(buf))
	case 4:
		return uint64(binary.BigEndian.Uint32(buf))
	case 8:
		return uint64(binary.BigEndian.Uint64(buf))
	default:
		return 0
	}
}

func BytesToIntBE(bits int32, buf []byte) int64 {
	return int64(BytesToUIntBE(bits, buf))
}
func Float32ToByteBE(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32BE(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToByteBE(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64BE(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

//save file
type tracefile struct {
	writer    *bufio.Writer
	File      *os.File
	Count     uint32
	FileName  string
	MaxSize   uint32
	FileIndex uint16
}

func (tf *tracefile) New(pre string, maxlength uint32) error {
	if tf.writer != nil {
		return errors.New("Already create a trace file instance")
	}
	t := time.Now().Format("060102150405")
	path, _ := os.Getwd()
	tf.FileName = filepath.Join(path, pre+t+"-1")
	tf.MaxSize = maxlength

	var err error
	tf.File, err = os.Create(tf.FileName) //create file
	if err != nil {
		return err
	}
	tf.FileIndex = 1
	tf.Count = 0
	tf.writer = bufio.NewWriter(tf.File)
	return nil
}
func (tf *tracefile) Write(bytes []byte) error {
	if tf.writer == nil {
		return errors.New("no valid trace file")
	}
	n, err := tf.writer.Write(bytes)

	if err != nil {
		if tf.File != nil {
			tf.File.Close()
		}
		return err
	}
	err = tf.writer.Flush()
	if err != nil {
		if tf.File != nil {
			tf.File.Close()
		}
		return err
	}
	tf.Count += uint32(n)
	if tf.Count > tf.MaxSize {
		tf.Close()
		filename := strings.Split(filepath.Base(tf.FileName), "-")[0]

		tf.FileIndex += 1
		path, _ := os.Getwd()
		tf.FileName = filepath.Join(path, filename+"-"+strconv.Itoa(int(tf.FileIndex)))

		tf.File, err = os.Create(tf.FileName) //create file
		if err != nil {
			return err
		}

		tf.Count = 0
		tf.writer = bufio.NewWriter(tf.File)
	}
	return nil
}
func (tf *tracefile) Close() {
	if tf.writer != nil {
		tf.writer.Flush()
	}
	if tf.File != nil {
		tf.File.Close()
	}
}
