package lib

import (
	"fmt"

	"github.com/go-ble/ble"
)

const (
	switchbotUUID    = "fd3d"
	nordicUUID       = "0d00"
	m5stackEnvUUID   = "2301"
	switchbotBot     = 0x48
	switchbotMeter   = 0x54
	switchbotMotion  = 0x73
	switchbotContact = 0x64
	switchbotCurtain = 0x63
	m5stackEnv       = 0x5A
	m5stackInroom    = 0x58
)

type Device struct {
	Addr    string `json:"macaddress"`
	Rawdata string `json:"sensorData"`
}
type SwitchbotDevice struct {
	Device
	Battery    byte   `json:"battery"`
	DeviceType string `json:"deviceType"`
}
type Bot struct {
	SwitchbotDevice
	Mode  byte `json:"mode"`
	State byte `json:"state"`
}
type Curtain struct {
	SwitchbotDevice
	WheterToAllowConnection byte `json:"wheterToAllowConnection"`
	CalibrationSituation    byte `json:"calibrationSituation"`
	MotionState             byte `json:"motionState"`
	Position                byte `json:"position"`
	LightLevel              byte `json:"lightLevel"`
	DeviceChain             byte `json:"deviceChain"`
}
type Meter struct {
	SwitchbotDevice
	TempertureHighAlert byte    `json:"tempertureHighAlert"`
	TempertureLowAlert  byte    `json:"tempertureLowAlert"`
	HumidityHighAlert   byte    `json:"humidityHighAlert"`
	HumidityLowAlert    byte    `json:"humidityLowAlert"`
	Temperature         float32 `json:"temperature"`
	Humidity            byte    `json:"humidity"`
	TemperatureScale    byte    `json:"temperatureScale"`
}

type Motion struct {
	SwitchbotDevice
	ScopeTested             byte  `json:"scopeTested"`
	RipState                byte  `json:"ripState"`
	LedState                byte  `json:"ledState"`
	IotState                byte  `json:"iotState"`
	SensingDistance         byte  `json:"sensingDistance"`
	LightIntensity          byte  `json:"lightIntensity"`
	SinceLastTriggerRipTime int32 `json:"sinceLastTriggerRipTime"`
}

type Contact struct {
	SwitchbotDevice
	ScopeTested     byte  `json:"scopeTested"`
	RipState        byte  `json:"ripState"`
	HalState        byte  `json:"halState"`
	LightLevel      byte  `json:"lightLevel"`
	RipUTC          int32 `json:"ripUTC"`
	HalUTC          int32 `json:"halUTC"`
	EntranceCount   byte  `json:"entranceCount"`
	GooutCount      byte  `json:"goOutCount"`
	ButtonPushCount byte  `json:"buttonPushCount"`
}
type InRoom struct {
	Device
	DeviceType string `json:"deviceType"`
	PushCount  byte   `json:"pushCount"`
	SendCount  byte   `json:"count"`
}
type M5StackEnv struct {
	Device
	DeviceType string `json:"deviceType"`
	Range      int32  `json:"range"`
	SendCount  byte   `json:"count"`
}

// FIXME
type Sensor interface {
	DeviceName() string
}

func (d *Bot) DeviceName() string {
	return "bot"
}

func (d *Curtain) DeviceName() string {
	return "curtain"
}
func (d *Contact) DeviceName() string {
	return "contact"
}
func (d *Motion) DeviceName() string {
	return "motion"
}
func (d *Meter) DeviceName() string {
	return "meter"
}
func (d *M5StackEnv) DeviceName() string {
	return "env"
}
func (d *InRoom) DeviceName() string {
	return "inroom"
}
func NewDevice(addr string, d []byte) *Device {
	return &Device{Addr: addr, Rawdata: fmt.Sprintf("%x", d)}
}
func NewSwitchbotDevice(addr string, sd ble.ServiceData) *SwitchbotDevice {
	return &SwitchbotDevice{Device: *NewDevice(addr, sd.Data),
		DeviceType: string(sd.Data[0] & 0b01111111),
		Battery:    sd.Data[2] & 0b01111111,
	}
}
func NewBot(addr string, sd ble.ServiceData) *Bot {
	return &Bot{
		SwitchbotDevice: *NewSwitchbotDevice(addr, sd),
		Mode:            sd.Data[1] >> 7,
		State:           (sd.Data[1] & 0b01000000) >> 6,
	}
}
func NewCurtain(addr string, sd ble.ServiceData) *Curtain {
	return &Curtain{
		SwitchbotDevice:         *NewSwitchbotDevice(addr, sd),
		WheterToAllowConnection: sd.Data[1] >> 7,
		CalibrationSituation:    (sd.Data[1] & 0b01000000) >> 6,
		MotionState:             sd.Data[3] >> 7,
		Position:                sd.Data[3] & 0b01111111,
		LightLevel:              (sd.Data[4] & 0b11110000) >> 4,
		DeviceChain:             sd.Data[4] & 0b00001111,
	}
}
func NewMeter(addr string, sd ble.ServiceData) *Meter {
	return &Meter{
		SwitchbotDevice:     *NewSwitchbotDevice(addr, sd),
		TempertureHighAlert: sd.Data[3] >> 7,
		TempertureLowAlert:  (sd.Data[3] & 0b01000000) >> 6,
		HumidityHighAlert:   (sd.Data[3] & 0b000100000) >> 5,
		HumidityLowAlert:    (sd.Data[3] & 0b000010000) >> 4,
		Temperature: func() float32 {
			positiveTempFlag := sd.Data[4] >> 7
			temperature := float32(sd.Data[3]&0b00001111)/10 + float32(sd.Data[4]&0b01111111)
			if positiveTempFlag == 0 {
				temperature = -temperature
			}
			return temperature
		}(),
		Humidity:         (sd.Data[5] & 0b011111111),
		TemperatureScale: sd.Data[5] >> 7,
	}
}
func NewMotion(addr string, sd ble.ServiceData) *Motion {
	return &Motion{
		SwitchbotDevice:         *NewSwitchbotDevice(addr, sd),
		ScopeTested:             sd.Data[1] >> 7,
		RipState:                (sd.Data[1] & 0b01000000) >> 6,
		LedState:                (sd.Data[5] & 0b00100000) >> 5,
		IotState:                (sd.Data[5] & 0b00010000) >> 4,
		SensingDistance:         (sd.Data[5] & 0b00001100) >> 2,
		LightIntensity:          (sd.Data[5] & 0b00000011),
		SinceLastTriggerRipTime: (int32(sd.Data[3]) << 8) + int32(sd.Data[4]),
	}
}

func NewContact(addr string, sd ble.ServiceData) *Contact {
	return &Contact{
		SwitchbotDevice: *NewSwitchbotDevice(addr, sd),
		ScopeTested:     sd.Data[1] >> 7,
		RipState:        (sd.Data[1] & 0b01000000) >> 6,
		HalState:        (sd.Data[3] & 0b00000110) >> 1,
		LightLevel:      sd.Data[3] & 0b00000001,
		RipUTC:          int32(sd.Data[3]&0b10000000)<<9 + int32(sd.Data[4])<<8 + int32(sd.Data[5]),
		HalUTC:          int32(sd.Data[3]&0b01000000)<<10 + int32(sd.Data[6])<<8 + int32(sd.Data[7]),
		EntranceCount:   (sd.Data[8] & 0b11000000) >> 6,
		GooutCount:      (sd.Data[8] & 0b00110000) >> 4,
		ButtonPushCount: sd.Data[8] & 0b00001111,
	}
}
func NewM5StackEnv(addr string, sd ble.ServiceData) *M5StackEnv {
	return &M5StackEnv{
		Device:     *NewDevice(addr, sd.Data),
		DeviceType: string(sd.Data[0] & 0b01111111),
		Range:      int32(sd.Data[1])<<8 + int32(sd.Data[2]),
		SendCount:  sd.Data[3],
	}
}
func NewInRoom(addr string, sd []byte) *InRoom {
	return &InRoom{
		Device:     *NewDevice(addr, sd),
		PushCount:  sd[0],
		DeviceType: "X",
	}
}
func GetDevice(addr string, sd ble.ServiceData) Sensor {
	if sd.UUID.String() == switchbotUUID || sd.UUID.String() == nordicUUID || sd.UUID.String() == m5stackEnvUUID {
		switch sd.Data[0] {
		case switchbotBot:
			return NewBot(addr, sd)
		case switchbotCurtain:
			return NewCurtain(addr, sd)
		case switchbotMeter:
			return NewMeter(addr, sd)
		case switchbotMotion:
			return NewMotion(addr, sd)
		case switchbotContact:
			return NewContact(addr, sd)
		case m5stackEnv:
			return NewM5StackEnv(addr, sd)
		default:
			return nil
		}
	}
	return nil
}
